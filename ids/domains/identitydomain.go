package domains

import (
	"bytes"
	"errors"
	"time"

	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/identifiers"
	"github.com/distributed-vision/go-resources/resolvers/mappingResolver"
)

type identityDomain struct {
	*domain
	idGenerator             ids.IdGenerator
	sequenceLeaseHolderId   ids.Identifier
	nextSequenceLeaseExpiry *time.Time
	sequenceRoot            []byte
	sequenceIncarnation     *uint32
}

func (scope *identityDomain) NewIdentityDomainFromJSON(json map[string]interface{}, opts map[string]interface{}) (ids.IdentityDomain, error) {
	id, err := base62.Decode(json["id"].(string))

	if err != nil {
		return nil, err
	}

	return &identityDomain{&domain{scope: opts["scope"].(ids.DomainScope), id: id, info: json}, nil, nil, nil, nil, nil}, nil
}

func NewIdentityDomain(scope ids.DomainScope, rootId []byte, incarnation *uint32, crcLength uint, info map[string]interface{}) (ids.IdentityDomain, error) {
	base, err := NewDomain(scope, rootId, incarnation, crcLength, info)

	if err != nil {
		return nil, err
	}

	return &identityDomain{
		domain: base.(*domain)}, nil
}

func (this *identityDomain) NewIdentifier() (ids.Identifier, error) {
	if this.idGenerator != nil {
		id, err := this.idGenerator.GenerateId()

		if err != nil {
			return nil, err
		}

		return this.ToId(id)
	}

	return nil, errors.New("Id creation not initialized")
}

func (this *identityDomain) ToId(id []byte) (ids.Identifier, error) {
	return identifiers.NewIdentifier(this, id)
}

func (this *identityDomain) GetMapping(from ids.Identifier, asof *time.Time) (ids.Identifier, error) {
	mchan, echan := this.ResolveMapping(from, asof)

	select {
	case mapping := <-mchan:
		return mapping, nil
	case err := <-echan:
		return nil, err
	}
}

func (this *identityDomain) ResolveMapping(from ids.Identifier, asof *time.Time) (chan ids.Identifier, chan error) {
	return mappingResolver.ResolveMapping(from, this, *asof)
}

func (this *identityDomain) IsRootOf(domain ids.Domain) bool {
	return this.incarnation == nil && this.crcLength == 0 && bytes.Equal(domain.IdRoot(), this.idRoot)
}

func (this *identityDomain) IsRoot() bool {
	return this.incarnation == nil && this.crcLength == 0
}

func (this *identityDomain) AsRoot() ids.IdentityDomain {
	domain, _ := NewIdentityDomain(this.Scope(), this.IdRoot(), nil, 0, nil)
	return domain
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
	id, _ := toId(this.ScopeId(), this.sequenceRoot, this.sequenceIncarnation, 0)
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
