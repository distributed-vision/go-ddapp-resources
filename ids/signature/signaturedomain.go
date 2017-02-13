package signature

import (
	"context"
	"fmt"

	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/domain"
	"github.com/distributed-vision/go-resources/ids/domainScope"
	"github.com/distributed-vision/go-resources/ids/domainType"
	"github.com/distributed-vision/go-resources/ids/identifier"
	"github.com/distributed-vision/go-resources/resolvers"
	"github.com/distributed-vision/go-resources/version/versionType"
)

func init() {
	domain.RegisterJSONUnmarshaller(domainType.SIGNATURE, unmarshalJSON)
}

type signatureDomain struct {
	ids.IdentityDomain
}

func unmarshalJSON(unmarshalContext context.Context, json map[string]interface{}) (ids.Domain, error) {
	rootId, err := base62.Decode(json["id"].(string))
	resolverInfo, hasResolverInfo := unmarshalContext.Value("resolverInfo").(resolvers.ResolverInfo)

	if err != nil {
		return nil, err
	}

	var incarnation *uint32
	var crcLength uint
	var versionType = versionType.UNVERSIONED

	info := make(map[interface{}]interface{})

	for key, value := range json {
		info[key] = value
	}

	var scope ids.DomainScope

	if hasResolverInfo && resolverInfo.Value("scopeId") != nil {
		//fmt.Printf("scopeid=%v\n", resolverInfo.Value("scopeId"))
		scope, err = domainScope.Get(unmarshalContext, domainScope.Selector{Id: resolverInfo.Value("scopeId").([]byte)})
		//fmt.Printf("scope=%v\n", scope)
		if err != nil {
			return nil, err
		}
	}

	if scope == nil {
		return nil, fmt.Errorf("Can't unmarshal domain for unknow scope")
	}

	base, err := identifier.NewIdentityDomain(scope, rootId, incarnation, crcLength, versionType, info)
	//fmt.Printf("base=%v\n", base)

	if err != nil {
		return nil, err
	}

	return &signatureDomain{base}, nil
}

func NewElements() (ids.SignatureElements, error) {
	return nil, nil
}
