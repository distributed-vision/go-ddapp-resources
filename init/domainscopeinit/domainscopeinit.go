package domainscopeinit

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/resolvers"
	"github.com/distributed-vision/go-resources/resolvers/fileresolver"
	"github.com/distributed-vision/go-resources/types/gotypeid"
)

func Init() {
	scopePath := os.Getenv("DV_DOMAIN_SCOPE_PATH")

	if scopePath == "" {
		scopePath = "../../../domain-scope"
	}

	entityType := gotypeid.IdOf(reflect.TypeOf((*ids.DomainScope)(nil)).Elem())

	resolverFactory, err := fileresolver.NewResolverFactory(
		resolvers.NewResolverInfo(fileresolver.PublicType, []ids.TypeIdentifier{entityType},
			map[interface{}]interface{}{
				"location": "scopeinfo.json",
				"paths":    filepath.SplitList(scopePath)}))

	if err != nil {
		panic(fmt.Sprintf("Unexpected error creating resolver factory: %s", err))
	}

	resolvers.RegisterFactory(resolverFactory)
}
