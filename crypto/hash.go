package crypto

import (
	"encoding/hex"
	"hash"
)

func HashBytes(h hash.Hash, str string) []byte {
	h.Write([]byte(str))
	return h.Sum(nil)
}

func HashString(h hash.Hash, str string) string {
	return hex.EncodeToString(HashBytes(h, str))
}
