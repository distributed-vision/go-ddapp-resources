package identifiers

import (
	"bytes"
	"errors"
	"hash/crc32"
	"time"

	"github.com/distributed-vision/go-resources/encoding"
	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/resolvers/domainScopeResolver"
	"github.com/distributed-vision/go-resources/resolvers/mappingResolver"
	"github.com/howeyc/crc16"
	"github.com/sigurn/crc8"
)

type identifier struct {
	value []byte
}

func NewIdentifier(domain ids.Domain, id []byte) (ids.Identifier, error) {

	if domain == nil {
		return nil, errors.New("Invalid domain: undefined")
	}

	if id == nil {
		return nil, errors.New("Invalid id: undefined")
	}

	crcLength := domain.CrcLength()
	domainId := domain.Id()

	if domainId == nil {
		return nil, errors.New("Invalid domain: id undefined")
	}

	value := bytes.Join([][]byte{domainId, id}, []byte{})

	if crcLength > 0 {
		crc, err := calculateCrc(value, crcLength)

		if err != nil {
			return nil, err
		}

		value = bytes.Join([][]byte{value, crc}, []byte{})
	}

	return WrapIdentifier(value), nil
}

func WrapIdentifier(id []byte) ids.Identifier {
	return &identifier{id}
}

func UnwrapIdentifier(id ids.Identifier) []byte {
	return id.ValueOf()
}

func AsLocator(id ids.Identifier) ids.Locator {
	locator, ok := id.(*locator)

	if ok {
		return locator
	}

	return getLocator(id)
}

func (this *identifier) Id() []byte {
	scopeLength := ScopeLength(this.value)
	domainLength := domainLength(this.value) + 1
	identifierLength, _ := identifierLength(this.value)
	return this.value[scopeLength+domainLength : scopeLength+domainLength+identifierLength]
}

func (this *identifier) ScopeId() []byte {
	id := this.value[:ScopeLength(this.value)]
	id[0] = id[0] & 0x3f
	return id
}

func (this *identifier) DomainId() []byte {
	return this.value[:ScopeLength(this.value)+domainLength(this.value)+1]
}

func (this *identifier) DomainIdRoot() []byte {
	return this.value[ScopeLength(this.value)+1 : domainLength(this.value)-incarnationLength(this.value)+2]
}

func (this *identifier) DomainIncarnation() *uint32 {
	return incarnationValue(this.value)
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
	selector := domainScopeResolver.Selector{ScopeId: this.ScopeId()}
	cres, cerr := domainScopeResolver.Resolve(selector)

	select {
	case scope = <-cres:
		break
	case err := <-cerr:
		ignore(err)
		break
	}
	return scope
}

func (this *identifier) Domain() ids.IdentityDomain {
	scope, _ := this.Scope().Get(this.DomainId())
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

	return bytes.Equal(this.value, UnwrapIdentifier(other))
}

func (this *identifier) CompareTo(o ids.Identifier) int {
	result := bytes.Compare(this.DomainId(), o.DomainId())

	if result != 0 {
		return result
	}

	return bytes.Compare(this.Id(), o.Id())
}

func (this *identifier) ToEncodedString(domainEncoding string, idEncoding string, seperator string) string {

	if seperator != "" {
		// this effectivelt removes the CRC from the Identifier
		// so format the domain as a non crc domain
		did := this.DomainId()
		did = bytes.Join([][]byte{[]byte{did[0] & 0x3f}, did[1:]}, []byte{})
		dom, _ := encoding.Encode(did, domainEncoding)
		id, _ := encoding.Encode(this.Id(), idEncoding)
		return dom + seperator + id
	}

	result, _ := encoding.Encode(this.value, domainEncoding)
	return result
}

func (this *identifier) ToString() string {
	return this.ToEncodedString("base62", "", "")
}

func (this *identifier) ValueOf() []byte {
	return this.value
}

func (this *identifier) Sign(signatureDomain ids.SignatureDomain) (ids.Signature, error) {
	return nil, nil
}

func (this *identifier) IsFor(typeId ids.TypeIdentifier) bool {
	return this.Domain().IsFor(typeId)
}

func (this *identifier) As(domain ids.IdentityDomain) (id ids.Identifier, err error) {
	ci, ce := mappingResolver.ResolveMapping(this, domain, time.Now())

	select {
	case id = <-ci:
		break
	case err = <-ce:
		break
	}

	return id, err
}

func (this *identifier) MapFrom(from ids.Identifier, after *time.Time, before *time.Time) {
	mappingResolver.AddMapping(from, this, after, before)
}

