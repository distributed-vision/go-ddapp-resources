package idsinit

import (
	"github.com/distributed-vision/go-resources/ids/signaturedomain"
	"github.com/distributed-vision/go-resources/init/domainscopeinit"
)

func Init() {
	domainscopeinit.Init()
	signaturedomain.Init()
}
