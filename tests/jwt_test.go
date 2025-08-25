package tests

import (
	"testing"

	jwtutil "github.com/mymi14s/goconda/utils/jwt"
)

func TestJWTRoundTrip(t *testing.T) {
	token, err := jwtutil.Generate("user@example.com")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	claims, err := jwtutil.Parse(token)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if claims.Email != "user@example.com" {
		t.Fatalf("claims mismatch: %+v", claims)
	}
}
