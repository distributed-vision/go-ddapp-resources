package domains

import (
	"bytes"
	"errors"

	"github.com/distributed-vision/go-resources/encoding"
	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/identifiers"
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

	var incarnationSlice = identifiers.IncarnationAsBytes(incarnation)

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
		incarnationLengthBits, err := identifiers.IncarnationLengthBits(incarnationSlice)
		if err != nil {
			return nil, err
		}
		unscopedlenSlice = []byte{byte(len(unscoped)&0xff) | incarnationLengthBits}
	}

	if crcLength > 0 {
		crcLengthBits, err := identifiers.CrcLengthBits(crcLength)
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
	return this.id[:identifiers.ScopeLength(this.id)]
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
