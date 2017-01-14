package encoding

import (
	"encoding/hex"
	"errors"

	"github.com/distributed-vision/go-resources/encoding/base10"
	"github.com/distributed-vision/go-resources/encoding/base36"
	"github.com/distributed-vision/go-resources/encoding/base58"
	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/encoding/base64url"
)

var encoders = map[string]func(src []byte) string{
	"base62":    base62.Encode,
	"base36":    base36.Encode,
	"base58":    base58.Encode,
	"base10":    base10.Encode,
	"base64url": base64url.Encode,
	"hex":       hex.EncodeToString}

func Encode(toEncode []byte, encoding string) (string, error) {
	encoder := encoders[encoding]

	if encoder != nil {
		return encoder(toEncode), nil
	}

	if encoding == "utf8" {
		return string(toEncode), nil
	}

	return "", errors.New("Unsupported encoding: " + encoding)
}
