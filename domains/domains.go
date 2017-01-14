package domains

import (
	"errors"
	"strings"

	"github.com/distributed-vision/go-resources/ids"
)

func FromJSON(json map[string]interface{}, opts map[string]interface{}) (ids.Domain, error) {
	/*switch toDomainType(json["domainType"].(string)) {
	case SCOPE:
		return NewDomainScopeFromJSON(json, opts), nil
	case IDENTITY:
		return _identityDomainFromJSON(json, opts), nil
	case SIGNATURE:
		return NewSignatureDomainFromJSON(json, opts), nil
	case SEQUENCE:
		return NewSequenceDomainFromJSON(json, opts), nil
	}*/

	return nil, errors.New("Unknown domain type: " + json["domainType"].(string))
}

func toDomainType(domainType string) ids.DomainType {
	switch strings.ToUpper(domainType) {
	case "SCOPE":
		return SCOPE
	case "INDENTITY":
		return IDENTITY
	case "SIGNATURE":
		return SIGNATURE
	case "SEQUENCE":
		return SEQUENCE
	}
}
