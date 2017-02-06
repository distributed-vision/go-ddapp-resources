package ids

import (
	"reflect"
	"sync"

	"github.com/distributed-vision/go-resources/encoding/encoderType"
	"github.com/distributed-vision/go-resources/version"
	"github.com/distributed-vision/go-resources/version/versionType"
)

type Domain interface {
	String() string
	Encode(encoder encoderType.EncoderType) string
	ToJSON() string
	IsFor(typeId TypeIdentifier) bool
	Matches(other Domain) bool
	Id() []byte
	ScopeId() []byte
	Scope() DomainScope
	IdRoot() []byte
	Incarnation() *uint32
	NewIncarnation(incarnation uint32, crcLength uint, info map[interface{}]interface{}) (Domain, error)
	CrcLength() uint
	VersionType() versionType.VersionType
	Name() string
	Description() string
	TypeId() TypeIdentifier
	Source() string
	InfoValue(interface{}) interface{}
	Equals(other Domain) bool
}

type DomainScopeVisibility int

type DomainScopeFormat int
type DomainType int

type DomainScope interface {
	Domain
	Visibility() DomainScopeVisibility
	Format() DomainScopeFormat
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
	VersionId() []byte
	DomainIdRoot() []byte
	DomainIncarnation() *uint32
	Checksum() []byte
	Scope() DomainScope
	Domain() IdentityDomain
	Version() version.Version
	HasVersion() bool
	IsUndefined() bool
	IsNull() bool
	IsValid() bool
	Value() []byte

	Sign(signatureDomain SignatureDomain) (Signature, error)

	Matches(other Identifier) bool
	Equals(other Identifier) bool

	As(domain IdentityDomain) (Identifier, error)
}

var typeInitFunctions = []func(){}
var typeInitMutex = sync.Mutex{}

func OnLocalTypeInit(initFunction func()) {
	typeInitMutex.Lock()
	if IdOfType == nil {
		typeInitFunctions = append(typeInitFunctions, initFunction)
	} else {
		initFunction()
	}
	typeInitMutex.Unlock()
}

func LocalTypeInit(idInitialiser func(gotype reflect.Type) TypeIdentifier) {
	typeInitMutex.Lock()
	if IdOfType == nil {
		IdOfType = idInitialiser
		for _, initFunction := range typeInitFunctions {
			initFunction()
		}
		typeInitFunctions = []func(){}
	}
	typeInitMutex.Unlock()
}

var IdOfType func(gotype reflect.Type) TypeIdentifier

//var NewIdentifier func(domain Domain, id []byte, ver version.Version) (Identifier, error)
//var ParseIdentifier func(id string) (Identifier, error)

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
	IsAssignableFrom(typeId TypeIdentifier) bool
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
