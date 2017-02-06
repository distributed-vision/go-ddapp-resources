package identifier

import (
	"os"
	"testing"

	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/ids/domain"
	"github.com/distributed-vision/go-resources/ids/domainScope"
)

func setup(t *testing.T) {
	if os.Getenv("DV_DOMAIN_SCOPE_PATH") == "" {
		os.Setenv("DV_DOMAIN_SCOPE_PATH", "../../../domain-scope")
	}
}

func TestHashDomainXX64(t *testing.T) {

	setup(t)

	scopeId, err := base62.Decode("1")

	if err != nil {
		t.Errorf("TestHashDomainXX64 Failed: Can't decode scopeId")
	}

	domainId, err := base62.Decode("0")

	if err != nil {
		t.Errorf("TestHashDomainXX64 Failed: Can't decode domainId")
	}

	scopeNameSelector := domainScope.Selector{Name: "global hash", Opts: domainScope.SelectorOpts{
		IgnoreCase:       true,
		IgnoreWhitespace: true}}

	domainNameSelector := domain.Selector{Name: "xx64", Opts: domain.SelectorOpts{
		IgnoreCase:       true,
		IgnoreWhitespace: true}}

	scope, err := domainScope.Get(domainScope.Selector{Id: scopeId})

	if err != nil {
		t.Errorf("TestHashDomainXX64 Failed to resolve Scope: %s", err)
	} else {
		domain, derr := domain.Get(domain.Selector{Scope: scope, Id: domainId})

		if derr != nil {
			t.Errorf("TestHashDomainXX64 Failed to resolve Domain: %s", derr)
		} else {
			if base62.Encode(domain.Id()) != "H7A" {
				t.Errorf("TestHashDomainXX64 Failed: domain.Id: expected: '%s' got '%s'", "H7A", base62.Encode(domain.Id()))
			}

			if domain.Name() != "xx64" {
				t.Errorf("TestHashDomainXX64 Failed: domain.Name: expected: '%s' got '%s'", "xx64", domain.Name())

			}

			if base62.Encode(domain.IdRoot()) != "0" {
				t.Errorf("TestHashDomainXX64 Failed: domain.IdRoot: expected: '%s' got '%s'", "0", base62.Encode(domain.IdRoot()))

			}

			if base62.Encode(domain.Scope().Id()) != "1" {
				t.Errorf("TestHashDomainXX64 Failed: domain.Scope.Id: expected: '%s' got '%s'", "1", base62.Encode(domain.Scope().Id()))
			}
		}
	}

	scope, err = domainScope.Get(scopeNameSelector)

	if err != nil {
		t.Errorf("TestHashDomainXX64 Failed to resolve Scope: %s", err)
	} else {
		domain, derr := domain.Get(domainNameSelector)
		if derr != nil {
			t.Errorf("TestHashDomainXX64 Failed to resolve Domain: %s", derr)
		} else {
			if base62.Encode(domain.Id()) != "H7A" {
				t.Errorf("TestHashDomainXX64 Failed: domain.Id: expected: '%s' got '%s'", "H7A", base62.Encode(domain.Id()))
			}

			if domain.Name() != "xx64" {
				t.Errorf("TestHashDomainXX64 Failed: domain.Name: expected: '%s' got '%s'", "xx64", domain.Name())

			}

			if base62.Encode(domain.IdRoot()) != "0" {
				t.Errorf("TestHashDomainXX64 Failed: domain.IdRoot: expected: '%s' got '%s'", "0", base62.Encode(domain.IdRoot()))

			}

			if base62.Encode(domain.Scope().Id()) != "1" {
				t.Errorf("TestHashDomainXX64 Failed: domain.Scope.Id: expected: '%s' got '%s'", "1", base62.Encode(domain.Scope().Id()))
			}
		}
	}
}

