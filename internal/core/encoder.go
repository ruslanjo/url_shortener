package core

import (
	"crypto/md5"
	"encoding/base64"
)

func HashString(s string) []byte {
	h := md5.New()
	h.Write([]byte(s))
	return h.Sum(nil)
}

func EncodeHash(hash []byte) string {
	encoded := base64.URLEncoding.EncodeToString(hash)
	return encoded
}
