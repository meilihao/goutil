// https://www.jianshu.com/p/a7b900e8e50a
// http://blog.gaoyuan.xyz/2017/01/05/2fa-a-programmers-perspective/
/*
共享密钥使用base32而非base64的原因如下：

base32编码的字符串，包含了大写英文字母和数字2-7。不会因字体显示问题，把1，8，0和’I’,‘B’, ‘O’混淆，更利于输入。
base32编码的字符串，出现在url中时，可以不用进行url编码处理（encode），便于直接使用生成二维码的web服务。
*/
package otp

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"hash"
	"net/url"
)

var (
	T0           int64 = 0
	X            int64 = 30
	Digits_Power       = []uint32{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000}
)

// google-authenticator is GenerateTOTP(secret, time.Now().Unix(), 6, "") with (T0,X) = (0,30), len(secret)=16
func GenerateTOTP(secret string, timestamp int64, count int, hashName string) (string, int64) {
	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		fmt.Println(err)
		return "", -1
	}

	var h func() hash.Hash

	switch hashName {
	case "sha256":
		h = sha256.New
	case "sha512":
		h = sha512.New
	default:
		h = sha1.New
	}

	if count > len(Digits_Power)-1 {
		count = 6
	}

	winSecond := (timestamp - T0) / X
	remaining := X - (timestamp % X)

	result := oneTimePassword(key, winSecond, count, h)
	if result == "" {
		remaining = -1
	}

	return result, remaining
}

func oneTimePassword(key []byte, value int64, count int, h func() hash.Hash) string {
	hmacHash := hmac.New(h, key)
	if err := binary.Write(hmacHash, binary.BigEndian, value); err != nil {
		return ""
	}
	tmp := hmacHash.Sum(nil)

	offset := tmp[len(tmp)-1] & 0x0f
	num := toUint32(tmp[offset : offset+4])
	num = num % Digits_Power[count]

	result := fmt.Sprintf("%d", num)
	for len(result) < count {
		result = "0" + result
	}

	return result
}

func toUint32(bytes []byte) uint32 {
	return (uint32(bytes[0]&0x7f) << 24) |
		(uint32(bytes[1]) << 16) |
		(uint32(bytes[2]) << 8) |
		uint32(bytes[3])
}

// min is 10
func GenerateSecret(num int) string {
	if num < 10 {
		num = 10
	}

	var i int
	var err error
	buf := make([]byte, num)

	for {
		if _, err = rand.Read(buf); err != nil {
			i++

			if i > 10 {
				panic("invalid random")
			}

			continue
		}

		break
	}

	return base32.StdEncoding.EncodeToString(buf)
}

// https://github.com/google/google-authenticator/wiki/Key-Uri-Format
func GenerateTOTPAuth(secret, user, issuer string) string {
	return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s",
		url.PathEscape(issuer), // not QueryEscape()
		url.PathEscape(user),
		url.PathEscape(secret),
		url.PathEscape(issuer),
	)
}
