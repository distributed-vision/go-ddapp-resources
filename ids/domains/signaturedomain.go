package domains

import "github.com/distributed-vision/go-resources/ids"

func NewElements() (ids.SignatureElements, error) {
	return nil, nil
}

type signatureDomain struct {
	ids.IdentityDomain
}

func (this *signatureDomain) CreateDomain(ids.SignatureElements) (ids.Signature, error) {
	return nil, nil
}
