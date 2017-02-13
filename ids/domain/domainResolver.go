package domain

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/domainType"
	"github.com/distributed-vision/go-resources/resolvers"
	"github.com/distributed-vision/go-resources/translators"
)

func domainTranslator(translationContext context.Context, fromId ids.Identifier, fromValue interface{}) (chan interface{}, chan error) {

	cres := make(chan interface{}, 1)
	cerr := make(chan error, 1)

	json := fromValue.(map[string]interface{})
	json["id"] = string(fromId.Id())

	toValue, err := unmarshalJSON(translationContext, json)
	//fmt.Printf("id: %+v val: %+v err: %s\n", fromId.Id(), toValue, err)

	if err != nil {
		cerr <- err
	} else {
		cres <- toValue
	}

	close(cres)
	close(cerr)

	return cres, cerr
}

var domainEntityType ids.TypeIdentifier

func init() {
	ids.OnLocalTypeInit(func() {

		if domainEntityType == nil {
			domainEntityType = ids.NewLocalTypeId(reflect.TypeOf((*ids.Domain)(nil)).Elem())
		}

		mapType := ids.NewLocalTypeId(reflect.TypeOf(map[string]interface{}{}))
		translators.Register(context.Background(), mapType, domainEntityType, domainTranslator)
	})
}

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
	domain, ok := candidate.(ids.Domain)

	if !ok {
		return false
	}

	if this.Id != nil && !bytes.Equal(this.Id, domain.IdRoot()) {
		return false
	}

	if this.Name != "" {
		if this.Opts.IgnoreCase {
			if strings.ToUpper(this.Name) != strings.ToUpper(domain.Name()) {
				return false
			}
		} else {
			if this.Name != domain.Name() {
				return false
			}
		}
	}

	return true
}

func (this *Selector) Key() interface{} {
	return base62.Encode(this.Id)
}

var entityType ids.TypeIdentifier

func (this *Selector) Type() ids.TypeIdentifier {
	if entityType == nil {
		entityType = ids.NewLocalTypeId(reflect.TypeOf((*ids.Domain)(nil)).Elem())
	}

	return entityType
}

func Get(resolutionContext context.Context, selector Selector) (domain ids.Domain, err error) {
	//fmt.Printf("scope=%+v\n", selector.Scope)
	res, err := resolvers.Get(resolutionContext, &selector)

	if err != nil {
		return nil, err
	}

	if domain, ok := res.(ids.Domain); ok {
		return domain, err
	}

	return nil, fmt.Errorf("Resolver returned invalid type, expected: ids.Domain got: %s", reflect.TypeOf(res))
}

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

	cresIn, cerrIn := resolvers.Resolve(resolutionContext, &selector)

	go func() {
		resolved := false
		for !resolved {
			select {
			case resIn, ok := <-cresIn:
				if ok {
					if domain, ok := resIn.(ids.Domain); ok {
						cResOut <- domain
					} else {
						cErrOut <- fmt.Errorf("Resolver returned invalid type, expected: ids.Domain got: %s", reflect.TypeOf(resIn))
					}
					resolved = true
				}
			case errIn, ok := <-cerrIn:
				if ok {
					cErrOut <- errIn
					resolved = true
				}
			}
		}
		close(cResOut)
		close(cErrOut)
	}()

	return cResOut, cErrOut
}

type unmarshaller func(unmarshalContext context.Context, json map[string]interface{}) (ids.Domain, error)

var unmarshalers map[ids.DomainType]unmarshaller = make(map[ids.DomainType]unmarshaller)

func RegisterJSONUnmarshaller(domainType ids.DomainType, unmarshaller unmarshaller) {
	unmarshalers[domainType] = unmarshaller
}

func unmarshalJSON(unmarshalContext context.Context, json map[string]interface{}) (ids.Domain, error) {
	dt, err := domainType.Parse(json["domainType"].(string))
	//fmt.Printf("dt=%v\n", dt)
	if err != nil {
		return nil, err
	}

	unmarshaler, ok := unmarshalers[dt]

	if !ok {
		return nil, errors.New("Unknown domain type: " + json["domainType"].(string))
	}

	return unmarshaler(unmarshalContext, json)
}
