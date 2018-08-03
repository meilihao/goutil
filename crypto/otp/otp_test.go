package otp

import (
	"testing"
	"time"
)

func TestGenerateSecret(t *testing.T) {
	secret1 := GenerateSecret(1)
	secret2 := GenerateSecret(10)
	secret3 := GenerateSecret(13)

	if !(len(secret1) == len(secret2) && len(secret2) < len(secret3)) {
		t.Error("bad GenerateSecret")
	} else {
		t.Log(secret1, secret2, secret3)
	}
}

func TestGenerateTOTPAuth(t *testing.T) {
	secret := "QTGS4FA5CWKPESL762TWAWYHS45NICLZ"
	user := "john@example.com"
	issuer := "ACME Co"
	target := "otpauth://totp/ACME%20Co:john@example.com?secret=QTGS4FA5CWKPESL762TWAWYHS45NICLZ&issuer=ACME%20Co"

	uri := GenerateTOTPAuth(secret, user, issuer)

	if uri != target {
		t.Errorf("bad GenerateTOTPAuth(%s,%s)", uri, target)
	} else {
		t.Log(uri)
	}
}

// https://authenticator.ppl.family/
func TestGenerateTOTP(t *testing.T) {
	secret := "WS6JDWWTG76UNFANLJEPM==="

	num, remaining := GenerateTOTP(secret, time.Now().Unix(), 6, "")

	if num == "" {
		t.Errorf("bad GenerateTOTP(%s,%d)", num, remaining)
	} else {
		t.Log(num, remaining)
	}
}
