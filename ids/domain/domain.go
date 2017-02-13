package domain

import (
	"bytes"
	"errors"
	"fmt"
	"hash/crc32"
	"reflect"

	"github.com/distributed-vision/go-resources/encoding"
	"github.com/distributed-vision/go-resources/encoding/encoderType"
	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/version/versionType"
	"github.com/howeyc/crc16"
	"github.com/sigurn/crc8"
)

var DomainPathKey = "DV_DOMAIN_PATH"

type domain struct {
	id          []byte
	scope       ids.DomainScope
	root        ids.Domain
	idRoot      []byte
	incarnation *uint32
	crcLength   uint
	versionType versionType.VersionType
	typeId      ids.TypeIdentifier
	info        map[interface{}]interface{}
}

var empty = []byte{}

func NewDomain(scope ids.DomainScope, idRoot []byte, incarnation *uint32, crcLength uint, versionType versionType.VersionType, info map[interface{}]interface{}) (ids.Domain, error) {
	scopeId := empty

	if scope != nil {
		scopeId = scope.Id()
	}

	id, err := ToId(scopeId, idRoot, incarnation, crcLength, versionType)
	//fmt.Printf("id=%v\n", id)
	if err != nil {
		return nil, err
	}

	return &domain{
		id:          id,
		scope:       scope,
		idRoot:      idRoot,
		incarnation: incarnation,
		crcLength:   crcLength,
		versionType: versionType,
		info:        info}, nil
}

func (this *domain) NewIncarnation(incarnation uint32, crcLength uint, info map[interface{}]interface{}) (ids.Domain, error) {

	id, err := ToId(this.scope.Id(), this.idRoot, &incarnation, crcLength, this.versionType)

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
		versionType: this.versionType,
		info:        info}, nil
}

var featureBit byte = (1 << 6)

func ToId(scopeId []byte, idRoot []byte, incarnation *uint32, crcLength uint, versionType versionType.VersionType) ([]byte, error) {

	var incarnationSlice = IncarnationAsBytes(incarnation)

	unscoped := bytes.Join([][]byte{idRoot, incarnationSlice}, empty)
	//fmt.Printf("unscoped=%v\n", unscoped)
	if len(unscoped) > 61 {
		return nil, errors.New("Id too Long: domain id unscoped binary length (idRoot+incarnation) must be < 61")
	}

	var unscopedlenSlice []byte
	featureSlice := empty

	if len(scopeId) == 0 {
		unscopedlenSlice = empty
	} else {
		incarnationLengthBits, err := IncarnationLengthBits(incarnationSlice)

		if err != nil {
			return nil, err
		}

		crcLengthBits, err := CrcLengthBits(crcLength)

		if err != nil {
			return nil, err
		}

		versionTypeBits, err := VersionTypeBits(versionType)

		if err != nil {
			return nil, err
		}

		if crcLengthBits > 0 || incarnationLengthBits > 0 || versionTypeBits > 0 {
			featureSlice = []byte{crcLengthBits | incarnationLengthBits | versionTypeBits}
		}

		if len(featureSlice) > 0 {
			unscopedlenSlice = []byte{byte(len(unscoped)&0xff) | featureBit}
		} else {
			unscopedlenSlice = []byte{byte(len(unscoped) & 0xff)}
		}
	}

	return bytes.Join([][]byte{scopeId, unscopedlenSlice, featureSlice, unscoped}, empty), nil
}

func DecodeId(encoderType encoderType.EncoderType, scopeId string, rootId string, features ...interface{}) ([]byte, error) {
	var scopeIdValue []byte
	var idRootValue []byte
	var incarnationValue *uint32 = nil
	var crcLengthValue uint = 0
	var versionTypeValue versionType.VersionType = versionType.UNVERSIONED
	var err error

	scopeIdValue, err = encoding.Decode(scopeId, encoderType)

	if err != nil {
		return nil, fmt.Errorf("Invalid scopeId encoding %s", err)
	}

	idRootValue, err = encoding.Decode(rootId, encoderType)

	if err != nil {
		return nil, fmt.Errorf("Invalid rootId encoding %s", err)
	}

	if len(features) > 0 {

		if feature, ok := features[0].(*uint32); ok {
			incarnationValue = feature
		} else if feature, ok := features[0].(uint32); ok {
			incarnationValue = &feature
		} else {
			return nil, fmt.Errorf("Invalid incarnation type expected: *uint32, got: %s", reflect.ValueOf(features[0]).Type())
		}
	}

	if len(features) > 1 {
		if feature, ok := features[1].(uint); ok {
			crcLengthValue = feature
		} else {
			return nil, fmt.Errorf("Invalid crcLength type expected: uint, got: %s", reflect.ValueOf(features[1]).Type())
		}
	}

	if len(features) > 2 {
		if feature, ok := features[2].(versionType.VersionType); ok {
			versionTypeValue = feature
		} else {
			return nil, fmt.Errorf("Invalid versionType type expected: versionType.VersionType, got: %s", reflect.ValueOf(features[1]).Type())
		}
	}

	return ToId(scopeIdValue, idRootValue, incarnationValue, crcLengthValue, versionTypeValue)
}

func MustDecodeId(encoderType encoderType.EncoderType, scopeId string, rootId string, features ...interface{}) []byte {
	id, err := DecodeId(encoderType, scopeId, rootId, features...)

	if err != nil {
		panic(fmt.Sprintf("Failed to encode id: %s", err))
	}

	return id
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
		this.id, _ = ToId(this.scope.Id(), this.idRoot, this.incarnation, this.crcLength, this.versionType)
	}

	return changed
}

