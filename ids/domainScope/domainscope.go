package domainScope

import (
	"context"
	"fmt"
	"strings"

	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/encoding/encoderType"
	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/domain"
	"github.com/distributed-vision/go-resources/ids/domainScopeFormat"
	"github.com/distributed-vision/go-resources/ids/domainScopeVisibility"
	"github.com/distributed-vision/go-resources/resolvers"
	"github.com/distributed-vision/go-resources/version"
	"github.com/distributed-vision/go-resources/version/versionType"
)

var DomainPathKey = string(domain.MustDecodeId(encoderType.BASE62, "2", "1")) + "-DOMAINPATH"

type scope struct {
	ids.Domain
	visibility ids.DomainScopeVisibility
	format     ids.DomainScopeFormat
}

func unmarshalJSON(unmarshalContext context.Context, json map[string]interface{}) (interface{}, error) {
	id, err := base62.Decode(json["id"].(string))

	if err != nil {
		return nil, err
	}

	if len(id) == 0 {
		return nil, fmt.Errorf("Empty id is invalid")
	}

	visibility, err := domainScopeVisibility.Parse(json["visibility"].(string))

	if err != nil {
		return nil, err
	}

	format, err := domainScopeFormat.Parse(json["format"].(string))

	if err != nil {
		return nil, err
	}

	/*
		if unmarshalContext.Value("paths")!=nil {
			unmarshalContext=context.WithValue(unmarshalContext, domain.DomainPathKey, path )
		}
	*/

	return NewDomainScope(
		id,
		toString(json, "name"),
		toString(json, "description"),
		visibility,
		format,
		toDomainInfo(unmarshalContext, id, json["domainInfo"].(map[string]interface{})))
}

func NewDomainScope(id []byte, name string, description string, visibility ids.DomainScopeVisibility, format ids.DomainScopeFormat, info map[interface{}]interface{}) (ids.DomainScope, error) {

	info["name"] = name
	info["description"] = description

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

func (this *scope) RegisterResolvers() error {
	resolverInfos := this.InfoValue("resolverInfo").([]resolvers.ResolverInfo)
	var errs = make([]error, 0)

	if resolverInfos != nil {

		for _, resolverInfo := range resolverInfos {
			factory, err := resolvers.NewResolverFactory(resolverInfo)

			if err == nil {
				resolvers.RegisterFactory(factory)
			} else {
				//fmt.Printf("factory err=%s\n", err)
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("Resolver registration fails: %v", errs)
	}

	return nil
}

func toString(json map[string]interface{}, field string) string {
	value := json[field]

	if value == nil {
		return ""
	}

	return value.(string)
}

var fileResolverDomainId = domain.MustDecodeId(encoderType.BASE62, "T", "0", uint32(0), uint(0), versionType.SEMANTIC)

func toDomainInfo(infoContext context.Context, scopeId []byte, infoIn map[string]interface{}) map[interface{}]interface{} {

	infoOut := make(map[interface{}]interface{})

	var sourcePaths []string

	if contextResInfo, ok := infoContext.Value("resolverInfo").(resolvers.ResolverInfo); ok {
		if paths, ok := contextResInfo.Value("paths").([]string); ok {
			sourcePaths = paths
		}
	}

	for key, value := range infoIn {

		if strings.ToUpper(key) == strings.ToUpper("resolverInfo") {
			resolverInfos, ok := value.([]interface{})
			if ok {
				resolverInfosOut := make([]resolvers.ResolverInfo, len(resolverInfos))

				for index, resolverInfo := range resolverInfos {
					var infoMap = map[interface{}]interface{}{}
					var resolverType ids.TypeIdentifier

					for key, value := range resolverInfo.(map[string]interface{}) {

						switch key {
						case "resolverType":
							resolverTypeParts := strings.Split(value.(string), "-")
							idPart := []byte(resolverTypeParts[0])

							if len(resolverTypeParts) > 1 {
								versionPart, err := version.Parse(resolverTypeParts[1])

								if err == nil {
									resolverType, err = ids.NewTypeId(fileResolverDomainId, idPart, versionPart)
								}
							} else {
								var err error
								resolverType, err = ids.NewTypeId(fileResolverDomainId, idPart, nil)

								if err != nil {
								}
							}
							break

						default:
							infoMap[key] = value
						}
					}

					if resolverType != nil {

						if sourcePaths != nil && "FileResolver" == string(resolverType.Id()) {
							if /*value*/ _, ok := infoMap["paths"]; !ok {
								infoMap["paths"] = sourcePaths
							} else {
								// TODO append source paths after predefined paths
							}
						}

						infoMap["scopeId"] = scopeId
						resolverInfosOut[index] = resolvers.NewResolverInfo(resolverType,
							[]ids.TypeIdentifier{domainEntityType}, infoMap)
					}
				}

				infoOut["resolverInfo"] = resolverInfosOut
			}
		} else {
			infoOut[key] = value
		}
	}

	return infoOut
}
