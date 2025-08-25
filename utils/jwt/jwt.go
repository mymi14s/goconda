package jwtutil

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"github.com/beego/beego/v2/server/web"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func secretAndIssuer() ([]byte, string, time.Duration) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// try config
		secret = web.AppConfig.DefaultString("jwtsecret", "")
	}
	if secret == "" {
		secret = "dev-secret-please-change"
	}
	issuer := web.AppConfig.DefaultString("jwtissuer", "goconda")
	ttlStr := web.AppConfig.DefaultString("jwtexp", "")
	ttl := 24 * time.Hour
	if ttlStr != "" {
		if d, err := time.ParseDuration(ttlStr); err == nil {
			ttl = d
		}
	}
	return []byte(secret), issuer, ttl
}

func randomJTI() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func Generate(email string) (string, error) {
	sec, iss, ttl := secretAndIssuer()
	now := time.Now()
	claims := Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    iss,
			Subject:   email,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			ID:        randomJTI(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(sec)
}

func Parse(tokenStr string) (*Claims, error) {
	sec, _, _ := secretAndIssuer()
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return sec, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
