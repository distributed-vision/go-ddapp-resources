package base62

import "github.com/distributed-vision/go-resources/encoding/basex"

var base62 = basex.NewEncoder("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

func Decode(toDecode string) ([]byte, error) {
	return base62.Decode(toDecode)
}

func Encode(toEncode []byte) string {
	return base62.Encode(toEncode)
}
