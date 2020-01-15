package crypto

import "testing"

func TestPassword(t *testing.T) {
	password := "123456"

	hp, _ := HashPassword(password)
	if hp == "" {
		t.Fatal("invalid password")
	}

	t.Log("hashed password: ", hp)
	if !CheckPasswordHash(password, hp) {
		t.Fatal("invalid password check")
	}
}
