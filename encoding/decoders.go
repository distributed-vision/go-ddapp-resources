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

var decoders = map[string]func(src string) ([]byte, error){
	"base62":    base62.Decode,
	"base36":    base36.Decode,
	"base58":    base58.Decode,
	"base10":    base10.Decode,
	"base64url": base64url.Decode,
	"hex":       hex.DecodeString}

func Decode(toDecode string, encoding string) ([]byte, error) {
	decoder := decoders[encoding]

	if decoder != nil {
		return decoder(toDecode)
	}

	if encoding == "utf8" {
		return []byte(toDecode), nil
	}

	return nil, errors.New("Unsupported encoding: " + encoding)
}
