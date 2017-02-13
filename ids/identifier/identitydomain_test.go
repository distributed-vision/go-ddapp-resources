package identifier_test

import (
	"context"
	"testing"

	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/ids/domain"
	"github.com/distributed-vision/go-resources/ids/domainScope"
)

func TestHashDomainXX64(t *testing.T) {

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

	scope, err := domainScope.Get(context.Background(), domainScope.Selector{Id: scopeId})

	if err != nil {
		t.Errorf("TestHashDomainXX64 Failed to resolve Scope: %s", err)
	} else {
		domain, derr := domain.Get(context.Background(), domain.Selector{Scope: scope, Id: domainId})

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

	scope, err = domainScope.Get(context.Background(), scopeNameSelector)

	if err != nil {
		t.Errorf("TestHashDomainXX64 Failed to resolve Scope: %s", err)
	} else {
		domain, derr := domain.Get(context.Background(), domainNameSelector)

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

func TestHashDomainXX32(t *testing.T) {

	scopeId, err := base62.Decode("1")

	if err != nil {
		t.Errorf("TestHashDomainXX32 Failed: Can't decode scopeId")
	}

	domainId, err := base62.Decode("1")

	if err != nil {
		t.Errorf("TestHashDomainXX32 Failed: Can't decode domainId")
	}

	scopeNameSelector := domainScope.Selector{Name: "global hash", Opts: domainScope.SelectorOpts{
		IgnoreCase:       true,
		IgnoreWhitespace: true}}

	domainNameSelector := domain.Selector{Name: "xx32", Opts: domain.SelectorOpts{
		IgnoreCase:       true,
		IgnoreWhitespace: true}}

	scope, err := domainScope.Get(context.Background(), domainScope.Selector{Id: scopeId})

	if err != nil {
		t.Errorf("TestHashDomainXX32 Failed to resolve Scope: %s", err)
	} else {
		domain, derr := domain.Get(context.Background(), domain.Selector{Scope: scope, Id: domainId})

		if derr != nil {
			t.Errorf("TestHashDomainXX32 Failed to resolve Domain: %s", derr)
		} else {
			if base62.Encode(domain.Id()) != "H7B" {
				t.Errorf("TestHashDomainXX32 Failed: domain.Id: expected: '%s' got '%s'", "H7B", base62.Encode(domain.Id()))
			}

			if domain.Name() != "xx32" {
				t.Errorf("TestHashDomainXX32 Failed: domain.Name: expected: '%s' got '%s'", "xx32", domain.Name())

			}

			if base62.Encode(domain.IdRoot()) != "1" {
				t.Errorf("TestHashDomainXX32 Failed: domain.IdRoot: expected: '%s' got '%s'", "1", base62.Encode(domain.IdRoot()))

			}

			if base62.Encode(domain.Scope().Id()) != "1" {
				t.Errorf("TestHashDomainXX32 Failed: domain.Scope.Id: expected: '%s' got '%s'", "1", base62.Encode(domain.Scope().Id()))
			}
		}
	}

	scope, err = domainScope.Get(context.Background(), scopeNameSelector)

	if err != nil {
		t.Errorf("TestHashDomainXX32 Failed to resolve Scope: %s", err)
	} else {
		domain, derr := domain.Get(context.Background(), domainNameSelector)
		if derr != nil {
			t.Errorf("TestHashDomainXX32 Failed to resolve Domain: %s", derr)
		} else {
			if base62.Encode(domain.Id()) != "H7B" {
				t.Errorf("TestHashDomainXX32 Failed: domain.Id: expected: '%s' got '%s'", "H7B", base62.Encode(domain.Id()))
			}

			if domain.Name() != "xx32" {
				t.Errorf("TestHashDomainXX32 Failed: domain.Name: expected: '%s' got '%s'", "xx32", domain.Name())

			}

			if base62.Encode(domain.IdRoot()) != "1" {
				t.Errorf("TestHashDomainXX32 Failed: domain.IdRoot: expected: '%s' got '%s'", "1", base62.Encode(domain.IdRoot()))
			}

			if base62.Encode(domain.Scope().Id()) != "1" {
				t.Errorf("TestHashDomainXX32 Failed: domain.Scope.Id: expected: '%s' got '%s'", "1", base62.Encode(domain.Scope().Id()))
			}
		}
	}
}

func TestHashDomainRipemd160(t *testing.T) {
	scopeId, err := base62.Decode("1")

	if err != nil {
		t.Errorf("TestHashDomainRipemd160 Failed: Can't decode scopeId")
	}

	domainId, err := base62.Decode("2")

	if err != nil {
		t.Errorf("TestHashDomainRipemd160 Failed: Can't decode domainId")
	}

	scopeNameSelector := domainScope.Selector{Name: "global hash", Opts: domainScope.SelectorOpts{
		IgnoreCase:       true,
		IgnoreWhitespace: true}}

	domainNameSelector := domain.Selector{Name: "ripemd160", Opts: domain.SelectorOpts{
		IgnoreCase:       true,
		IgnoreWhitespace: true}}

	scope, err := domainScope.Get(context.Background(), domainScope.Selector{Id: scopeId})

	if err != nil {
		t.Errorf("TestHashDomainRipemd160 Failed to resolve Scope: %s", err)
	} else {
		domain, derr := domain.Get(context.Background(), domain.Selector{Scope: scope, Id: domainId})

		if derr != nil {
			t.Errorf("TestHashDomainRipemd160 Failed to resolve Domain: %s", derr)
		} else {
			if base62.Encode(domain.Id()) != "H7C" {
				t.Errorf("TestHashDomainRipemd160 Failed: domain.Id: expected: '%s' got '%s'", "H7C", base62.Encode(domain.Id()))
			}

			if domain.Name() != "ripemd160" {
				t.Errorf("TestHashDomainRipemd160 Failed: domain.Name: expected: '%s' got '%s'", "ripemd160", domain.Name())

			}

			if base62.Encode(domain.IdRoot()) != "2" {
				t.Errorf("TestHashDomainRipemd160 Failed: domain.IdRoot: expected: '%s' got '%s'", "2", base62.Encode(domain.IdRoot()))

			}

			if base62.Encode(domain.Scope().Id()) != "1" {
				t.Errorf("TestHashDomainRipemd160 Failed: domain.Scope.Id: expected: '%s' got '%s'", "1", base62.Encode(domain.Scope().Id()))
			}
		}
	}

	scope, err = domainScope.Get(context.Background(), scopeNameSelector)

	if err != nil {
		t.Errorf("TestHashDomainRipemd160 Failed to resolve Scope: %s", err)
	} else {
		domain, derr := domain.Get(context.Background(), domainNameSelector)
		if derr != nil {
			t.Errorf("TestHashDomainRipemd160 Failed to resolve Domain: %s", derr)
		} else {
			if base62.Encode(domain.Id()) != "H7C" {
				t.Errorf("TestHashDomainRipemd160 Failed: domain.Id: expected: '%s' got '%s'", "H7C", base62.Encode(domain.Id()))
			}

			if domain.Name() != "ripemd160" {
				t.Errorf("TestHashDomainRipemd160 Failed: domain.Name: expected: '%s' got '%s'", "ripemd160", domain.Name())

			}

			if base62.Encode(domain.IdRoot()) != "2" {
				t.Errorf("TestHashDomainRipemd160 Failed: domain.IdRoot: expected: '%s' got '%s'", "2", base62.Encode(domain.IdRoot()))
			}

			if base62.Encode(domain.Scope().Id()) != "1" {
				t.Errorf("TestHashDomainRipemd160 Failed: domain.Scope.Id: expected: '%s' got '%s'", "1", base62.Encode(domain.Scope().Id()))
			}
		}
	}
}

func TestHashDomainSha256(t *testing.T) {
	scopeId, err := base62.Decode("1")

	if err != nil {
		t.Errorf("TestHashDomainSha256 Failed: Can't decode scopeId")
	}

	domainId, err := base62.Decode("3")

	if err != nil {
		t.Errorf("TestHashDomainSha256 Failed: Can't decode domainId")
	}

	scopeNameSelector := domainScope.Selector{Name: "global hash", Opts: domainScope.SelectorOpts{
		IgnoreCase:       true,
		IgnoreWhitespace: true}}

	domainNameSelector := domain.Selector{Name: "sha256", Opts: domain.SelectorOpts{
		IgnoreCase:       true,
		IgnoreWhitespace: true}}

	scope, err := domainScope.Get(context.Background(), domainScope.Selector{Id: scopeId})

	if err != nil {
		t.Errorf("TestHashDomainSha256 Failed to resolve Scope: %s", err)
	} else {
		domain, derr := domain.Get(context.Background(), domain.Selector{Scope: scope, Id: domainId})

		if derr != nil {
			t.Errorf("TestHashDomainSha256 Failed to resolve Domain: %s", derr)
		} else {
			if base62.Encode(domain.Id()) != "H7D" {
				t.Errorf("TestHashDomainSha256 Failed: domain.Id: expected: '%s' got '%s'", "H7D", base62.Encode(domain.Id()))
			}

			if domain.Name() != "sha256" {
				t.Errorf("TestHashDomainSha256 Failed: domain.Name: expected: '%s' got '%s'", "sha256", domain.Name())

			}

			if base62.Encode(domain.IdRoot()) != "3" {
				t.Errorf("TestHashDomainSha256 Failed: domain.IdRoot: expected: '%s' got '%s'", "3", base62.Encode(domain.IdRoot()))

			}

			if base62.Encode(domain.Scope().Id()) != "1" {
				t.Errorf("TestHashDomainSha256 Failed: domain.Scope.Id: expected: '%s' got '%s'", "1", base62.Encode(domain.Scope().Id()))
			}
		}
	}

	scope, err = domainScope.Get(context.Background(), scopeNameSelector)

	if err != nil {
		t.Errorf("TestHashDomainSha256 Failed to resolve Scope: %s", err)
	} else {
		domain, derr := domain.Get(context.Background(), domainNameSelector)
		if derr != nil {
			t.Errorf("TestHashDomainSha256 Failed to resolve Domain: %s", derr)
		} else {
			if base62.Encode(domain.Id()) != "H7D" {
				t.Errorf("TestHashDomainSha256 Failed: domain.Id: expected: '%s' got '%s'", "H7D", base62.Encode(domain.Id()))
			}

			if domain.Name() != "sha256" {
				t.Errorf("TestHashDomainSha256 Failed: domain.Name: expected: '%s' got '%s'", "sha256", domain.Name())

			}

			if base62.Encode(domain.IdRoot()) != "3" {
				t.Errorf("TestHashDomainSha256 Failed: domain.IdRoot: expected: '%s' got '%s'", "3", base62.Encode(domain.IdRoot()))
			}

			if base62.Encode(domain.Scope().Id()) != "1" {
				t.Errorf("TestHashDomainSha256 Failed: domain.Scope.Id: expected: '%s' got '%s'", "1", base62.Encode(domain.Scope().Id()))
			}
		}
	}
}

/*
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