func (this *domain) String() string {
	return this.Encode(encoderType.BASE62)
}

func (this *domain) Encode(encoder encoderType.EncoderType) string {
	str, _ := encoding.Encode(this.id, encoder)
	return str
}

func (this *domain) ToJSON() string {
	return this.Encode(encoderType.BASE62)
}

func (this *domain) IsFor(typeId ids.TypeIdentifier) bool {
	isFor := typeId.IsAssignableFrom(this.typeId)
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
	return this.id[:ScopeLength(this.id)]
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

func (this *domain) VersionType() versionType.VersionType {
	return this.versionType
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

func (this *domain) InfoValue(key interface{}) interface{} {
	//fmt.Printf("info=%v\n", this.info)
	if this.info != nil {
		value, ok := this.info[key]
		if ok {
			return value
		}
	}

	return nil
}

func IncarnationLengthBits(incarnation []byte) (byte, error) {
	switch len(incarnation) {
	case 0:
		return 0, nil
	case 1:
		return 1, nil
	case 2:
		return 2, nil
	case 4:
		return 3, nil
	default:
		return 0, errors.New("Invalid incarnation length")
	}
}

func incarnationBitsLength(bits byte) uint {
	bits = (bits & 0x03)
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
		return 1 << 2, nil
	case 16:
		return 2 << 2, nil
	case 32:
		return 3 << 2, nil
	default:
		return 0, errors.New("Invalid crc length")
	}
}

func CrcLengthValue(value []byte) (uint, error) {

	featureSlice := featureSlice(value)

	if featureSlice == nil {
		return 0, nil
	}

	bits := (featureSlice[0] >> 2) & 0x03
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

func VersionTypeBits(verType versionType.VersionType) (byte, error) {
	switch verType {
	case versionType.UNVERSIONED:
		return 0, nil
	case versionType.NUMERIC:
		return 1 << 4, nil
	case versionType.SEMANTIC:
		return 2 << 4, nil
	default:
		return 0, errors.New("Invalid version type")
	}
}

func VersionTypeValue(value []byte) (versionType.VersionType, error) {

	featureSlice := featureSlice(value)

	if featureSlice == nil {
		return versionType.UNVERSIONED, nil
	}

	bits := (featureSlice[0] >> 4) & 0x03
	switch bits {
	case 0:
		return versionType.UNVERSIONED, nil
	case 1:
		return versionType.NUMERIC, nil
	case 2:
		return versionType.SEMANTIC, nil
	}

	return versionType.UNVERSIONED, errors.New("Invalid version type")
}

func VersionLength(value []byte) uint {
	lengthLength := VersionLengthLength(value)

	if lengthLength > 0 {
		domainOffset := ScopeLength(value) + 1 + featureSliceLength(value)
		domainLength := DomainLength(value)
		return uint(value[domainOffset+domainLength])
	}

	return 0
}

func NumericVersionValue(versionValue []byte) uint32 {
	vlen := len(versionValue)
	if vlen > 0 {
		if vlen == 1 {
			return uint32(versionValue[0])
		}
		if vlen == 2 {
			return uint32(ntohs(versionValue, 0))
		}
		return ntohl(versionValue, 0)
	}

	return 0
}

func VersionLengthLength(value []byte) uint {
	vt, err := VersionTypeValue(value)
	if err == nil {
		if vt != versionType.UNVERSIONED {
			return 1
		}
	}
	return 0
}

func ScopeLength(value []byte) uint {
	switch value[0] & 0x3f {
	// TODO handle named scopes
	default:
		return 1
	}
}

func IncarnationAsBytes(incarnation *uint32) []byte {
	var incarnationSlice []byte

	if incarnation == nil {
		incarnationSlice = empty
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

func IncarnationValue(value []byte) *uint32 {
	incLen := IncarnationLength(value)
	if incLen > 0 {
		incOffset := ScopeLength(value) +
			DomainLength(value) + featureSliceLength(value) - incLen + 1
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

func IncarnationLength(value []byte) uint {
	featureSlice := featureSlice(value)

	if featureSlice != nil {
		return incarnationBitsLength(featureSlice[0])
	}

	return 0
}

func DomainLength(value []byte) uint {
	return uint(value[ScopeLength(value)] & 0x3f)
}

func DomainOffset(value []byte) uint {
	return ScopeLength(value) + 1 + featureSliceLength(value)
}

func IdRootValue(value []byte) []byte {
	domainOffset := DomainOffset(value)
	return value[domainOffset : domainOffset+DomainLength(value)-IncarnationLength(value)]
}

func featureSlice(value []byte) []byte {
	if value[ScopeLength(value)]>>6 > 0 {
		featurePos := ScopeLength(value) + 1
		return value[featurePos : featurePos+1]
	}

	return nil
}

func featureSliceLength(value []byte) uint {
	if value[ScopeLength(value)]>>6 > 0 {
		return 1
	}

	return 0
}

var crc8Table *crc8.Table = crc8.MakeTable(crc8.CRC8_MAXIM)
var crc16Table *crc16.Table = crc16.MakeTable(crc16.IBM)
var crc32Table *crc32.Table = crc32.MakeTable(crc32.IEEE)

func CrcCalc(value []byte, crcLength uint) ([]byte, error) {
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
