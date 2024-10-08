package web

import (
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var chars []rune = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

func RandStringRunes(n int) string {
	length := len(chars)
	b := make([]rune, n)
	for i := range b {
		b[i] = chars[rand.Intn(length)]
	}
	return string(b)
}

func GenerateRefreshToken() string {
	return RandStringRunes(32)
}

func GenerateAccessToken(secret string) string {
	expirationTime := time.Now().Add(time.Hour).UTC()
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

func parseAccessToken(tokenString string, secret string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
}

func validateToken(token *jwt.Token) bool {
	validator := jwt.NewValidator()
	err := validator.Validate(token.Claims)
	if err != nil {
		return false
	}
	return token.Valid
}

func ParseAndValidateToken(tokenString string, secret string) bool {
	token, err := parseAccessToken(tokenString, secret)
	if err != nil {
		return false
	}
	return validateToken(token)
}
