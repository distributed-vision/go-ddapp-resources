package publictypeid

import (
	"github.com/distributed-vision/go-resources/encoding/encoderType"
	"github.com/distributed-vision/go-resources/ids/domain"
	"github.com/distributed-vision/go-resources/version/versionType"
)

var ResolverDomain = domain.MustDecodeId(encoderType.BASE62, "T", "0", uint32(0), uint(0), versionType.SEMANTIC)
