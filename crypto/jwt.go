package crypto

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("123456")

type Claims struct {
	Uid    int64 `json:"uid"`
	RoleId int32 `json:"role_id"`
	jwt.StandardClaims
}

func GenerateToken(uid int64, roleId int32, expireTime time.Time) (string, error) {
	claims := Claims{
		uid,
		roleId,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			//Issuer : "chen",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
