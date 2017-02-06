package domain

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/domainType"
	"github.com/distributed-vision/go-resources/resolvers"
)

type SelectorOpts struct {
	IgnoreCase       bool
	IgnoreWhitespace bool
}

type Selector struct {
	Scope ids.DomainScope
	Id    []byte
	Name  string
	Opts  SelectorOpts
}

func (this *Selector) Test(candidate interface{}) bool {
	return true
}

func (this *Selector) Key() interface{} {
	return this.Id
}

var entityType ids.TypeIdentifier

func (this *Selector) Type() ids.TypeIdentifier {
	if entityType == nil {
		entityType = ids.IdOfType(reflect.TypeOf((*ids.Domain)(nil)).Elem())
	}

	return entityType
}

func Get(selector Selector) (domain ids.Domain, err error) {
	cres, cerr := Resolve(context.Background(), selector)

	if cres == nil || cerr == nil {
		return nil, fmt.Errorf("domain.Get Failed: domain.Resolve channels are undefined")
	}

	select {
	case domain = <-cres:
		if domain != nil {
			break
		}
	case err = <-cerr:
		if err != nil {
			break
		}
	}

	return domain, err
}

var domainResolvers []resolvers.Resolver
var domainMaps = make(map[string]map[string]ids.Domain)
var resolveMutex = &sync.Mutex{}

func Resolve(resolutionContext context.Context, selector Selector) (<-chan ids.Domain, <-chan error) {

	var result ids.Domain

	resolveMutex.Lock()
	domainMap, ok := domainMaps[base62.Encode(selector.Scope.Id())]

	if ok {
		result = domainMap[base62.Encode(selector.Id)]
	}

	resolveMutex.Unlock()

	cResOut := make(chan ids.Domain, 1)
	cErrOut := make(chan error, 1)

	if result != nil {
		cResOut <- result
		close(cResOut)
		close(cErrOut)
		return cResOut, cErrOut
	}

	if domainResolvers == nil {
		resolveMutex.Lock()

		if domainResolvers == nil {
			domainResolvers = make([]resolvers.Resolver, 0)
			// TODO only need to do thsi once - probably should just call resolve and
			// move the merge to there
			for _, resolverInfo := range selector.Scope.InfoValue("resolverInfo").(map[interface{}]resolvers.ResolverInfo) {
				resolverFactory, err := resolvers.NewResolverFactory(resolverInfo)

				if err != nil {
					cresolver, cerr := resolverFactory.New(resolutionContext)

					select {
					case resolver := <-cresolver:
						if resolver != nil {
							domainResolvers = append(domainResolvers, resolver)
						}
					case err := <-cerr:
						if err != nil {
							continue
						}
					case <-resolutionContext.Done():
						close(cResOut)
						close(cErrOut)
						return cResOut, cErrOut
					}
				}
			}
		}

		resolveMutex.Unlock()
	}

	var wg sync.WaitGroup
	mergeContext, cancel := context.WithCancel(resolutionContext)
	errors := make([]error, 0)
	resultMutex := &sync.Mutex{}

	resolve := func(resolver resolvers.Resolver) {
		defer wg.Done()

		cres, cerr := resolver.Resolve(mergeContext, &selector)

		if cres == nil || cerr == nil {
			errors = append(errors, fmt.Errorf("domain.Get Failed: domain.Resolve channels are undefined"))
		}

		select {
		case res := <-cres:
			if res != nil {
				if result == nil {
					resultMutex.Lock()
					if result == nil {
						result = res.(ids.Domain)
						cancel()
					}
					resultMutex.Unlock()
				}
			}
		case err := <-cerr:
			if err != nil {
				errors = append(errors, err)
			}
		case <-mergeContext.Done():
		}
	}

	wg.Add(len(domainResolvers))
	for _, r := range domainResolvers {
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

//	_resolverFactories.set('github', GitHubResolver.createResolver)
//	_resolverFactories.set('file', FileResolver.createResolver)

type unmarshaller func(json map[string]interface{}, opts map[string]interface{}) (ids.Domain, error)

var unmarshalers map[ids.DomainType]unmarshaller = make(map[ids.DomainType]unmarshaller)

func RegisterJSONUnmarshaller(domainType ids.DomainType, unmarshaller unmarshaller) {
	unmarshalers[domainType] = unmarshaller
}

func unmarshalJSON(json map[string]interface{}, opts map[string]interface{}) (ids.Domain, error) {
	dt, err := domainType.Parse(json["domainType"].(string))

	if err != nil {
		return nil, err
	}

	unmarshaler, ok := unmarshalers[dt]

	if !ok {
		return nil, errors.New("Unknown domain type: " + json["domainType"].(string))
	}

	return unmarshaler(json, opts)
}
