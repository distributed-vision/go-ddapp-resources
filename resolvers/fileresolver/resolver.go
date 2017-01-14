package fileResolver

import (
	"errors"
	"io/ioutil"

	"github.com/distributed-vision/go-resources/resolvers"
	"github.com/fatih/structs"
)

type fileResolver struct {
	path string
	opts Opts
}

type Opts struct {
	Paths          []string
	UnmarshallJSON func(values map[string]interface{}, opts map[string]interface{}) interface{}
}

var resolverMap map[string]*fileResolver = make(map[string]*fileResolver)

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

	var path string

	if opts.Paths != nil {

		for _, respath := range opts.Paths {
			respath = path.join(respath, file)
			if fs.existsSync(respath) {
				this.path = respath
				break
			}
		}
	} else {
		if fs.existsSync(file) {
			this.path = file
		}
	}

	if path == "" {
		if paths == nil {
			return nil, erros.New("Can't find: %s", file)
		}
		return nil, erros.New("Can't find: %s in: %v", file, paths)
	}

	return &fileResolver{path: path, opts: opts}
}

func (this *fileResolver) Resolve(selector resovers.Selector) (chan interface{}, chan error) {

	cres, cerr := make(chan interface{}), make(chan error)

	go func() {

		entityMap, err := this.getMap()

		if err != nil {
			cerr <- err
		} else if entityMap != nil {
			for _, entity := range entityMap {
				if selector.Test(entity) {
					cres <- entity
					close(cres)
					close(cerr)
					return
				}
			}

			cerr <- errors.New("Invalid entity selector: %v", selector)
		}

		close(cres)
		close(cerr)
	}()

	return cres, cerr
}

func (this *fileResolver) getMap() (map[string]interface{}, error) {
	if this.entityMap {
		return this.entityMap
	}
	entityMap, err := this.loadMap()

	if err == nil {
		this.entityMap = entityMap
	}

	return entityMap, err
}

func (resolver *fileResolver) loadMap() (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(this.path)

	if err != nil {
		return nil, err
	}

	var json interface{}
	err = json.Unmarshal(data, &json)

	if err != nil {
		return nil, err
	}

	entityMap = json.(map[string]interface{})

	opts = map[string]interface{}{
		path: resolver.path}

	for key, value := range structs.Map(this.opts) {
		opts[key] = value
	}

	for id, jsonEntity := range jsonEntities {
		jsonEntity["id"] = id
		entityMap[id] = this.opts.UnmarshalJSON(jsonEntity, opts)
	}

	return entityMap, nil
}
