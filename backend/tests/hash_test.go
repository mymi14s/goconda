package tests

import (
	"testing"

	"github.com/mymi14s/goconda/utils/hash"
)

func TestHashAndCheck(t *testing.T) {
	h, err := hash.HashPassword("secret123")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	if !hash.CheckPassword("secret123", h) {
		t.Fatal("expected password to match")
	}
	if hash.CheckPassword("wrong", h) {
		t.Fatal("expected password to NOT match")
	}
}
