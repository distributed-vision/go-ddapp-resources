package mappingResolver

import (
	"time"

	"github.com/distributed-vision/go-resources/ids"
)

var MaxTime = time.Unix(1<<63-62135596801, 999999999)
var MinTime = time.Time{}

func ResolveMapping(id ids.Identifier, to ids.IdentityDomain, at time.Time) (chan ids.Identifier, chan error) {
	return nil, nil
}

func AddMapping(from ids.Identifier, to ids.Identifier, after *time.Time, before *time.Time) error {
	return nil
}
