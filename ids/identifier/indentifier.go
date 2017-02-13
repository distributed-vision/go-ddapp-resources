package identifier

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/distributed-vision/go-resources/encoding"
	"github.com/distributed-vision/go-resources/encoding/encoderType"
	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/domain"
	"github.com/distributed-vision/go-resources/ids/domainScope"
	"github.com/distributed-vision/go-resources/ids/mappings"
	"github.com/distributed-vision/go-resources/version"
	"github.com/distributed-vision/go-resources/version/versionType"
)

func Init() {
}

type identifier struct {
	value []byte
}

func New(domainValue interface{}, id []byte, idVersion version.Version) (ids.Identifier, error) {

	var crcLength uint
	var domainId []byte
	var idVersionType versionType.VersionType
	var err error

	if domainValue == nil {
		return nil, errors.New("Invalid domain: undefined")
	} else {
		switch t := domainValue.(type) {
		case ids.Domain:
			dom := domainValue.(ids.Domain)
			crcLength = dom.CrcLength()
			domainId = dom.Id()
			idVersionType = dom.VersionType()
		case []byte:
			domainId = domainValue.([]byte)
			crcLength, err = identifierCrcLength(domainId)
			idVersionType, err = identifierVersionType(domainId)

			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("Invalid domain type: %s", t)
		}
	}

	if id == nil {
		return nil, errors.New("Invalid id: undefined")
	}

	if domainId == nil {
		return nil, errors.New("Invalid domain: id undefined")
	}

	var value []byte

	switch idVersionType {
	case versionType.UNVERSIONED:
		value = bytes.Join([][]byte{domainId, id}, []byte{})
		break
	case versionType.NUMERIC:
		if nver, ok := idVersion.(version.NumericVersion); !ok {
			return nil, errors.New("Expected numeric version")
		} else {
			value = bytes.Join([][]byte{domainId, []byte{nver.ByteLength()}, id, nver.Bytes()}, []byte{})
		}
		break
	case versionType.SEMANTIC:
		if sver, ok := idVersion.(*version.SemanticVersion); !ok {
			return nil, errors.New("Expected semantic version")
		} else {
			verbytes := sver.Bytes()
			blen := len(verbytes)

			if blen > 255 {
				return nil, errors.New("Version length > 255")
			}
			value = bytes.Join([][]byte{domainId, []byte{byte(blen & 0xff)}, id, verbytes}, []byte{})
		}
		break
	default:
		return nil, errors.New("Unknown version type")
	}

	if crcLength > 0 {
		crc, err := calculateCrc(value, crcLength)

		if err != nil {
			return nil, err
		}

		value = bytes.Join([][]byte{value, crc}, []byte{})
	}

	return Wrap(value), nil
}

func Wrap(id []byte) ids.Identifier {
	return &identifier{id}
}

func Unwrap(id ids.Identifier) []byte {
	return id.Value()
}

func Parse(id string) (ids.Identifier, error) {
	return nil, nil
}

func AsLocator(id ids.Identifier) ids.Locator {
	locator, ok := id.(*locator)

	if ok {
		return locator
	}

	return getLocator(id)
}

func (this *identifier) Id() []byte {
	domainOffset := domain.DomainOffset(this.value)
	domainLength := domain.DomainLength(this.value) + domain.VersionLengthLength(this.value)
	versionLength := domain.VersionLength(this.value)
	identifierLength, _ := identifierLength(this.value)
	//fmt.Printf("do=%v, dl=%v, vl=%v, il=%v\n", domainOffset, domainLength, versionLength, identifierLength)
	return this.value[domainOffset+domainLength : domainOffset+domainLength+identifierLength-versionLength]
}

func (this *identifier) ScopeId() []byte {
	id := this.value[:domain.ScopeLength(this.value)]
	return id
}

func (this *identifier) DomainId() []byte {
	return this.value[:domain.ScopeLength(this.value)+domain.DomainLength(this.value)+1]
}

func (this *identifier) HasVersion() bool {
	return domain.VersionLengthLength(this.value) > 0
}

