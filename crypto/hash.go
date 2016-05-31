package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(str string) []byte {
	tmp := md5.Sum([]byte(str))
	return tmp[:]
}

func Md5String(str string) string {
	return hex.EncodeToString(Md5(str))
}
