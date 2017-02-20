package domain

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/distributed-vision/go-resources/encoding"
	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/encoding/encodertype"
	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/resolvers"
	"github.com/distributed-vision/go-resources/util/hton"
	"github.com/distributed-vision/go-resources/util/ntoh"
	"github.com/distributed-vision/go-resources/version/versiontype"
)

var DomainPathKey = "DV_DOMAIN_PATH"

type domain struct {
	id          []byte
	scope       ids.DomainScope
	root        ids.Domain
	idRoot      []byte
	incarnation *uint32
	crcLength   uint
	versionType versiontype.VersionType
	typeId      ids.TypeIdentifier
	info        map[interface{}]interface{}
}

func KeyExtractor(entity ...interface{}) (interface{}, bool) {
	if len(entity) > 0 {
		if domain, ok := entity[0].(ids.Domain); ok {
			id, err := ToId(domain.ScopeId(), domain.IdRoot(), nil, 0, versiontype.UNVERSIONED)

			if err == nil {
				return base62.Encode(id), true
			}
		}
	}
	return nil, false
}

func Await(cres chan ids.Domain, cerr chan error) (result ids.Domain, err error) {
	if cres == nil || cerr == nil {
		return nil, fmt.Errorf("Await Failed: channels are undefined")
	}

	resolved := false
	for !resolved {
		select {
		case res, ok := <-cres:
			if ok {
				result = res
				resolved = true
			}
		case error, ok := <-cerr:
			if ok {
				err = error
				resolved = true
			}
		}
	}

	return result, err
}

var empty = []byte{}

func New(scopeId []byte, idRoot []byte, incarnation *uint32, crcLength uint, versionType versiontype.VersionType, infos ...map[interface{}]interface{}) (ids.Domain, error) {

	id, err := ToId(scopeId, idRoot, incarnation, crcLength, versionType)
	//fmt.Printf("id=%v\n", id)
	if err != nil {
		return nil, err
	}

	var info map[interface{}]interface{}

	if len(infos) > 0 {
		info = make(map[interface{}]interface{})
		for _, infoVal := range infos {
			for key, value := range infoVal {
				info[key] = value
			}
		}
	}

	return &domain{
		id:          id,
		scope:       nil,
		idRoot:      idRoot,
		incarnation: incarnation,
		crcLength:   crcLength,
		versionType: versionType,
		info:        info}, nil
}

func Wrap(id []byte) ids.Domain {
	crcLength, _ := CrcLengthValue(id)
	versionType, _ := VersionTypeValue(id)

	return &domain{
		id:          id,
		scope:       nil,
		idRoot:      IdRootValue(id),
		incarnation: IncarnationValue(id),
		crcLength:   crcLength,
		versionType: versionType}
}

func WithIncarnation(root ids.Domain, incarnation uint32, crcLength uint, infos ...map[interface{}]interface{}) (ids.Domain, error) {

	id, err := ToId(root.ScopeId(), root.IdRoot(), &incarnation, crcLength, root.VersionType())

	if err != nil {
		return nil, err
	}

	var info map[interface{}]interface{}

	if len(infos) > 0 {
		info = make(map[interface{}]interface{})
		for _, infoVal := range infos {
			for key, value := range infoVal {
				info[key] = value
			}
		}
	}

	return &domain{
		id:          id,
		scope:       nil,
		root:        root,
		idRoot:      root.IdRoot(),
		incarnation: &incarnation,
		crcLength:   crcLength,
		versionType: root.VersionType(),
		info:        info}, nil
}

func WithCrc(root ids.Domain, crcLength uint, infos ...map[interface{}]interface{}) (ids.IdentityDomain, error) {

	id, err := ToId(root.ScopeId(), root.IdRoot(), nil, crcLength, root.VersionType())

	if err != nil {
		return nil, err
	}

	var info map[interface{}]interface{}

	if len(infos) > 0 {
		info = make(map[interface{}]interface{})
		for _, infoVal := range infos {
			for key, value := range infoVal {
				info[key] = value
			}
		}
	}

	var incarnation uint32

	if root.Incarnation() != nil {
		incarnation = *root.Incarnation()
	}

	return &domain{
		id:          id,
		scope:       nil,
		root:        root,
		idRoot:      root.IdRoot(),
		incarnation: &incarnation,
		crcLength:   crcLength,
		versionType: root.VersionType(),
		info:        info}, nil
}

var featureBit byte = (1 << 6)

