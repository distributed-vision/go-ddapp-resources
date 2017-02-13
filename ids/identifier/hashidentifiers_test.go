package identifier_test

import (
	"os"
	"testing"

	"github.com/distributed-vision/go-resources/init/idsinit"
)

/*
const random = require('../utils/random');
const crc = require('crc-hash')
const XXHash = require('xxhash')
const XXHash64 = XXHash.XXHash64
*/

func TestMain(m *testing.M) {
	idsinit.Init()
	os.Exit(m.Run())
}

func TestCreateXX64Id(t *testing.T) {
	/*    let scopeId = base62.decode('1')
	      let domainId = base62.decode('0')

	      return DomainScopeResolver.resolve(scopeId)
	        .then(scope => scope.resolve(domainId))
	        .then(domain => {
	          let value = Buffer.from(random.genRandomString(64))
	          let xx = new XXHash64(0xCAFEBABE);
	          xx.update(value)
	          let xxdigest = new Buffer(8);
	          xx.digest(xxdigest)

	          let id = new Identifier(domain, xxdigest)

	          should.ok(id.domainId.equals(domain.id))
	          should.ok(id.id.equals(xxdigest))
	          should.ok(id.scopeId.equals(domain.scopeId))
	          should(id.checksum).be.null
	        })*/
}

func TestCreateXX64IdWithCrc(t *testing.T) {
	/*let scopeId = base62.decode('1')
	  let domainId = base62.decode('0')

	  return DomainScopeResolver.resolve(scopeId)
	    .then(scope => scope.resolve(domainId))
	    .then(domain => {

	      let crc8 = new IdentityDomain(domain.scope, domain, undefined, 8)

	      let value = Buffer.from(random.genRandomString(64))
	      let xx = new XXHash64(0xCAFEBABE);
	      xx.update(value)
	      let xxdigest = new Buffer(8);
	      xx.digest(xxdigest)

	      let id = new Identifier(domain, xxdigest)
	      let crcid = new Identifier(crc8, xxdigest)
	      let cs = crc.createHash("crc8").update(
	        crcid.value.slice(0, crcid.value.length - 1)).digest();

	      should.ok(id.domainId.equals(domain.id))
	      should.ok(id.id.equals(xxdigest))
	      should.ok(id.scopeId.equals(domain.scopeId))
	      should(id.checksum).be.null
	      should.ok(id.isValid)

	      should.ok(crcid.domainId.equals(crc8.id))
	      should.ok(crcid.domainIdRoot.equals(crc8.idRoot))
	      should.ok(crcid.domainIdRoot.equals(id.domainIdRoot))
	      should.ok(crcid.id.equals(xxdigest))
	      should.ok(crcid.scopeId.equals(domain.scopeId))
	      should.ok(crcid.checksum.equals(cs))
	      should.ok(crcid.isValid)

	      should(id.toString('base62', 'base62')).equal(crcid.toString('base62', 'base62'))
	    })*/
}

func TestCreateXX64IdWithIncarnationCrc(t *testing.T) {
	/*let scopeId = base62.decode('1')
	  let domainId = base62.decode('0')

	  return DomainScopeResolver.resolve(scopeId)
	    .then(scope => scope.resolve(domainId))
	    .then(domain => {

	      let crc16_30 = new IdentityDomain(domain.scope, domain, 30, 16)

	      let value = Buffer.from(random.genRandomString(64))
	      let xx = new XXHash64(0xCAFEBABE);
	      xx.update(value)
	      let xxdigest = new Buffer(8);
	      xx.digest(xxdigest)

	      let id = new Identifier(domain, xxdigest)
	      let crcid = new Identifier(crc16_30, xxdigest)
	      let cs = crc.createHash("crc16").update(
	        crcid.value.slice(0, crcid.value.length - 2)).digest();

	      should.ok(id.domainId.equals(domain.id))
	      should.ok(id.id.equals(xxdigest))
	      should.ok(id.scopeId.equals(domain.scopeId))
	      should(id.checksum).be.null
	      should.ok(id.isValid)

	      should.ok(crcid.domainId.equals(crc16_30.id))
	      should(crcid.domainIncarnation).equal(crc16_30.incarnation)
	      should.ok(crcid.id.equals(xxdigest))
	      should.ok(crcid.scopeId.equals(domain.scopeId))
	      should.ok(crcid.checksum.equals(cs))
	      should.ok(crcid.isValid)
	    })
	})*/
}
