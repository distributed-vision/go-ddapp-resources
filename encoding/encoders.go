package encoding

import (
	"encoding/hex"
	"errors"

	"github.com/distributed-vision/go-resources/encoding/base10"
	"github.com/distributed-vision/go-resources/encoding/base36"
	"github.com/distributed-vision/go-resources/encoding/base58"
	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/encoding/base64url"
	"github.com/distributed-vision/go-resources/encoding/encoderType"
)

var encoders = map[encoderType.EncoderType]func(src []byte) string{
	encoderType.BASE62:    base62.Encode,
	encoderType.BASE36:    base36.Encode,
	encoderType.BASE58:    base58.Encode,
	encoderType.BASE10:    base10.Encode,
	encoderType.BASE64URL: base64url.Encode,
	encoderType.HEX:       hex.EncodeToString}

func Encode(toEncode []byte, encoding encoderType.EncoderType) (string, error) {
	encoder := encoders[encoding]

	if encoder != nil {
		return encoder(toEncode), nil
	}

	if encoding == encoderType.UTF8 {
		return string(toEncode), nil
	}

	return "", errors.New("Unsupported encoding: " + encoding.String())
}
