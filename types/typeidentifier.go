package types

import (
	"errors"
	"fmt"

	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/identifier"
	"github.com/distributed-vision/go-resources/version"
)

var typeIds = make(map[string]ids.TypeIdentifier)

func MustNewId(domain interface{}, id []byte, idVersion version.Version) ids.TypeIdentifier {
	typeid, err := NewId(domain, id, idVersion)

	if err != nil {
		panic(fmt.Sprintf("Can't create new type id: %s", err))
	}

	return typeid
}

func NewId(domain interface{}, id []byte, idVersion version.Version) (ids.TypeIdentifier, error) {

	base, err := identifier.New(domain, id, idVersion)

	if err != nil {
		return nil, err
	}

	typeTd := &typeIdentifier{base}

	if registeredId, ok := typeIds[string(typeTd.Value())]; ok {
		return registeredId, nil
	}

	RegisterId(typeTd)
	return typeTd, nil
}

type typeIdentifier struct {
	ids.Identifier
}

func (this *typeIdentifier) Elements() (ids.SignatureElements, error) {
	return nil, errors.New("TODO")
}

func (this *typeIdentifier) IsAssignableFrom(typeId ids.TypeIdentifier) bool {
	return false
}

func RegisterId(typeid ids.TypeIdentifier) {
	typeIds[string(typeid.Value())] = typeid
}
