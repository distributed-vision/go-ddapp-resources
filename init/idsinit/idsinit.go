package idsinit

import (
	"github.com/distributed-vision/go-resources/ids/signature"
	"github.com/distributed-vision/go-resources/init/domainscopeinit"
)

func Init() {
	domainscopeinit.Init()
	signature.Init()
}
