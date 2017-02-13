package domainScope_test

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/ids/domainScope"
	"github.com/distributed-vision/go-resources/ids/domainScopeFormat"
	"github.com/distributed-vision/go-resources/ids/domainScopeVisibility"
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
	cres, cerr := domainScope.Resolve(context.Background(), domainScope.Selector{Id: scopeId})

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
	scopeNameSelector := domainScope.Selector{Name: "untyped",
		Opts: domainScope.SelectorOpts{
			IgnoreCase:       true,
			IgnoreWhitespace: true}}

	scope, err := domainScope.Get(context.Background(), domainScope.Selector{Id: scopeId})

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

			if scope.Visibility() != domainScopeVisibility.PUBLIC {
				t.Errorf("ResolveUntypedScope Failed: Visibility: expected: '%v' got '%v'", domainScopeVisibility.PUBLIC, scope.Visibility())

			}

			if scope.Format() != domainScopeFormat.FIXED {
				t.Errorf("ResolveUntypedScope Failed: Format: expected: '%v' got '%v'", domainScopeFormat.FIXED, scope.Format())
			}
		}
	}

	scope, err = domainScope.Get(context.Background(), scopeNameSelector)

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

			if scope.Visibility() != domainScopeVisibility.PUBLIC {
				t.Errorf("ResolveUntypedScope Failed: Visibility: expected: '%v got '%v'", domainScopeVisibility.PUBLIC, scope.Visibility())

			}

			if scope.Format() != domainScopeFormat.FIXED {
				t.Errorf("ResolveUntypedScope Failed: Format: expected: '%v' got '%v'", domainScopeFormat.FIXED, scope.Format())
			}
		}
	}
}
