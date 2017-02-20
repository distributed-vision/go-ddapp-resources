package signaturedomain

import (
	"context"
	"fmt"

	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/domain"
	"github.com/distributed-vision/go-resources/ids/domainscope"
	"github.com/distributed-vision/go-resources/ids/domaintype"
	"github.com/distributed-vision/go-resources/ids/identitydomain"
	"github.com/distributed-vision/go-resources/resolvers"
	"github.com/distributed-vision/go-resources/version/versiontype"
)

func Init() {
}

func init() {
	domain.RegisterJSONUnmarshaller(domaintype.SIGNATURE, unmarshalJSON)
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
	var versionType = versiontype.UNVERSIONED

	info := make(map[interface{}]interface{})

	for key, value := range json {
		info[key] = value
	}

	var scope ids.DomainScope

	if hasResolverInfo && resolverInfo.Value("scopeId") != nil {
		//fmt.Printf("scopeid=%v\n", resolverInfo.Value("scopeId"))
		scope, err = domainscope.Get(unmarshalContext, domainscope.Selector{Id: resolverInfo.Value("scopeId").([]byte)})
		//fmt.Printf("scope=%v\n", scope)
		if err != nil {
			return nil, err
		}
	}

	if scope == nil {
		return nil, fmt.Errorf("Can't unmarshal domain for unknow scope")
	}

	base, err := identitydomain.New(scope, rootId, incarnation, crcLength, versionType, info)
	//fmt.Printf("base=%v\n", base)

	if err != nil {
		return nil, err
	}

	return &signatureDomain{base}, nil
}

func NewElements() (ids.SignatureElements, error) {
	return nil, nil
}
