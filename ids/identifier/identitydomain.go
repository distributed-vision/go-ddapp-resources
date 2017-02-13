package identifier

import (
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/domain"
	"github.com/distributed-vision/go-resources/ids/domainType"
	"github.com/distributed-vision/go-resources/ids/mappings"
	"github.com/distributed-vision/go-resources/version"
	"github.com/distributed-vision/go-resources/version/versionType"
)

func init() {
	domain.RegisterJSONUnmarshaller(domainType.IDENTITY, unmarshalJSON)
}

type identityDomain struct {
	ids.Domain
	idGenerator             ids.IdGenerator
	sequenceLeaseHolderId   ids.Identifier
	nextSequenceLeaseExpiry *time.Time
	sequenceRoot            []byte
	sequenceIncarnation     *uint32
}

func unmarshalJSON(unmarshalContext context.Context, json map[string]interface{}) (ids.Domain, error) {
	rootId, err := base62.Decode(json["id"].(string))

	if err != nil {
		return nil, err
	}

	var incarnation *uint32
	var crcLength uint
	var versionType = versionType.UNVERSIONED

	info := make(map[interface{}]interface{})

	for key, value := range json {
		info[key] = value
	}

	return NewIdentityDomain(unmarshalContext.Value("scope").(ids.DomainScope), rootId, incarnation, crcLength, versionType, info)
}

func NewIdentityDomain(scope ids.DomainScope, rootId []byte, incarnation *uint32, crcLength uint, versionType versionType.VersionType, info map[interface{}]interface{}) (ids.IdentityDomain, error) {
	base, err := domain.NewDomain(scope, rootId, incarnation, crcLength, versionType, info)

	if err != nil {
		return nil, err
	}

	return &identityDomain{
		Domain: base}, nil
}

func (this *identityDomain) NewIdentifier() (ids.Identifier, error) {
	if this.idGenerator != nil {
		id, err := this.idGenerator.GenerateId()

		if err != nil {
			return nil, err
		}

		var nilVersion version.Version
		return this.ToId(id, nilVersion)
	}

	return nil, errors.New("Id creation not initialized")
}

func (this *identityDomain) ToId(id []byte, ver version.Version) (ids.Identifier, error) {
	return New(this, id, ver)
}

func (this *identityDomain) GetMapping(from ids.Identifier, asof *time.Time) (ids.Identifier, error) {
	mchan, echan := this.ResolveMapping(context.Background(), from, asof)

	select {
	case mapping := <-mchan:
		return mapping, nil
	case err := <-echan:
		return nil, err
	}
}

func (this *identityDomain) ResolveMapping(resolutionContext context.Context, from ids.Identifier, asof *time.Time) (chan ids.Identifier, chan error) {
	return mappings.Resolve(resolutionContext, mappings.Selector{From: from, To: this, At: *asof})
}

func (this *identityDomain) IsRootOf(domain ids.Domain) bool {
	return this.Incarnation() == nil && this.CrcLength() == 0 && bytes.Equal(domain.IdRoot(), this.IdRoot())
}

func (this *identityDomain) IsRoot() bool {
	return this.Incarnation() == nil && this.CrcLength() == 0
}

func (this *identityDomain) AsRoot() ids.IdentityDomain {
	root, _ := NewIdentityDomain(this.Scope(), this.IdRoot(), nil, 0, this.VersionType(), nil)
	return root
}

func (this *identityDomain) HasCurrentMapping(from ids.Identifier) bool {
	now := time.Now()
	return this.HasMapping(from, &now)
}

func (this *identityDomain) HasMapping(from ids.Identifier, asof *time.Time) bool {
	_, err := this.GetMapping(from, asof)
	return err == nil
}

func (this *identityDomain) SetIfChanged(root []byte, incarnation *int32, name string, description string, source string, typeId ids.TypeIdentifier) bool {

	//changed = super.setIfChanged(root, incarnation)

	/*
		if (type!=null) {
			if (!type.equals(this.type)) {
				this.type=type;
				changed=true;
			}
		}

		if (!(name==null||"".equals(name))) {
			if (!name.equals(this.name)) {
				this.name = name;
				changed=true;
			}
		}

		if (!(description==null||"".equals(description))) {
			if (!description.equals(this.description)) {
				this.description = description;
				changed=true;
			}
		}

		if (!(source==null||"".equals(source))) {
			if (!source.equals(this.source)) {
				this.source = source;
				changed=true;
			}
		}
	*/

	return false //changed
}

func (this *identityDomain) SequenceLeaseHolderId() ids.Identifier {
	return this.sequenceLeaseHolderId
}

func (this *identityDomain) HasSequenceLeaseHolderId() bool {
	return this.sequenceLeaseHolderId != nil
}

func (this *identityDomain) NextSequenceLeaseExpiry() *time.Time {
	return this.nextSequenceLeaseExpiry
}

func (this *identityDomain) SequenceDomain() ids.SequenceDomain {
	//return DomainStore.getSequenceDomain((IdentityDomainInfo)domainInfo);
	return nil
}

func (this *identityDomain) SequenceRoot() []byte {
	if this.sequenceRoot == nil {
		return nil
	}
	return this.sequenceRoot
}

func (this *identityDomain) SequenceDomainId() []byte {
	id, _ := domain.ToId(this.ScopeId(), this.sequenceRoot, this.sequenceIncarnation, 0, versionType.UNVERSIONED)
	return id
}

func (this *identityDomain) SequenceIncarnation() uint32 {
	if this.sequenceIncarnation == nil {
		return 0
	}
	return *this.sequenceIncarnation
}

/*
	aquireSequenceDomain( int incarnation, Period leasePeriod, SequenceNumber<Type> lastSequence ) {

		// TODO - this really needs a lock - so should be in a transaction

		IdentityDomainInfo identityDomainInfo=(IdentityDomainInfo)domainInfo;
		Date nextExpiry=identityDomainInfo.nextSequenceLeaseExpiry;

		if (nextExpiry==null||LeaseHolderIdentifier.VM_ID.equals(getSequenceLeaseHolderId())||
				nextExpiry.getTime()<System.currentTimeMillis()) {

			if (identityDomainInfo.sequenceIncarnation==null||incarnation>identityDomainInfo.sequenceIncarnation)
				identityDomainInfo.sequenceIncarnation=incarnation;
			identityDomainInfo.nextSequenceLeaseExpiry=DateTime.now().plus(leasePeriod).toDate();
			identityDomainInfo.sequenceLeaseHolderId=LeaseHolderIdentifier.VM_ID;
			DomainStore.post(domainInfo);

		} else {
			throw new LeaseAquisitionException(identityDomainInfo.nextSequenceLeaseExpiry);
		}

		return DomainStore.getSequenceDomain(identityDomainInfo,identityDomainInfo.getSequenceRoot(),identityDomainInfo.sequenceIncarnation);
	}

	void releaseSequenceDomain() {
		((IdentityDomainInfo)domainInfo).nextSequenceLeaseExpiry=null;
		((IdentityDomainInfo)domainInfo).sequenceLeaseHolderId=null;
		DomainStore.post(domainInfo);
	}

	getNextSequenceNumber() {
		return getSequenceDomain().nextSequenceNumber();
	}
*/

func (this *identityDomain) RegisterMappingProvider(to ids.IdentityDomain, provider func() (ids.IdentityDomain, error)) (ids.IdentityDomain, error) {
	//MAPPING_PROVIDERS.set(to, provider)
	return nil, nil
}

func (this *identityDomain) UnregisterMappingProvider(to ids.IdentityDomain) {
	//MAPPING_PROVIDERS.delete(to)
}
