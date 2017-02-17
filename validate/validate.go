package validate

import (
    "regexp"
)

var (
    // --- string ---
    emailReg=regexp.MustCompile("^[a-zA-Z0-9.!#$%&’*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$")
)

func IsEmail(s string) bool{
    return emailReg.MatchString(s)
}