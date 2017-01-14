package domains

import (
	"errors"

	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/domainType"
)

func FromJSON(json map[string]interface{}, opts map[string]interface{}) (ids.Domain, error) {
	dt, err := domainType.Parse(json["domainType"].(string))

	if err != nil {
		return nil, err
	}

	switch dt {
	case domainType.SCOPE:
		//return NewDomainScopeFromJSON(json, opts), nil
	case domainType.IDENTITY:
		//return _identityDomainFromJSON(json, opts), nil
	case domainType.SIGNATURE:
		//return NewSignatureDomainFromJSON(json, opts), nil
	case domainType.SEQUENCE:
		//return NewSequenceDomainFromJSON(json, opts), nil

	}

	return nil, errors.New("Unknown domain type: " + json["domainType"].(string))
}
