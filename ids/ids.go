package ids

type Domain interface {
	ToString() string
	ToEncodedString(encoder string) string
	ToJSON() string
	IsFor(typeId TypeIdentifier) bool
	Matches(other Domain) bool
	Id() []byte
	ScopeId() []byte
	Scope() DomainScope
	IdRoot() []byte
	Incarnation() *uint32
	NewIncarnation(incarnation uint32, crcLength uint, info map[string]interface{}) (Domain, error)
	CrcLength() uint
	Name() string
	Description() string
	TypeId() TypeIdentifier
	Source() string
	Equals(other Domain) bool
}

type DomainScopeVisibility int
type DomainScopeFormat int
type DomainType int

type DomainScope interface {
	Domain
	Visibility() DomainScopeVisibility
	Format() DomainScopeFormat
	IdBits() int

	Resolve(domainId []byte) (chan Domain, chan error)
	Get(domainId []byte) (Domain, error)
}

type IdentityDomain interface {
	Domain
}

type SignatureDomain interface {
	IdentityDomain
	CreateSignature(elements interface{}) (Signature, error)
}

type SequenceDomain interface {
	IdentityDomain
}

type Identifier interface {
	Id() []byte
	ScopeId() []byte
	DomainId() []byte
	DomainIdRoot() []byte
	DomainIncarnation() *uint32
	Checksum() []byte
	Scope() DomainScope
	Domain() IdentityDomain
	IsUndefined() bool
	IsNull() bool
	IsValid() bool
	ValueOf() []byte

	Sign(signatureDomain SignatureDomain) (Signature, error)

	Matches(other Identifier) bool
	Equals(other Identifier) bool

	As(domain IdentityDomain) (Identifier, error)
}

type SignatureElements interface {
	SignatureBytes() []byte
	Signature() Signature
}

type Signature interface {
	Identifier
	Elements() (SignatureElements, error)
}

type TypeIdentifier interface {
	Signature
	IsAssignableFrom(typeId TypeIdentifier) (bool, error)
}

type Locator interface {
	Identifier

	Get() (interface{}, error)
	GetAs(typeId TypeIdentifier) (interface{}, error)

	Resolve() (chan interface{}, chan error)
	ResolveAs(typeId TypeIdentifier) (chan interface{}, chan error)
}

type IdGenerator interface {
	GenerateId() ([]byte, error)
}
