package mappings

import (
	"context"
	"time"

	"github.com/distributed-vision/go-resources/ids"
)

var MaxTime = time.Unix(1<<63-62135596801, 999999999)
var MinTime = time.Time{}

type mapping struct {
	FromId   ids.Identifier
	ToDomain ids.IdentityDomain
	From     time.Time
	To       time.Time
	ToId     ids.Identifier
}

type Selector struct {
	From ids.Identifier
	To   ids.IdentityDomain
	At   time.Time
}

func (this *Selector) Test(candidate interface{}) bool {
	mapping, ok := candidate.(*mapping)

	if ok {
		return this.At.After(mapping.From) && this.At.Before(mapping.To) &&
			this.From.Equals(mapping.FromId) && this.To.Equals(mapping.ToDomain)
	}

	return false
}

func (this *Selector) Key() interface{} {
	return this
}

func Get(resolutionContext context.Context, selector Selector) (ids.Identifier, error) {
	return nil, nil
}

func Resolve(resolutionContext context.Context, selector Selector) (chan ids.Identifier, chan error) {
	return nil, nil
}

func Add(from ids.Identifier, to ids.Identifier, after *time.Time, before *time.Time) error {
	return nil
}
