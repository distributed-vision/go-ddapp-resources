package domains

import (
	"path"

	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/domainScopeFormat"
	"github.com/distributed-vision/go-resources/ids/domainScopeVisibility"
	"github.com/distributed-vision/go-resources/resolvers/domainResolver"
)

type domainScope struct {
	*domain
	visibility ids.DomainScopeVisibility
	format     ids.DomainScopeFormat
	idbits     int
}

func NewDomainScopeFromJSON(json map[string]interface{}, opts map[string]interface{}) (ids.DomainScope, error) {
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
		json["name"].(string),
		json["description"].(string),
		visibility,
		format,
		json["idbits"].(int),
		toDomainInfo(json["domainInfo"].(map[string]interface{}), opts))
}

func NewDomainScope(id []byte, name string, description string, visibility ids.DomainScopeVisibility, format ids.DomainScopeFormat, idbits int, info map[string]interface{}) (ids.DomainScope, error) {
	base, err := NewDomain(nil, id, nil, 0, info)

	if err != nil {
		return nil, err
	}

	return &domainScope{
		domain:     base.(*domain),
		visibility: visibility,
		format:     format,
		idbits:     idbits}, nil
}

func (this *domainScope) Visibility() ids.DomainScopeVisibility {
	return this.visibility
}

func (this *domainScope) Format() ids.DomainScopeFormat {
	return this.format
}

func (this *domainScope) IdBits() int {
	return this.idbits
}

func (this *domainScope) Get(domainId []byte) (domain ids.Domain, err error) {
	cd, ce := domainResolver.Resolve(this, domainId)

	select {
	case domain = <-cd:
		break
	case err = <-ce:
		break
	}

	return domain, err
}

func (this *domainScope) Resolve(domainId []byte) (chan ids.Domain, chan error) {
	return domainResolver.Resolve(this, domainId)
}

type location struct {
	resolver string
	locator  string
}

func toDomainInfo(info map[string]interface{}, opts map[string]interface{}) map[string]interface{} {
	var sourcePath string
	var hasSourcePath bool

	if opts != nil {
		sourcePath, hasSourcePath = opts["path"].(string)
	}

	locations, hasLocations := info["locations"]

	if hasSourcePath && hasLocations {
		for _, location := range locations.([]location) {
			if location.resolver == "file" {
				if !path.IsAbs(location.locator) {
					location.locator = path.Join(path.Dir(sourcePath), location.locator)
				}
			}
		}
	}

	return info
}
