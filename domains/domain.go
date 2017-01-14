package domains

import (
	"bytes"
	"errors"
	"hash/crc32"

	"github.com/howeyc/crc16"
	"github.com/sigurn/crc8"

	"github.com/distributed-vision/go-resources/encoding"
	"github.com/distributed-vision/go-resources/ids"
)

type domain struct {
	id          []byte
	scope       ids.DomainScope
	root        ids.Domain
	idRoot      []byte
	incarnation *uint32
	crcLength   uint
	typeId      ids.TypeIdentifier
	info        map[string]interface{}
}

func NewDomain(scope ids.DomainScope, idRoot []byte, incarnation *uint32, crcLength uint, info map[string]interface{}) (ids.Domain, error) {

	id, err := toId(scope.Id(), idRoot, incarnation, crcLength)

	if err != nil {
		return nil, err
	}

	return &domain{
		id:          id,
		scope:       scope,
		idRoot:      idRoot,
		incarnation: incarnation,
		crcLength:   crcLength,
		info:        info}, nil
}

func (this *domain) NewIncarnation(incarnation uint32, crcLength uint, info map[string]interface{}) (ids.Domain, error) {

	id, err := toId(this.scope.Id(), this.idRoot, &incarnation, crcLength)

	if err != nil {
		return nil, err
	}

	return &domain{
		id:          id,
		scope:       this.Scope(),
		root:        this,
		idRoot:      this.idRoot,
		incarnation: &incarnation,
		crcLength:   crcLength,
		info:        info}, nil
}

func toId(scopeId []byte, idRoot []byte, incarnation *uint32, crcLength uint) ([]byte, error) {

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

	if len(incarnationSlice) > 4 {
		return nil, errors.New("Incarnation too Long: Incarnation binary length must be <= 4")
	}

	unscoped := bytes.Join([][]byte{idRoot, incarnationSlice}, []byte{})

	if len(unscoped) > 61 {
		return nil, errors.New("Id too Long: domain id unscoped binary length (idRoot+incarnation) must be < 61")
	}

	var unscopedlenSlice []byte

	if len(scopeId) == 0 {
		unscopedlenSlice = []byte{}
	} else {
		incarnationLengthBits, err := incarnationLengthBits(incarnationSlice)
		if err != nil {
			return nil, err
		}
		unscopedlenSlice = []byte{byte(len(unscoped)&0xff) | incarnationLengthBits}
	}

	if crcLength > 0 {
		crcLengthBits, err := crcLengthBits(crcLength)
		if err != nil {
			return nil, err
		}
		scopeId = bytes.Join([][]byte{[]byte{scopeId[0] | crcLengthBits}, scopeId[1:]}, []byte{})
	}

	return bytes.Join([][]byte{scopeId, unscopedlenSlice, unscoped}, []byte{}), nil
}

func (this *domain) Equals(other ids.Domain) bool {

	if this == other.(*domain) {
		return true
	}

	if other == nil {
		return false
	}

	if this.id == nil {
		if other.Id() != nil {
			return false
		}
	} else if !bytes.Equal(this.id, other.Id()) {
		return false
	}

	return true
}

func (this *domain) CompareTo(other ids.Domain) int {
	return bytes.Compare(this.id, other.Id())
}

func (this *domain) SetIfChanged(idRoot []byte, incarnation *uint32) bool {

	changed := false

	if idRoot != nil {
		if !bytes.Equal(idRoot, this.idRoot) {
			this.idRoot = idRoot
			changed = true
		}
	}

	if incarnation != nil {
		if incarnation != this.incarnation {
			this.incarnation = incarnation
			changed = true
		}
	}

	if changed {
		this.id, _ = toId(this.scope.Id(), this.idRoot, this.incarnation, this.crcLength)
	}

	return changed
}

func (this *domain) ToString() string {
	return this.ToEncodedString("base62")
}

func (this *domain) ToEncodedString(encoder string) string {
	str, _ := encoding.Encode(this.id, encoder)
	return str
}

func (this *domain) ToJSON() string {
	return this.ToEncodedString("base62")
}

func (this *domain) IsFor(typeId ids.TypeIdentifier) bool {
	isFor, _ := typeId.IsAssignableFrom(this.typeId)
	return isFor
}

func (this *domain) Matches(domain ids.Domain) bool {
	idRoot := domain.IdRoot()
	scope := domain.Scope()

	if this.idRoot == nil {
		return bytes.Equal(this.scope.Id(), scope.Id())
	}
	return bytes.Equal(this.idRoot, idRoot)
}

func (this *domain) Id() []byte {
	return this.id
}

func (this *domain) ScopeId() []byte {
	return this.id[:scopeLength(this.id)]
}

func (this *domain) Scope() ids.DomainScope {
	return this.scope
}

func (this *domain) IdRoot() []byte {
	return this.idRoot
}

func (this *domain) Incarnation() *uint32 {
	return this.incarnation
}

func (this *domain) CrcLength() uint {
	return this.crcLength
}

func (this *domain) Name() string {
	name, ok := this.info["name"]

	if ok {
		return name.(string)
	}

	if this.root != nil {
		return this.root.Name()
	}

	return ""
}

func (this *domain) Description() string {
	description, ok := this.info["description"]

	if ok {
		return description.(string)
	}

	if this.root != nil {
		return this.root.Description()
	}

	return ""
}

func (this *domain) TypeId() ids.TypeIdentifier {
	if this.typeId != nil {
		return this.typeId
	}
	if this.root != nil && this.root.TypeId() != nil {
		return this.root.TypeId()
	}

	return nil
}

func (this *domain) Source() string {
	source, ok := this.info["source"]

	if ok {
		return source.(string)
	}

	if this.root != nil {
		return this.root.Source()
	}

	return ""
}

func scopeLength(value []byte) uint {
	switch value[0] & 0x3f {
	// TODO handle named scopes
	default:
		return 1
	}
}

func domainLength(value []byte) uint {
	return uint(value[scopeLength(value)] & 0x3f)
}

func incarnationLength(value []byte) uint {
	return incarnationBitsLength(value[scopeLength(value)])
}

func incarnationValue(value []byte) uint32 {
	incLen := incarnationLength(value)

	if incLen > 0 {
		incOffset := scopeLength(value) +
			domainLength(value) - incLen + 1

		if incLen == 1 {
			return uint32(value[incOffset])
		}
		if incLen == 2 {
			return uint32(ntohs(value, int(incOffset)))
		}
		return ntohl(value, int(incOffset))
	}

	return 0
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

	length := uint(len(value)) - scopeLength(value) -
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
	scopeLength := scopeLength(value)
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

func incarnationLengthBits(incarnation []byte) (byte, error) {
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

func crcLengthBits(crcLength uint) (byte, error) {
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

type DomainType int

const (
	SCOPE DomainType = iota
	IDENTITY
	SIGNATURE
	SEQUENCE
)
