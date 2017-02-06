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

var decoders = map[encoderType.EncoderType]func(src string) ([]byte, error){
	encoderType.BASE62:    base62.Decode,
	encoderType.BASE36:    base36.Decode,
	encoderType.BASE58:    base58.Decode,
	encoderType.BASE10:    base10.Decode,
	encoderType.BASE64URL: base64url.Decode,
	encoderType.HEX:       hex.DecodeString,
}

func Decode(toDecode string, encoding encoderType.EncoderType) ([]byte, error) {
	decoder := decoders[encoding]

	if decoder != nil {
		return decoder(toDecode)
	}

	if encoding == encoderType.UTF8 {
		return []byte(toDecode), nil
	}

	return nil, errors.New("Unsupported encoding: " + encoding.String())
}
