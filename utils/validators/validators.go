package validators

import (
	"errors"
	"net/mail"
	"strings"
)

func ValidateEmail(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return errors.New("email is required")
	}
	_, err := mail.ParseAddress(s)
	return err
}

func RequireNonEmpty(fields map[string]string) error {
	for k, v := range fields {
		if strings.TrimSpace(v) == "" {
			return errors.New(k + " is required")
		}
	}
	return nil
}

func IsEmailValid(s string) bool {
	return ValidateEmail(s) == nil
}
