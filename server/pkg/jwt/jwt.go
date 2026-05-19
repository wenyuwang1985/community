package jwt

import (
	"fmt"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    int64  `json:"user_id"`
	TokenType string `json:"token_type"`
	gojwt.RegisteredClaims
}

func GenerateAccessToken(userID int64, secret string, expireSeconds int) (string, error) {
	return generateToken(userID, secret, expireSeconds, "access")
}

func GenerateRefreshToken(userID int64, secret string, expireSeconds int) (string, error) {
	return generateToken(userID, secret, expireSeconds, "refresh")
}

func generateToken(userID int64, secret string, expireSeconds int, tokenType string) (string, error) {
	claims := Claims{
		UserID:    userID,
		TokenType: tokenType,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(time.Duration(expireSeconds) * time.Second)),
			IssuedAt:  gojwt.NewNumericDate(time.Now()),
		},
	}

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenStr string, secret string) (*Claims, error) {
	token, err := gojwt.ParseWithClaims(tokenStr, &Claims{}, func(t *gojwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*gojwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
