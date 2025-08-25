package tests

import (
	"testing"

	"github.com/mymi14s/goconda/utils/validators"
)

func TestEmail(t *testing.T) {
	if err := validators.ValidateEmail("john@example.com"); err != nil {
		t.Fatalf("valid email failed: %v", err)
	}
	if err := validators.ValidateEmail("nope"); err == nil {
		t.Fatal("expected error for invalid email")
	}
}
