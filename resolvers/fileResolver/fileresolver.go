package fileResolver

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"

	"github.com/distributed-vision/go-resources/encoding/encoderType"
	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/domain"
	"github.com/distributed-vision/go-resources/ids/identifier"
	"github.com/distributed-vision/go-resources/ids/mappings"
	"github.com/distributed-vision/go-resources/resolvers"
	"github.com/distributed-vision/go-resources/translators"
	"github.com/distributed-vision/go-resources/types"
	"github.com/distributed-vision/go-resources/types/gotypeid"
	"github.com/fatih/structs"
)

var contentType ids.TypeIdentifier = gotypeid.IdOf(reflect.TypeOf(map[string]interface{}{}))

type fileResolver struct {
	path      string
	opts      Opts
	entityMap map[string]interface{}
}

type Opts struct {
	Paths         []string
	UnmarshalJSON func(values map[string]interface{}, opts map[string]interface{}) (interface{}, error)
}

var resolverMap map[string]*fileResolver = make(map[string]*fileResolver)

var resolverType ids.TypeIdentifier = gotypeid.IdOf(reflect.TypeOf(fileResolver{}))

func init() {
	baseType := types.IdOf("T0-FileResolver")
	mappings.Add(resolverType, baseType, nil, nil)
	//translators.Register(context.Background(), resolverType, baseType, baseTypeTranslator)
	//resolverFactories.Register(baseType, &factory{})
}

func baseTypeTranslator(translationContext context.Context, resoverInfo interface{}) (chan interface{}, chan error) {
	return nil, nil
}

func ResolverType() ids.TypeIdentifier {
	return resolverType
}

type factory struct {
	resolverInfo resolvers.ResolverInfo
}

func NewResolverFactory(resolverInfo resolvers.ResolverInfo) (resolvers.ResolverFactory, error) {
	return &factory{resolverInfo}, nil
}

func (this *factory) New(resolutionContext context.Context) (chan resolvers.Resolver, chan error) {
	resolver, err := New(this.resolverInfo.Value("file").(string), Opts{Paths: this.resolverInfo.Value("paths").([]string)})

	cres, cerr := make(chan resolvers.Resolver, 1), make(chan error, 1)
	if err != nil {
		cerr <- err
	} else {
		cres <- resolver
	}

	close(cres)
	close(cerr)

	return cres, cerr
}

func (this *factory) ResolverType() ids.TypeIdentifier {
	return resolverType
}

func (this *factory) ResolverInfo() resolvers.ResolverInfo {
	return this.resolverInfo
}

func New(locator string, opts Opts) (resolvers.Resolver, error) {
	resolver, ok := resolverMap[locator]

	if ok {
		return resolver, nil
	}

	resolver, err := newResolver(locator, opts)

	if err != nil {
		return nil, err
	}

	resolverMap[locator] = resolver

	return resolver, nil
}

func newResolver(file string, opts Opts) (*fileResolver, error) {

	var filePath string

	if opts.Paths != nil {

		for _, respath := range opts.Paths {
			respath = path.Join(respath, file)

			_, err := os.Stat(respath)

			if err == nil {
				filePath = respath
				break
			}
		}
	} else {
		_, err := os.Stat(file)

		if err == nil {
			filePath = file
		}
	}

	if filePath == "" {
		if opts.Paths == nil {
			return nil, fmt.Errorf("Can't find: %s", file)
		}

		return nil, fmt.Errorf("Can't find: %s in: %v", file, opts.Paths)
	}

	return &fileResolver{path: filePath, opts: opts}, nil
}

func (this *fileResolver) Resolve(context context.Context, selector resolvers.Selector) (chan interface{}, chan error) {

	cres, cerr := make(chan interface{}), make(chan error)

	go func() {
		entityMap, err := this.getMap(context, selector.Type())

		if err != nil {
			cerr <- err
		} else if entityMap != nil {
			entity, ok := entityMap[selector.Key().(string)]

			if ok {
				if selector.Test(entity) {
					cres <- entity
					close(cres)
					close(cerr)
					return
				}
			}

			for _, entity := range entityMap {
				if selector.Test(entity) {
					cres <- entity
					close(cres)
					close(cerr)
					return
				}
			}

			cerr <- fmt.Errorf("Invalid entity selector: %+v", selector)
		}

		close(cres)
		close(cerr)
	}()

	return cres, cerr
}

func (this *fileResolver) getMap(context context.Context, targetType ids.TypeIdentifier) (map[string]interface{}, error) {
	if this.entityMap != nil {
		return this.entityMap, nil
	}

	entityMap, err := this.loadMap(context, targetType)

	if err == nil {
		this.entityMap = entityMap
	}

	return entityMap, err
}

var unypedLocalDomain []byte = domain.MustDecodeId(encoderType.BASE62, "3", "")

func (this *fileResolver) loadMap(context context.Context, targetType ids.TypeIdentifier) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(this.path)

	if err != nil {
		return nil, err
	}

	var jsonEntities interface{}
	err = json.Unmarshal(data, &jsonEntities)

	if err != nil {
		return nil, err
	}

	jsonEntityMap := jsonEntities.(map[string]interface{})

	opts := map[string]interface{}{
		"path": this.path}

	for key, value := range structs.Map(this.opts) {
		opts[key] = value
	}

	//fmt.Printf("Content Type=%+v, Target Type %+v\n", contentType, targetType)

	entityMap := make(map[string]interface{})

	for id, jsonEntity := range jsonEntityMap {
		if contentType != targetType {
			//fmt.Printf("id=%s\n", id)
			entityId, err := identifier.New(unypedLocalDomain, []byte(id), nil)
			//identifier.Parse(id)

			if err != nil {
				return nil, err
			}

			cTransRes, cTransErr := translators.Translate(context, contentType,
				entityId, jsonEntity, targetType)

			if cTransRes == nil || cTransErr == nil {
				return nil, fmt.Errorf("fileResolver.loadMap Failed: translators.Translate channels are undefined")
			}

			hasTranslation := false
			for !hasTranslation {
				select {
				case entity, ok := <-cTransRes:
					if ok {
						entityMap[string(entityId.Value())] = entity
						hasTranslation = true
					}
				case err, ok := <-cTransErr:
					if ok {
						return nil, err
					}
				}
			}

		} else {
			jsonEntity.(map[string]interface{})["id"] = id
			entityMap[id] = jsonEntity
		}
	}

	return entityMap, nil
}
