package domain_test

import (
	"os"
	"testing"

	"github.com/distributed-vision/go-resources/init/domainscopeinit"
)

func TestMain(m *testing.M) {

	domainscopeinit.Init()

	os.Exit(m.Run())
}