func ToId(scopeId []byte, idRoot []byte, incarnation *uint32, crcLength uint, versionType versiontype.VersionType) ([]byte, error) {

	var incarnationSlice = IncarnationAsBytes(incarnation)

	unscoped := bytes.Join([][]byte{idRoot, incarnationSlice}, empty)
	//fmt.Printf("unscoped=%v\n", unscoped)
	if len(unscoped) > 61 {
		return nil, errors.New("Id too Long: domain id unscoped binary length (idRoot+incarnation) must be < 61")
	}

	var unscopedlenSlice []byte
	featureSlice := empty

	if scopeId == nil || len(scopeId) == 0 {
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

func DecodeId(encoderType encodertype.EncoderType, scopeId string, rootId string, features ...interface{}) ([]byte, error) {
	var scopeIdValue []byte
	var idRootValue []byte
	var incarnationValue *uint32 = nil
	var crcLengthValue uint = 0
	var versionTypeValue versiontype.VersionType = versiontype.UNVERSIONED
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
		if feature, ok := features[2].(versiontype.VersionType); ok {
			versionTypeValue = feature
		} else {
			return nil, fmt.Errorf("Invalid versionType type expected: versionType.VersionType, got: %s", reflect.ValueOf(features[1]).Type())
		}
	}

	return ToId(scopeIdValue, idRootValue, incarnationValue, crcLengthValue, versionTypeValue)
}

func MustDecodeId(encoderType encodertype.EncoderType, scopeId string, rootId string, features ...interface{}) []byte {
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
	return this.Encode(encodertype.BASE62)
}

func (this *domain) Encode(encoderType encodertype.EncoderType) string {
	str, _ := encoding.Encode(this.id, encoderType)
	return str
}

func (this *domain) ToJSON() string {
	return this.Encode(encodertype.BASE62)
}

func (this *domain) IsFor(typeId ids.TypeIdentifier) bool {
	isFor := typeId.IsAssignableFrom(this.typeId)
	return isFor
}

func (this *domain) Matches(domain ids.Domain) bool {
	idRoot := domain.IdRoot()
	scopeId := domain.ScopeId()

	if this.idRoot == nil {
		return bytes.Equal(this.scope.Id(), scopeId)
	}
	return bytes.Equal(this.idRoot, idRoot)
}

func (this *domain) Id() []byte {
	return this.id
}

func (this *domain) ScopeId() []byte {
	return ScopeId(this.id)
}

func (this *domain) Scope() ids.DomainScope {
	if this.scope == nil {
		if this.root != nil {
			return this.root.Scope()
		} else {
			res, err := resolvers.Get(context.Background(), &scopeSelector{this.ScopeId()})
			if err == nil {
				this.scope = res.(ids.DomainScope)
			}
		}
	}
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

func (this *domain) IsRootOf(domain ids.Domain) bool {
	return this.IsRoot() && bytes.Equal(domain.IdRoot(), this.IdRoot())
}

func (this *domain) IsRoot() bool {
	return this.incarnation == nil && this.crcLength == 0
}

func (this *domain) VersionType() versiontype.VersionType {
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

func VersionTypeBits(versionType versiontype.VersionType) (byte, error) {
	switch versionType {
	case versiontype.UNVERSIONED:
		return 0, nil
	case versiontype.NUMERIC:
		return 1 << 4, nil
	case versiontype.SEMANTIC:
		return 2 << 4, nil
	default:
		return 0, errors.New("Invalid version type")
	}
}

func VersionTypeValue(value []byte) (versiontype.VersionType, error) {

	featureSlice := featureSlice(value)

	if featureSlice == nil {
		return versiontype.UNVERSIONED, nil
	}

	bits := (featureSlice[0] >> 4) & 0x03
	switch bits {
	case 0:
		return versiontype.UNVERSIONED, nil
	case 1:
		return versiontype.NUMERIC, nil
	case 2:
		return versiontype.SEMANTIC, nil
	}

	return versiontype.UNVERSIONED, errors.New("Invalid version type")
}

func VersionLengthLength(value []byte) uint {
	vt, err := VersionTypeValue(value)
	if err == nil {
		if vt != versiontype.UNVERSIONED {
			return 1
		}
	}
	return 0
}

var extensionBit byte = (1 << 6)

func ScopeLength(value []byte) uint {
	if (value[0] & extensionBit) != 0 {
		return uint(value[1] + 2)
	}
	return 1
}

func ScopeId(value []byte) []byte {
	return value[:ScopeLength(value)]
}

func IncarnationAsBytes(incarnation *uint32) []byte {
	var incarnationSlice []byte

	if incarnation == nil {
		incarnationSlice = empty
	} else if *incarnation < 0xff {
		incarnationSlice = []byte{byte(*incarnation & 0xff)}
	} else if *incarnation < 0xffff {
		buf := [2]byte{0, 0}
		incarnationSlice = hton.U16(buf[:], 0, uint16(*incarnation&0xffff))
	} else {
		buf := [4]byte{0, 0, 0, 0}
		incarnationSlice = hton.U32(buf[:], 0, *incarnation)
	}

	return incarnationSlice
}

func IncarnationValue(value []byte) *uint32 {
	incLen := IncarnationLength(value)
	if incLen > 0 {
		incOffset := ScopeLength(value) +
			DomainLength(value) + FeatureSliceLength(value) - incLen + 1
		if incLen == 1 {
			res := uint32(value[incOffset])
			return &res
		}
		if incLen == 2 {
			res := uint32(ntoh.U16(value, int(incOffset)))
			return &res
		}
		res := ntoh.U32(value, int(incOffset))
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
	return ScopeLength(value) + 1 + FeatureSliceLength(value)
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

func FeatureSliceLength(value []byte) uint {
	if value[ScopeLength(value)]>>6 > 0 {
		return 1
	}
	return 0
}
