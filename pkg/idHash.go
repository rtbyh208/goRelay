package pkg

import (
	"crypto/md5"
	"encoding/base64"
)

func IDHash(id string) string {
	hash := md5.Sum([]byte(id))
	return base64.StdEncoding.EncodeToString(hash[:])
}
