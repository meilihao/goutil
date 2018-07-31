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
	key  []byte
	h    func() hash.Hash
	diff int64
}

func NewTimeX(h func() hash.Hash, key string, diff int64) *TimeX {
	return &TimeX{
		key:  []byte(key),
		h:    h,
		diff: diff,
	}
}

func (x *TimeX) Generate(s string) string {
	return x.generate(s, time.Now().Unix())
}

func (x *TimeX) generate(s string, now int64) string {
	tmp := fmt.Sprintf("%s-%d", s, now)

	h := hmac.New(x.h, x.key)
	h.Write([]byte(tmp))

	return tmp + "-" + hex.EncodeToString(h.Sum(nil))
}

func (x *TimeX) Parse(s string) (string, error) {
	tmp := strings.Split(s, "-")

	if len(tmp) != 3 {
		return "", errors.New("no TimeX")
	}

	now, _ := strconv.ParseInt(tmp[1], 10, 64)

	if s != x.generate(tmp[0], now) {
		return "", errors.New("invalid TimeX")
	}

	current := time.Now().Unix()
	if !(now > 0 && (current-now <= x.diff || now-current <= x.diff)) {
		return "", errors.New("invalid timestamp")
	}

	return tmp[0], nil
}
