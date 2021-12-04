package data

import (
	"encoding/base64"
)

func Encode(code string) string {
	return base64.StdEncoding.EncodeToString([]byte(code))
}

func Decode(code string) (str []byte, err error) {
	str, err = base64.StdEncoding.DecodeString(code)
	return
}
