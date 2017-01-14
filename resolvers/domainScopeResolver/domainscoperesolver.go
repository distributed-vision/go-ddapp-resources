package domainScopeResolver

import (
	"os"

	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/resolvers"
)

type SelectorOpts struct {
	IgnoreCase       bool
	IgnoreWhitespace bool
}

type Selector struct {
	ScopeId []byte
	Name    string
	Opts    SelectorOpts
}

func (this *Selector) Test(candidate interface{}) bool {

	return true
}

var resolver fileResolver.Resolver

func Resolve(selector resolvers.Selector) (chan ids.DomainScope, chan error) {
	scopePath := os.Getenv("DOMAIN_SCOPE_PATH")

	if resolver == nil {
		resolver = fileResolver.NewFileResolver("scopeinfo.json", fileResolver.Opts{
			entityFromJSON: fromJSON,
			paths:          scopePath})
	}

	cresOut := make(chan ids.DomainScope)
	cerrOut := make(chan error)

	cresIn, cerrIn := resolver.Resolve(selector)

	go func() {
		select {
		case resIn := <-cresIn:
			cresOut <- ids.DomainScope(resIn)
			break
		case errIn := <-cerrIn:
			cerrOut <- errIn
			break
		}

		close(cresIn)
		close(cerrIn)
	}()

	return cresOut, cerrOut
}

func fromJSON(json map[string]interface{}, opts map[string]interface{}) interface{} {
	return nil
}
