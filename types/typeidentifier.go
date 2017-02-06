package types

import "github.com/distributed-vision/go-resources/ids"

var typeIds = make(map[string]ids.TypeIdentifier)

func IdOf(identifierSpecification interface{}) ids.TypeIdentifier {
	return nil
}

func RegisterId(typeid ids.TypeIdentifier) {
	typeIds[string(typeid.Value())] = typeid
}
