package signature

import (
	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/identifier"
	"github.com/distributed-vision/go-resources/resolvers/signatureResolver"
)

type signature struct {
	ids.Identifier
	elements ids.SignatureElements
}

func NewSignature(domain ids.SignatureDomain, id []byte) (ids.Signature, error) {
	base, err := identifier.NewIdentifier(domain, id)

	if err != nil {
		return nil, err
	}

	return &signature{base, nil}, nil
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