func (this *identifier) VersionId() []byte {
	versionLength := domain.VersionLength(this.value)
	if versionLength == 0 {
		return nil
	}
	domainOffset := domain.DomainOffset(this.value)
	domainLength := domain.DomainLength(this.value) + domain.VersionLengthLength(this.value)
	identifierLength, _ := identifierLength(this.value)
	return this.value[domainOffset+domainLength+identifierLength-versionLength:]
}

func (this *identifier) Version() version.Version {
	versionId := this.VersionId()

	if versionId == nil {
		return nil
	}

	vtype, err := domain.VersionTypeValue(this.value)

	if err == nil {
		switch vtype {
		case versionType.NUMERIC:
			return version.NumericVersion(domain.NumericVersionValue(versionId))
		case versionType.SEMANTIC:
			result, err := version.Parse(string(versionId))
			if err == nil {
				return result
			}
		}
	}

	return nil
}

func (this *identifier) DomainIdRoot() []byte {
	return domain.IdRootValue(this.value)
}

func (this *identifier) DomainIncarnation() *uint32 {
	return domain.IncarnationValue(this.value)
}

func (this *identifier) Checksum() []byte {
	crcLength, _ := identifierCrcLength(this.value)
	return this.value[uint(len(this.value))-crcLength:]
}

func (this *identifier) sign(signatureDomain ids.SignatureDomain) (ids.Signature, error) {
	crcLength, _ := identifierCrcLength(this.value)
	signatureBytes := this.value[:uint(len(this.value))-crcLength]
	signatureBytes[0] = signatureBytes[0] & 0x3f
	return signatureDomain.CreateSignature(struct {
		domainId       []byte
		id             []byte
		signatureBytes []byte
	}{domainId: this.DomainId(),
		id:             this.Id(),
		signatureBytes: signatureBytes})
}

func (this *identifier) Scope() (scope ids.DomainScope) {
	selector := domainScope.Selector{Id: this.ScopeId()}
	scope, _ = domainScope.Get(context.Background(), selector)

	return scope
}

func (this *identifier) Domain() ids.IdentityDomain {
	scope, _ := domain.Get(context.Background(), domain.Selector{Scope: this.Scope(), Id: this.DomainId()})
	return scope
}

func (this *identifier) IsUndefined() bool {
	return this.DomainId() == nil
}

func (this *identifier) IsNull() bool {
	return this.value == nil
}

func (this *identifier) IsValid() bool {
	return isValid(this.value)
}

func (this *identifier) Matches(other ids.Identifier) bool {

	if this.Equals(other) {
		return true
	}

	if bytes.Equal(this.DomainId(), other.DomainId()) {
		return false
	}

	as, err := other.As(this.Domain())
	if err != nil {
		return false
	}

	return this.Equals(as)
}

func (this *identifier) Equals(other ids.Identifier) bool {
	if other == nil {
		return false
	}

	return bytes.Equal(this.value, Unwrap(other))
}

func (this *identifier) CompareTo(o ids.Identifier) int {
	result := bytes.Compare(this.DomainId(), o.DomainId())

	if result != 0 {
		return result
	}

	return bytes.Compare(this.Id(), o.Id())
}

func (this *identifier) Encode(seperator string, encoders ...encoderType.EncoderType) string {
	domainEncoder := encoderType.BASE62
	idEncoder := encoderType.BASE62
	versionEncoder := encoderType.BASE62

	if len(encoders) > 0 {
		domainEncoder = encoders[0]

		if len(encoders) > 1 {
			idEncoder = encoders[1]
		} else {
			idEncoder = domainEncoder
		}

		if len(encoders) > 2 {
			versionEncoder = encoders[2]
		} else {
			versionEncoder = domainEncoder
		}
	}

	if seperator != "" {
		// this effectivelt removes the CRC from the Identifier
		// so format the domain as a non crc domain
		did := this.DomainId()
		did = bytes.Join([][]byte{[]byte{did[0] & 0x3f}, did[1:]}, []byte{})
		dom, _ := encoding.Encode(did, domainEncoder)
		id, _ := encoding.Encode(this.Id(), idEncoder)

		result := dom + seperator + id

		if this.HasVersion() {
			version, _ := encoding.Encode(this.VersionId(), versionEncoder)
			result = result + seperator + version
		}

		return result
	}

	result, _ := encoding.Encode(this.value, domainEncoder)
	return result
}

