package resolvers

import (
	"context"
	"fmt"
	"sync"

	"github.com/distributed-vision/go-resources/ids"
)

type Selector interface {
	Type() ids.TypeIdentifier
	Key() interface{}
	Test(entity interface{}) bool
}

type Resolver interface {
	Resolve(resolutionContext context.Context, selector Selector) (chan interface{}, chan error)
}

type ResolverInfo interface {
	EntityTypes() []ids.TypeIdentifier
	Matches(selector Selector) bool
	Value(key interface{}) interface{}
	WithValue(key, value interface{}) ResolverInfo
	WithValues(values map[interface{}]interface{}) ResolverInfo
}

type ResolverFactory interface {
	ResolverType() ids.TypeIdentifier
	ResolverInfo() ResolverInfo
	New(resolutionContext context.Context) (chan Resolver, chan error)
}

func Get(selector Selector) (entity interface{}, err error) {
	cres, cerr := Resolve(context.Background(), selector)

	if cres == nil || cerr == nil {
		return nil, fmt.Errorf("resolvers.Get Failed: resolvers.Resolve channels are undefined")
	}

	select {
	case entity = <-cres:
	case err = <-cerr:
	}

	return entity, err
}

type registryEntry struct {
	entityType ids.TypeIdentifier
	factory    ResolverFactory
	resolver   Resolver
}

var resolverRegistry map[string][]*registryEntry = make(map[string][]*registryEntry)
var resolverRegistryMutex = &sync.Mutex{}

func Resolve(resolutionContext context.Context, selector Selector) (chan interface{}, chan error) {
	cResOut := make(chan interface{}, 1)
	cErrOut := make(chan error, 1)

	var resolverEntries = []*registryEntry{}
	var err error

	if selector != nil {
		//fmt.Printf("reolvers.Resolve: %s: %+v\n", selector.Type(), selector)
		resolverRegistryMutex.Lock()
		if entries, ok := resolverRegistry[string(selector.Type().Value())]; ok {
			for _, entry := range entries {
				if entry.factory.ResolverInfo().Matches(selector) {
					resolverEntries = append(resolverEntries, entry)
				}
			}
		} else {
			err = fmt.Errorf("resolvers.Resolve Failed: no resolver for entity type=%s", selector.Type())
		}
		resolverRegistryMutex.Unlock()
	} else {
		err = fmt.Errorf("resolvers.Resolve Failed: selector cannot be nil")
	}

	if err != nil {
		cErrOut <- err
		close(cResOut)
		close(cErrOut)
		return cResOut, cErrOut
	}

	if len(resolverEntries) == 0 {
		cErrOut <- fmt.Errorf("resolvers.Resolve Failed: no resolver for entity type=%s", selector.Type())
		close(cResOut)
		close(cErrOut)
		return cResOut, cErrOut
	}

	var wg sync.WaitGroup
	mergeContext, cancel := context.WithCancel(resolutionContext)
	errors := make([]error, 0)
	resultMutex := &sync.Mutex{}
	var result interface{}

	resolve := func(resolverEntry *registryEntry) {
		defer wg.Done()

		if resolverEntry.resolver == nil {
			cres, cerr := resolverEntry.factory.New(mergeContext)

			for resolverEntry.resolver == nil {
				select {
				case resolver, ok := <-cres:
					if ok {
						resolverEntry.resolver = resolver
					}
				case err, ok := <-cerr:
					if ok {
						errors = append(errors, err)
						return
					}
				case <-mergeContext.Done():
					return
				}
			}
		}

		cres, cerr := resolverEntry.resolver.Resolve(mergeContext, selector)

		if cres == nil || cerr == nil {
			errors = append(errors, fmt.Errorf("domain.Get Failed: domain.Resolve channels are undefined"))
		}

		for result == nil {
			select {
			case res, ok := <-cres:
				if ok {
					if result == nil {
						resultMutex.Lock()
						if result == nil {
							result = res
							cancel()
						}
						resultMutex.Unlock()
					}
					return
				}
			case err, ok := <-cerr:
				if ok {
					errors = append(errors, err)
					return
				}
			case <-mergeContext.Done():
				return
			}
		}
	}

	wg.Add(len(resolverEntries))
	for _, r := range resolverEntries {
		go resolve(r)
	}

	go func() {
		wg.Wait()
		if result != nil {
			cResOut <- result
		} else if len(errors) > 0 {
			cErrOut <- fmt.Errorf("Resolve failed with the folling errors %v", errors)
		}
		close(cResOut)
		close(cErrOut)
	}()

	return cResOut, cErrOut
}

func RegisterFactory(resolverFactory ResolverFactory) {
	resolverRegistryMutex.Lock()
	defer resolverRegistryMutex.Unlock()

	for _, entityType := range resolverFactory.ResolverInfo().EntityTypes() {

		if entries, ok := resolverRegistry[string(entityType.Value())]; ok {
			for _, entry := range entries {
				if entry.factory.ResolverType().Equals(resolverFactory.ResolverType()) {
					// ignore duplicate entries
					// TODO look at resolver info to determine duplication
					return
				}
			}

			resolverRegistry[string(entityType.Value())] = append(entries, &registryEntry{entityType, resolverFactory, nil})
		} else {
			entries = []*registryEntry{&registryEntry{entityType, resolverFactory, nil}}
			resolverRegistry[string(entityType.Value())] = entries
		}

	}
}

func NewResolverFactory(resolverInfo ResolverInfo) (ResolverFactory, error) {
	return nil, nil
}

type resolverInfo struct {
	entityTypes []ids.TypeIdentifier
	//entityDomains []ids.Domain
	values map[interface{}]interface{}
	parent ResolverInfo
}

func NewResolverInfo(entityTypes []ids.TypeIdentifier, values map[interface{}]interface{}) ResolverInfo {
	return &resolverInfo{entityTypes, values, nil}
}

func (this *resolverInfo) EntityTypes() []ids.TypeIdentifier {
	return this.entityTypes
}

func (this *resolverInfo) Value(key interface{}) (res interface{}) {
	if this.values != nil {
		res = this.values[key]
	}

	if res == nil && this.parent != nil {
		res = this.parent.Value(key)
	}

	return res
}

func (this *resolverInfo) Matches(selector Selector) bool {
	return true
}

func (this *resolverInfo) WithValue(key, value interface{}) ResolverInfo {
	values := make(map[interface{}]interface{})
	values[key] = value
	return &resolverInfo{this.entityTypes, values, this}
}

func (this *resolverInfo) WithValues(values map[interface{}]interface{}) ResolverInfo {
	return &resolverInfo{this.entityTypes, values, this}
}
