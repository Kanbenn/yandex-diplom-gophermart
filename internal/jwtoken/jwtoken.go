package jwtoken

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UID int
}

const tokenExp = 7 * 24 * time.Hour
const secretKey = "topsecretkey"

func ParseToken(tokenString string) (uid int, err error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return 0, err
	}
	return claims.UID, nil
}

func MakeToken(uid int) (token string, err error) {
	jt := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UID: uid,
	})

	tokenString, err := jt.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
