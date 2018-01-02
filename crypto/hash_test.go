package crypto

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"hash"
	"testing"

	"golang.org/x/crypto/sha3"
)

func TestHash(t *testing.T) {
	input := "123456"

	cases := []struct {
		h    hash.Hash
		name string
		str  string
		want string
	}{
		{md5.New(), "md5", input, "e10adc3949ba59abbe56e057f20f883e"},
		{sha256.New(), "sha256", input, "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92"},
		{sha3.New224(), "sha3-224", input, "6be790258b73da9441099c4cb6aeec1f0c883152dd74e7581b70a648"},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%s", c.name), func(t *testing.T) {
			bs := HashBytes(c.h, c.str)
			if len(bs) == 0 {
				t.Fatalf("%s can't sum.", c.name)
			}

			c.h.Reset() // need reset,because use the same hash func.

			s := HashString(c.h, c.str)
			if s != c.want {
				t.Errorf("%s(\"%s\") got %s, want %s", c.name, c.str, s, c.want)
			}
		})
	}
}
