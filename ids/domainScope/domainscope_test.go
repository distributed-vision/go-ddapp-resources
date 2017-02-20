package domainscope_test

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/encoding/encodertype"
	"github.com/distributed-vision/go-resources/ids/domain"
	"github.com/distributed-vision/go-resources/ids/domainscope"
	"github.com/distributed-vision/go-resources/ids/domainscopeformat"
	"github.com/distributed-vision/go-resources/ids/domainscopevisibility"
	"github.com/distributed-vision/go-resources/init/domainscopeinit"
)

func TestMain(m *testing.M) {

	GetUninitialisedResolver()
	domainscopeinit.Init()

	os.Exit(m.Run())
}

func ignore(value interface{}) {
}

func GetUninitialisedResolver() {

	scopeId, _ := base62.Decode("0")
	cres, cerr := domainscope.Resolve(context.Background(), domainscope.Selector{Id: scopeId})

	if cres == nil || cerr == nil {
		panic("GetUninitialisedResolver Failed: domainScopeResolver.Resolve channels are undefined")
	}

	select {
	case err := <-cerr:
		ignore(err)
	case scope := <-cres:
		if scope != nil {
			panic("GetUninitialisedResolver Failed: resolution should fail")
		}
	}
}

func TestGetUntypedScope(t *testing.T) {
	scopeId, _ := base62.Decode("0")
	scopeNameSelector := domainscope.Selector{Name: "untyped",
		Opts: domainscope.SelectorOpts{
			IgnoreCase:       true,
			IgnoreWhitespace: true}}

	scope, err := domainscope.Get(context.Background(), domainscope.Selector{Id: scopeId})

	if err != nil {
		t.Errorf("ResolveUntypedScope Failed With Error: %s", err)
	} else {
		if scope == nil {
			t.Errorf("ResolveUntypedScope Failed: scope is undefined")
		} else {
			if !bytes.Equal(scope.Id(), scopeId) {
				t.Errorf("ResolveUntypedScope Failed: Id: expected: '%v' got '%v'", scopeId, scope.Id())
			}

			if scope.Name() != "Untyped" {
				t.Errorf("ResolveUntypedScope Failed: Name: expected: '%s' got '%s'", "Untyped", scope.Name())
			}

			if scope.Visibility() != domainscopevisibility.PUBLIC {
				t.Errorf("ResolveUntypedScope Failed: Visibility: expected: '%v' got '%v'", domainscopevisibility.PUBLIC, scope.Visibility())

			}

			if scope.Format() != domainscopeformat.FIXED {
				t.Errorf("ResolveUntypedScope Failed: Format: expected: '%v' got '%v'", domainscopeformat.FIXED, scope.Format())
			}
		}
	}

	scope, err = domainscope.Get(context.Background(), scopeNameSelector)

	if err != nil {
		t.Errorf("ResolveUntypedScope Failed With Error: %s", err)
	} else {
		if scope == nil {
			t.Errorf("ResolveUntypedScope Failed: scope is undefined")
		} else {
			if !bytes.Equal(scope.Id(), scopeId) {
				t.Errorf("ResolveUntypedScope Failed: Id: expected: '%v' got '%v'", scopeId, scope.Id())
			}

			if scope.Name() != "Untyped" {
				t.Errorf("ResolveUntypedScope Failed: Name: expected: '%s' got '%s'", "Untyped", scope.Name())
			}

			if scope.Visibility() != domainscopevisibility.PUBLIC {
				t.Errorf("ResolveUntypedScope Failed: Visibility: expected: '%v got '%v'", domainscopevisibility.PUBLIC, scope.Visibility())

			}

			if scope.Format() != domainscopeformat.FIXED {
				t.Errorf("ResolveUntypedScope Failed: Format: expected: '%v' got '%v'", domainscopeformat.FIXED, scope.Format())
			}
		}
	}
}

func TestScopeIdFormatting(t *testing.T) {
	scopeId, err := domainscope.DecodeId(encodertype.BASE62, "0", "123")

	if err != nil {
		t.Errorf("TestScopeIdFormatting Failed With Error: %s", err)
	}

	if base62.Encode(scopeId) != "1AfsQd" {
		t.Errorf("TestScopeIdFormatting unexpected scopeId: expected %s, got: %s", "1AfsQd", base62.Encode(scopeId))
	}

	if scopeId[0] != (1 << 6) {
		t.Errorf("TestScopeIdFormatting base encoding failed: expected %v, got: %v", (1 << 6), scopeId[0])
	}

	if scopeId[1] != 2 {
		t.Errorf("TestScopeIdFormatting encoding length incorrect: expected %v, got: %v", 2, scopeId[1])
	}

	if domain.ScopeLength(scopeId) != 4 {
		t.Errorf("TestScopeIdFormatting scope length incorrect: expected %v, got: %v", 4, domain.ScopeLength(scopeId))
	}

	if !bytes.Equal(scopeId[2:2+scopeId[1]], base62.MustDecode("123")) {
		t.Errorf("TestScopeIdFormatting extension encoding failed: expected %v, got: %v", base62.MustDecode("123"), scopeId[2:2+scopeId[1]])
	}

}
