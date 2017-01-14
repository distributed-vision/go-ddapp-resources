package domains

import (
	"os"
	"testing"

	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/resolvers/domainScopeResolver"
)

func setup() {
	if os.Getenv("DV_DOMAIN_PATH") == "" {
		os.Setenv("DV_DOMAIN_PATH", "../../../../distributed-vision/scope")
	}
}

func TestResolveUntypedScope(t *testing.T) {
	setup()
	scopeId, _ := base62.Decode("0")
	/*scopeNameSelector := domainScopeResolver.Selector{Name: "untyped",
	Opts: domainScopeResolver.SelectorOpts{
		IgnoreCase:       true,
		IgnoreWhitespace: true}}
	*/
	cres, cerr := domainScopeResolver.Resolve(
		domainScopeResolver.Selector{ScopeId: scopeId})

	if cres == nil || cerr == nil {
		t.Errorf("ResolveUntypedScope Failed: domainScopeResolver.Resolve channels are undefined")
		return
	}

	select {
	case scope := <-cres:
		if scope == nil {
			t.Errorf("ResolveUntypedScope Failed: scope is undefined")
		}
		/*
			if !bytes.Equal(scope.Id(), scopeId) {
				t.Errorf("ResolveUntypedScope Failed: Id: expected: '%s' got '%s'", scopeId, scope.Id())
			}

			if scope.Name() != "Untyped" {
				t.Errorf("ResolveUntypedScope Failed: Name: expected: '%s' got '%s'", "Untyped", scope.Name())
			}

			if scope.Visibility() != domainScopeVisibility.PUBLIC {
				t.Errorf("ResolveUntypedScope Failed: Name: expected: '%v' got '%v'", domainScopeVisibility.PUBLIC, scope.Visibility())

			}

			if scope.Format() != domainScopeFormat.FIXED {
				t.Errorf("ResolveUntypedScope Failed: Name: expected: '%v' got '%v'", domainScopeFormat.FIXED, scope.Format())
			}

			cres, cerr = domainScopeResolver.Resolve(scopeNameSelector)

			select {
			case scope := <-cres:
				if !bytes.Equal(scope.Id(), scopeId) {
					t.Errorf("ResolveUntypedScope Failed: Id: expected: '%s' got '%s'", scopeId, scope.Id())
				}

				if scope.Name() != "Untyped" {
					t.Errorf("ResolveUntypedScope Failed: Name: expected: '%s' got '%s'", "Untyped", scope.Name())
				}

				if scope.Visibility() != domainScopeVisibility.PUBLIC {
					t.Errorf("ResolveUntypedScope Failed: Name: expected: '%v got '%v'", domainScopeVisibility.PUBLIC, scope.Visibility())

				}

				if scope.Format() != domainScopeFormat.FIXED {
					t.Errorf("ResolveUntypedScope Failed: Name: expected: '%v' got '%v'", domainScopeFormat.FIXED, scope.Format())
				}

				break
			case err := <-cerr:
				t.Errorf("ResolveUntypedScope Failed With Error: %s", err)
				break
			}*/
		break
	case err := <-cerr:
		t.Errorf("ResolveUntypedScope Failed With Error: %s", err)
		break
	}

}
