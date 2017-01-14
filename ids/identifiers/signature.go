package identifiers

import (
	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/resolvers/signatureResolver"
)

type signature struct {
	*identifier
	elements ids.SignatureElements
}

func NewSignature(domain ids.SignatureDomain, id []byte) (ids.Signature, error) {
	base, err := NewIdentifier(domain, id)

	if err != nil {
		return nil, err
	}

	return &signature{base.(*identifier), nil}, nil
}

func (this *signature) Elements() (ids.SignatureElements, error) {
	if this.elements == nil {
		elements, err := signatureResolver.GetElements(this)

		if err != nil {
			return nil, err
		}

		this.elements = elements
	}

	return this.elements, nil
}
