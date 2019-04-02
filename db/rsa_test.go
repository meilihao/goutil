package db

import (
	"testing"

	"github.com/wenzhenxi/gorsa"
)

const (
	rsaPublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsF9X3iIjh703Q6BDNOUK
...
yQIDAQAB
-----END PUBLIC KEY-----`
	rsaPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAsF9X3iIjh703Q6BDNOUKeNSJxX7Yo4OfK9Y6h/dKVB6LE/IN
...
316ambqL4FMYhlEO0xwo66tkHMq5sTg2tOQVgD5UY9HazoQa0yJditM=
-----END RSA PRIVATE KEY-----`
)

func TestRSA(t *testing.T) {
	plaintext := "Lorem Ipsum"

	prienctypt, err := gorsa.PriKeyEncrypt(plaintext, rsaPrivateKey)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(prienctypt)
	}

	pubdecrypt, err := gorsa.PublicDecrypt(prienctypt, rsaPublicKey)
	if err != nil {
		t.Fatal(err)
	}
	if pubdecrypt != plaintext {
		t.Fatal(`解密失败`)
	}
}
