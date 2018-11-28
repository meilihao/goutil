package time

import (
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/blake2b"
)

// use blake2b
type TimeXB struct {
	key      []byte
	h        func([]byte) (hash.Hash, error)
	diff     int64
	trim2Len int // 截留长度
}

// h : blake2b
// size = len(key), size: 32, 48, 64
// 8<=trimLen<len(hash) ? hash[:trimLen] : hash
func NewTimeXB(size int, key string, diff int64, trimLen int) (*TimeXB, error) {
	var h func([]byte) (hash.Hash, error)

	switch size {
	case blake2b.Size:
		size = blake2b.Size
		h = blake2b.New512
	case blake2b.Size384:
		size = blake2b.Size384
		h = blake2b.New384
	default:
		size = blake2b.Size256
		h = blake2b.New256
	}

	if len(key) != size {
		return nil, errors.New("blake2b: size need equal len(key)")
	}

	return &TimeXB{
		key:      []byte(key),
		h:        h,
		diff:     diff,
		trim2Len: trimLen,
	}, nil
}

func (x *TimeXB) Generate(s string) string {
	return x.generate(s, time.Now().Unix())
}

func (x *TimeXB) generate(s string, now int64) string {
	tmp := fmt.Sprintf("%s-%d", s, now)

	h, _ := x.h(x.key)
	h.Write([]byte(tmp))

	result := hex.EncodeToString(h.Sum(nil))
	if x.trim2Len >= 8 && x.trim2Len < len(result) {
		result = result[:x.trim2Len]
	}

	return tmp + "-" + result
}

func (x *TimeXB) Parse(s string) (string, int64, error) {
	tmp := strings.Split(s, "-")

	if len(tmp) != 3 {
		return "", 0, errors.New("No TimeXB")
	}

	now, _ := strconv.ParseInt(tmp[1], 10, 64)

	if s != x.generate(tmp[0], now) {
		return "", 0, errors.New("Invalid TimeXB")
	}

	return tmp[0], now, nil
}
