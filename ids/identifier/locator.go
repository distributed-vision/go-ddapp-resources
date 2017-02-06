package identifier

import "github.com/distributed-vision/go-resources/ids"

type locator struct {
	*identifier
	entity interface{}
}

func NewLocator(id ids.Identifier) ids.Locator {
	return nil
}
