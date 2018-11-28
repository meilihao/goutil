package time

import (
	"crypto/hmac"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"strconv"
	"strings"
	"time"
)

type TimeX struct {
	key      []byte
	h        func() hash.Hash
	trim2Len int
}

// h: md5, shaX
// 8<=trimLen<len(hash) ? hash[:trimLen] : hash
func NewTimeX(h func() hash.Hash, key string, trim2Len int) *TimeX {
	return &TimeX{
		key:      []byte(key),
		h:        h,
		trim2Len: trim2Len,
	}
}

func (x *TimeX) Generate(s string) string {
	return x.generate(s, time.Now().Unix())
}

func (x *TimeX) generate(s string, now int64) string {
	tmp := fmt.Sprintf("%s-%d", s, now)

	h := hmac.New(x.h, x.key)
	h.Write([]byte(tmp))

	result := hex.EncodeToString(h.Sum(nil))
	if x.trim2Len >= 8 && x.trim2Len < len(result) {
		result = result[:x.trim2Len]
	}

	return tmp + "-" + result
}

func (x *TimeX) Parse(s string) (string, int64, error) {
	tmp := strings.Split(s, "-")

	if len(tmp) != 3 {
		return "", 0, errors.New("No TimeX")
	}

	now, _ := strconv.ParseInt(tmp[1], 10, 64)

	if s != x.generate(tmp[0], now) {
		return "", 0, errors.New("Invalid TimeX")
	}

	return tmp[0], now, nil
}