func (this *identifier) MapTo(to ids.Identifier, after *time.Time, before *time.Time) {
	mappingResolver.AddMapping(this, to, after, before)
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

func IncarnationAsBytes(incarnation *uint32) []byte {
	var incarnationSlice []byte

	if incarnation == nil {
		incarnationSlice = []byte{}
	} else if *incarnation < 0xff {
		incarnationSlice = []byte{byte(*incarnation & 0xff)}
	} else if *incarnation < 0xffff {
		buf := [2]byte{0, 0}
		incarnationSlice = htons(buf[:], 0, uint16(*incarnation&0xffff))
	} else {
		buf := [4]byte{0, 0, 0, 0}
		incarnationSlice = htonl(buf[:], 0, *incarnation)
	}

	return incarnationSlice
}

func ScopeLength(value []byte) uint {
	switch value[0] & 0x3f {
	// TODO handle named scopes
	default:
		return 1
	}
}

func domainLength(value []byte) uint {
	return uint(value[ScopeLength(value)] & 0x3f)
}

func incarnationLength(value []byte) uint {
	return incarnationBitsLength(value[ScopeLength(value)])
}

func incarnationValue(value []byte) *uint32 {
	incLen := incarnationLength(value)

	if incLen > 0 {
		incOffset := ScopeLength(value) +
			domainLength(value) - incLen + 1

		if incLen == 1 {
			res := uint32(value[incOffset])
			return &res
		}
		if incLen == 2 {
			res := uint32(ntohs(value, int(incOffset)))
			return &res
		}
		res := ntohl(value, int(incOffset))
		return &res
	}

	return nil
}

func identifierCrcLength(value []byte) (uint, error) {
	crcLength, err := crcLength(value[0])
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

	length := uint(len(value)) - ScopeLength(value) -
		(domainLength(value) + 1) - identifierCrcLength
	if length < 0 {
		return 0, nil
	}
	return length, nil
}

func calculateCrc(value []byte, crcLength uint) ([]byte, error) {
	return crcCalc(value, crcLength)
}

func isValid(value []byte) bool {
	scopeLength := ScopeLength(value)
	domainLength := domainLength(value)
	incarnationLength := incarnationLength(value)
	crcLength, err := identifierCrcLength(value)

	if err != nil {
		return false
	}

	crcCalc, err := crcCalc(value[:uint(len(value))-crcLength], crcLength*8)

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

	locatorKey := string(id.ValueOf())
	locator, ok := locators[locatorKey]

	if ok {
		return locator
	}

	locator = NewLocator(id)
	locators[locatorKey] = locator
	return locator
}

func ntohl(buffer []byte, index int) uint32 {
	return (uint32(0xff&buffer[index]) << 24) |
		(uint32(0xff&buffer[index+1]) << 16) |
		(uint32(0xff&buffer[index+2]) << 8) |
		uint32(0xff&buffer[index+3])
}

func ntohs(buffer []byte, index int) uint16 {
	return uint16(0xff&buffer[index])<<8 |
		uint16(0xff&buffer[index+1])
}

func htonl(buffer []byte, index int, value uint32) []byte {
	buffer[index] = byte(0xff & (value >> 24))
	buffer[index+1] = byte(0xff & (value >> 16))
	buffer[index+2] = byte(0xff & (value >> 8))
	buffer[index+3] = byte(0xff & (value))
	return buffer
}

func htons(buffer []byte, index int, value uint16) []byte {
	buffer[index] = byte(0xff & (value >> 8))
	buffer[index+1] = byte(0xff & (value))
	return buffer
}

func IncarnationLengthBits(incarnation []byte) (byte, error) {
	switch len(incarnation) {
	case 0:
		return 0, nil
	case 1:
		return 1 << 6, nil
	case 2:
		return 1 << 6, nil
	case 4:
		return 3 << 6, nil
	default:
		return 0, errors.New("Invalid incarnation length")
	}
}

func incarnationBitsLength(bits byte) uint {
	bits = ((bits & 0xff) >> 6)
	if bits < 3 {
		return uint(bits)
	}
	return 4
}

func CrcLengthBits(crcLength uint) (byte, error) {
	switch crcLength {
	case 0:
		return 0, nil
	case 8:
		return 1 << 6, nil
	case 16:
		return 2 << 6, nil
	case 32:
		return 3 << 6, nil
	default:
		return 0, errors.New("Invalid crc length")
	}
}

func crcLength(bits byte) (int, error) {
	bits = ((bits & 0xff) >> 6)
	switch bits {
	case 0:
		return 0, nil
	case 1:
		return 8, nil
	case 2:
		return 16, nil
	case 3:
		return 32, nil
	default:
		return 0, errors.New("Invalid crc length")
	}
}

var crc8Table *crc8.Table = crc8.MakeTable(crc8.CRC8_MAXIM)
var crc16Table *crc16.Table = crc16.MakeTable(crc16.IBM)
var crc32Table *crc32.Table = crc32.MakeTable(crc32.IEEE)

func crcCalc(value []byte, crcLength uint) ([]byte, error) {
	switch crcLength {
	case 0:
		return make([]byte, 0), nil
	case 8:
		buf := [1]byte{crc8.Checksum(value, crc8Table)}
		return buf[:], nil
	case 16:
		buf := make([]byte, 2)
		return htons(buf, 0, crc16.Checksum(value, crc16Table)), nil
	case 32:
		buf := make([]byte, 4)
		return htonl(buf, 0, crc32.Checksum(value, crc32Table)), nil
	default:
		return nil, errors.New("Invalid crc length")
	}
}

func ignore(err error) {

}
