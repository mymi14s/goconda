package hash

import "golang.org/x/crypto/bcrypt"

func HashPassword(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(b), err
}

func CheckPassword(pw, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pw)) == nil
}

func Make(pw string) (string, error) {
	return HashPassword(pw)
}

func Check(pw, hashed string) bool {
	return CheckPassword(pw, hashed)
}
