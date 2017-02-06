package domainScope

import (
	"strings"

	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/domain"
	"github.com/distributed-vision/go-resources/ids/domainScopeFormat"
	"github.com/distributed-vision/go-resources/ids/domainScopeVisibility"
	"github.com/distributed-vision/go-resources/resolvers"
	"github.com/distributed-vision/go-resources/version/versionType"
)

type scope struct {
	ids.Domain
	visibility ids.DomainScopeVisibility
	format     ids.DomainScopeFormat
}

func unmarshalJSON(json map[string]interface{}, opts map[string]interface{}) (interface{}, error) {
	id, err := base62.Decode(json["id"].(string))

	if err != nil {
		return nil, err
	}

	visibility, err := domainScopeVisibility.Parse(json["visibility"].(string))

	if err != nil {
		return nil, err
	}

	format, err := domainScopeFormat.Parse(json["format"].(string))

	if err != nil {
		return nil, err
	}

	return NewDomainScope(
		id,
		toString(json, "name"),
		toString(json, "description"),
		visibility,
		format,
		toDomainInfo(json["domainInfo"].(map[string]interface{}), opts))
}

func NewDomainScope(id []byte, name string, description string, visibility ids.DomainScopeVisibility, format ids.DomainScopeFormat, info map[interface{}]interface{}) (ids.DomainScope, error) {

	info["name"] = name
	info["description"] = description
	//fmt.Printf("info:%+v", info)
	base, err := domain.NewDomain(nil, id, nil, 0, versionType.UNVERSIONED, info)

	if err != nil {
		return nil, err
	}

	return &scope{
		Domain:     base,
		visibility: visibility,
		format:     format}, nil
}

func (this *scope) Visibility() ids.DomainScopeVisibility {
	return this.visibility
}

func (this *scope) Format() ids.DomainScopeFormat {
	return this.format
}

func toString(json map[string]interface{}, field string) string {
	value := json[field]

	if value == nil {
		return ""
	}

	return value.(string)
}

func toDomainInfo(infoIn map[string]interface{}, opts map[string]interface{}) map[interface{}]interface{} {

	infoOut := make(map[interface{}]interface{})

	//var sourcePath string
	var hasSourcePath bool

	if opts != nil {
		/*sourcePath*/ _, hasSourcePath = opts["path"].(string)
	}

	for key, value := range infoIn {

		if strings.ToUpper(key) == strings.ToUpper("resoverInfo") {
			if hasSourcePath {
				resolverInfos, ok := value.([]resolvers.ResolverInfo)
				if ok {
					resolverInfosOut := make([]resolvers.ResolverInfo, len(resolverInfos))

					for index, resolverInfo := range resolverInfos { /*
							if fileResolver.ResolverType().IsAssignableFrom(resolverInfo.ResolverType()) {
								locator := resolverInfo.Value("locator").(string)
								if !path.IsAbs(locator) {
									resolverInfosOut[index] =
										resolverInfo.WithValue("locator", path.Join(path.Dir(sourcePath), locator))
								}
							} else {*/
						resolverInfosOut[index] = resolverInfo
						/*}*/
					}
					infoOut["resolverInfo"] = resolverInfosOut
				}
			} else {
				infoOut[key] = value
			}
		}
	}

	return infoOut
}
