package crypto

import (
	"fmt"
	"testing"
	"time"
)

func TestJWT(t *testing.T) {
	cases := []struct {
		name       string
		uid        int64
		roleId     int32
		expireTime time.Time
		hasErr     bool
	}{
		{"ok", 1, 2, time.Now().Add(3 * time.Hour), false},
		{"timeout", 1, 2, time.Now().Add(-3 * time.Hour), true},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%s", c.name), func(t *testing.T) {
			token, err := GenerateToken(c.uid, c.roleId, c.expireTime)
			if err != nil {
				t.Fatal(err)
			} else {
				t.Log("token: ", token)
			}

			rawToken, err := ParseToken(token)
			if (err != nil) != c.hasErr {
				t.Fatal(err)
			} else {
				t.Log("raw token: ", rawToken)
			}
		})
	}
}