/*
	func TestHashDomainXX32(t *testing.T) {

    let scopeId = base62.decode('1')
    let domainId = base62.decode('1')
    let scopeNameSelector = new Selector({name: 'global hash'}, {
      ignoreCase: true,
      ignoreWhitespace: true
    })
    let domainNameSelector = new Selector({name: 'xx32'}, {
      ignoreCase: true,
      ignoreWhitespace: true
    })

    return DomainScopeResolver.resolve(scopeId)
      .then(scope => scope.resolve(domainId))
      .then(domain => {
        should(base62.encode(domain.id)).equal('H7B')
        should(domain.name).equal('xx32')
        should(base62.encode(domain.idRoot)).equal('1')
        should(base62.encode(domain.scope.id)).equal('1')
      })
      .then(() => {
        return DomainScopeResolver.resolve(scopeNameSelector)
          .then(scope => scope.resolve(domainNameSelector))
          .then(domain => {
            should(base62.encode(domain.id)).equal('H7B')
            should(domain.name).equal('xx32')
            should(base62.encode(domain.idRoot)).equal('1')
            should(base62.encode(domain.scope.id)).equal('1')
          })
      })
  })

  it("resolve ripemd160", function() {
    let scopeId = base62.decode('1')
    let domainId = base62.decode('2')
    let scopeNameSelector = new Selector({name: 'global hash'}, {
      ignoreCase: true,
      ignoreWhitespace: true
    })
    let domainNameSelector = new Selector({name: 'ripemd160'}, {
      ignoreCase: true,
      ignoreWhitespace: true
    })

    return DomainScopeResolver.resolve(scopeId)
      .then(scope => scope.resolve(domainId))
      .then(domain => {
        should(base62.encode(domain.id)).equal('H7C')
        should(domain.name).equal('ripemd160')
        should(base62.encode(domain.idRoot)).equal('2')
        should(base62.encode(domain.scope.id)).equal('1')
      })
      .then(() => {
        return DomainScopeResolver.resolve(scopeNameSelector)
          .then(scope => scope.resolve(domainNameSelector))
          .then(domain => {
            should(base62.encode(domain.id)).equal('H7C')
            should(domain.name).equal('ripemd160')
            should(base62.encode(domain.idRoot)).equal('2')
            should(base62.encode(domain.scope.id)).equal('1')
          })
      })
  })

  it("resolve sha256", function() {
    let scopeId = base62.decode('1')
    let domainId = base62.decode('3')
    let scopeNameSelector = new Selector({name: 'global hash'}, {
      ignoreCase: true,
      ignoreWhitespace: true
    })
    let domainNameSelector = new Selector({name: 'sha256'}, {
      ignoreCase: true,
      ignoreWhitespace: true
    })

    return DomainScopeResolver.resolve(scopeId)
      .then(scope => scope.resolve(domainId))
      .then(domain => {
        should(base62.encode(domain.id)).equal('H7D')
        should(domain.name).equal('sha256')
        should(base62.encode(domain.idRoot)).equal('3')
        should(base62.encode(domain.scope.id)).equal('1')
      })
      .then(() => {
        return DomainScopeResolver.resolve(scopeNameSelector)
          .then(scope => scope.resolve(domainNameSelector))
          .then(domain => {
            should(base62.encode(domain.id)).equal('H7D')
            should(domain.name).equal('sha256')
            should(base62.encode(domain.idRoot)).equal('3')
            should(base62.encode(domain.scope.id)).equal('1')
          })
      })
  })
})

describe('Domain with Incarnations', function() {
  let scopeId = base62.decode('1')
  let domainId = base62.decode('3')

  it("create domain with incarnation", function() {
    return DomainScopeResolver.resolve(scopeId)
      .then(scope => scope.resolve(domainId))
      .then(domain => {
        let incarnation0 = new IdentityDomain(domain.scope, domain, 0)
        should(base62.encode(incarnation0.id)).equal('1QXwm')
        should(base62.encode(incarnation0.idRoot)).equal('3')
        should(incarnation0.incarnation).equal(0)
        should(incarnation0.name).equal('sha256')

        let incarnation300 = new IdentityDomain(domain.scope, domain, 300)
        should(base62.encode(incarnation300.id)).equal('5ukZHI')
        should(base62.encode(incarnation300.idRoot)).equal('3')
        should(incarnation300.incarnation).equal(300)
        should(incarnation300.name).equal('sha256')

        let incarnation300000 = new IdentityDomain(domain.scope, domain, 300000)
        should(base62.encode(incarnation300000.id)).equal('2HRBKvjii')
        should(base62.encode(incarnation300000.idRoot)).equal('3')
        should(incarnation300000.incarnation).equal(300000)
        should(incarnation300000.name).equal('sha256')

        should(domain.isRoot).be.true

        should(incarnation0.isRoot).be.false
        should(incarnation300.isRoot).be.false
        should(incarnation300000.isRoot).be.false

        should(domain.isRootOf(incarnation0)).be.true
        should(domain.isRootOf(incarnation300)).be.true
        should(domain.isRootOf(incarnation300000)).be.true

        should(domain.matches(incarnation0)).be.true
        should(domain.matches(incarnation300)).be.true
        should(domain.matches(incarnation300000)).be.true

        should(incarnation0.matches(domain)).be.true
        should(incarnation0.matches(incarnation300)).be.true
        should(incarnation0.matches(incarnation300000)).be.true

        should(incarnation300.matches(domain)).be.true
        should(incarnation300.matches(incarnation0)).be.true
        should(incarnation300.matches(incarnation300000)).be.true

        should(incarnation300000.matches(domain)).be.true
        should(incarnation300000.matches(incarnation0)).be.true
        should(incarnation300000.matches(incarnation300)).be.true
      })
  })
})

describe('Domain with Crc', function() {
  let scopeId = base62.decode('1')
  let domainId = base62.decode('3')

  it("create domain with crc", function() {
    return DomainScopeResolver.resolve(scopeId)
      .then(scope => scope.resolve(domainId))
      .then(domain => {
        let crc8 = new IdentityDomain(domain.scope, domain, undefined, 8)
        should(base62.encode(crc8.id)).equal('HsFH')
        should(base62.encode(crc8.idRoot)).equal('3')
        should(crc8.crcLength).equal(8)
        should(crc8.name).equal('sha256')

        let crc16 = new IdentityDomain(domain.scope, domain, undefined, 16)
        should(base62.encode(crc16.id)).equal('ZTNL')
        should(base62.encode(crc16.idRoot)).equal('3')
        should(crc16.crcLength).equal(16)
        should(crc16.name).equal('sha256')

        let crc32 = new IdentityDomain(domain.scope, domain, undefined, 32)
        should(base62.encode(crc32.id)).equal('r4VP')
        should(base62.encode(crc32.idRoot)).equal('3')
        should(crc32.crcLength).equal(32)
        should(crc32.name).equal('sha256')

        should(domain.isRoot).be.true

        should(crc8.isRoot).be.false
        should(crc16.isRoot).be.false
        should(crc32.isRoot).be.false

        should(domain.isRootOf(crc8)).be.true
        should(domain.isRootOf(crc16)).be.true
        should(domain.isRootOf(crc32)).be.true

        should(domain.matches(crc8)).be.true
        should(domain.matches(crc16)).be.true
        should(domain.matches(crc32)).be.true

        should(crc8.matches(domain)).be.true
        should(crc8.matches(crc16)).be.true
        should(crc8.matches(crc32)).be.true

        should(crc16.matches(domain)).be.true
        should(crc16.matches(crc8)).be.true
        should(crc16.matches(crc32)).be.true

        should(crc32.matches(domain)).be.true
        should(crc32.matches(crc8)).be.true
        should(crc32.matches(crc32)).be.true
      })
  })
})
*/