func (this *identifier) String() string {
	return this.Encode("", encoderType.BASE62)
}

func (this *identifier) Value() []byte {
	return this.value
}

func (this *identifier) Sign(signatureDomain ids.SignatureDomain) (ids.Signature, error) {
	return nil, nil
}

func (this *identifier) IsFor(typeId ids.TypeIdentifier) bool {
	return this.Domain().IsFor(typeId)
}

func (this *identifier) As(domain ids.IdentityDomain) (id ids.Identifier, err error) {
	ci, ce := mappings.Resolve(context.Background(), mappings.Selector{From: this, To: domain, At: time.Now()})

	select {
	case id = <-ci:
		break
	case err = <-ce:
		break
	}

	return id, err
}

func (this *identifier) MapFrom(from ids.Identifier, after *time.Time, before *time.Time) {
	mappings.Add(from, this, after, before)
}

func (this *identifier) MapTo(to ids.Identifier, after *time.Time, before *time.Time) {
	mappings.Add(this, to, after, before)
}

func (this *identifier) MapBetween(id ids.Identifier, after *time.Time, before *time.Time) {
	this.MapFrom(id, after, before)
	this.MapTo(id, after, before)
}

func (this *identifier) Get() (interface{}, error) {
	return AsLocator(this).Get()
}

func (this *identifier) GetAs(typeId ids.TypeIdentifier) (interface{}, error) {
	return AsLocator(this).GetAs(typeId)
}

func (this *identifier) Resolve() (chan interface{}, chan error) {
	return AsLocator(this).Resolve()
}

func (this *identifier) ResolveAs(typeId ids.TypeIdentifier) (chan interface{}, chan error) {
	return AsLocator(this).ResolveAs(typeId)
}

func (this *identifier) TypeId() ids.TypeIdentifier {
	return this.Domain().TypeId()
}

func identifierVersionType(value []byte) (versionType.VersionType, error) {
	return domain.VersionTypeValue(value)
}

func identifierCrcLength(value []byte) (uint, error) {
	crcLength, err := domain.CrcLengthValue(value)
	if err != nil {
		return 0, err
	}
	return uint(crcLength / 8), nil
}

func identifierLength(value []byte) (uint, error) {
	identifierCrcLength, err := identifierCrcLength(value)

	if err != nil {
		return 0, err
	}

	length := uint(len(value)) - domain.DomainOffset(value) -
		domain.DomainLength(value) - domain.VersionLengthLength(value) -
		identifierCrcLength

	if length < 0 {
		return 0, nil
	}
	return length, nil
}

func calculateCrc(value []byte, crcLength uint) ([]byte, error) {
	return domain.CrcCalc(value, crcLength)
}

func isValid(value []byte) bool {
	scopeLength := domain.ScopeLength(value)
	domainLength := domain.DomainLength(value)
	incarnationLength := domain.IncarnationLength(value)
	crcLength, err := identifierCrcLength(value)

	if err != nil {
		return false
	}

	crcCalc, err := domain.CrcCalc(value[:uint(len(value))-crcLength], crcLength*8)

	if err != nil {
		return false
	}

	return crcLength <= 4 &&
		incarnationLength <= 4 &&
		scopeLength <= uint(len(value))-domainLength-crcLength &&
		domainLength > incarnationLength &&
		domainLength <= uint(len(value))-domainLength-crcLength &&
		bytes.Equal(value[uint(len(value))-crcLength:], crcCalc)
}

// TODO - this should be an LRU cache
var locators = make(map[string]ids.Locator)

func getLocator(id ids.Identifier) ids.Locator {

	locatorKey := string(id.Value())
	locator, ok := locators[locatorKey]

	if ok {
		return locator
	}

	locator = NewLocator(id)
	locators[locatorKey] = locator
	return locator
}

func ignore(err error) {

}
