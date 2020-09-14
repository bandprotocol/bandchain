package executor

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GetSingedToken(jwtSecretKey string) string {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Minute * 10).Unix()})
	tokenString, _ := token.SignedString([]byte(jwtSecretKey))
	return tokenString
}
