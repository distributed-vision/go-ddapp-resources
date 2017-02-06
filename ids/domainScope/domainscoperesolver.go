package domainScope

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/encoding/encoderType"
	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/domain"
	"github.com/distributed-vision/go-resources/resolvers"
	"github.com/distributed-vision/go-resources/translators"
)

var entityType ids.TypeIdentifier

var untypedDomain []byte = domain.MustDecodeId(encoderType.BASE62, "3", "")

func domainScopeTranslator(translationContext context.Context, fromId ids.Identifier, fromValue interface{}) (chan interface{}, chan error) {

	cres := make(chan interface{}, 1)
	cerr := make(chan error, 1)

	json := fromValue.(map[string]interface{})
	json["id"] = string(fromId.Id())
	toValue, err := unmarshalJSON(json, map[string]interface{}{})
	//fmt.Printf("val: %+v err: %s\n", toValue, err)

	if err != nil {
		cerr <- err
	} else {
		cres <- toValue
	}

	close(cres)
	close(cerr)

	return cres, cerr
}

func init() {
	ids.OnLocalTypeInit(func() {
		if entityType == nil {
			entityType = ids.IdOfType(reflect.TypeOf((*ids.DomainScope)(nil)).Elem())
		}

		mapType := ids.IdOfType(reflect.TypeOf(map[string]interface{}{}))
		translators.Register(context.Background(), mapType, entityType, domainScopeTranslator)
	})
}

type SelectorOpts struct {
	IgnoreCase       bool
	IgnoreWhitespace bool
}

type Selector struct {
	Id   []byte
	Name string
	Opts SelectorOpts
}

func (this *Selector) Test(candidate interface{}) bool {
	scope, ok := candidate.(ids.DomainScope)

	if !ok {
		return false
	}

	if this.Id != nil && !bytes.Equal(this.Id, scope.Id()) {
		return false
	}

	if this.Name != "" {
		if this.Opts.IgnoreCase {
			if strings.ToUpper(this.Name) != strings.ToUpper(scope.Name()) {
				return false
			}
		} else {
			if this.Name != scope.Name() {
				return false
			}
		}
	}

	return true
}

func (this *Selector) Key() interface{} {
	return base62.Encode(this.Id)
}

func (this *Selector) Type() ids.TypeIdentifier {
	return entityType
}

func Get(selector Selector) (ids.DomainScope, error) {
	res, err := resolvers.Get(&selector)

	if scope, ok := res.(ids.DomainScope); ok {
		return scope, err
	}

	return nil, fmt.Errorf("Resolver returned invalid type, expected: ids.DomainScope got: %s", reflect.TypeOf(res))
}

func Resolve(resolutionContext context.Context, selector Selector) (chan ids.DomainScope, chan error) {
	cresOut := make(chan ids.DomainScope, 1)
	cerrOut := make(chan error, 1)

	cresIn, cerrIn := resolvers.Resolve(resolutionContext, &selector)

	go func() {
		select {
		case resIn := <-cresIn:
			if scope, ok := resIn.(ids.DomainScope); ok {
				cresOut <- scope
			} else {
				cerrOut <- fmt.Errorf("Resolver returned invalid type, expected: ids.DomainScope got: %s", reflect.TypeOf(resIn))
			}
		case errIn := <-cerrIn:
			cerrOut <- errIn
		}

		close(cresOut)
		close(cerrOut)
	}()

	return cresOut, cerrOut
}
