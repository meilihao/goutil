package rand

import (
	crand "crypto/rand"
	"encoding/hex"
)

func RandString(n int) (string, error) {
	b := make([]byte, n)
	_, err := crand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
