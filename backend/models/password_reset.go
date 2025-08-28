package models

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type PasswordResetToken struct {
	Token     string     `orm:"size(64);pk" json:"token"`
	Email     string     `orm:"size(191)" json:"email"`
	ExpiresAt time.Time  `orm:"type(datetime)" json:"expires_at"`
	UsedAt    *time.Time `orm:"null;type(datetime)" json:"used_at"`
	CreatedAt time.Time  `orm:"auto_now_add;type(datetime)" json:"created_at"`
}

func (t *PasswordResetToken) TableName() string { return "password_reset_token" }

func randomResetToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func CreatePasswordResetToken(email string, ttl time.Duration) (*PasswordResetToken, error) {
	tok, err := randomResetToken(16)
	if err != nil {
		return nil, err
	}
	t := &PasswordResetToken{
		Token:     tok,
		Email:     email,
		ExpiresAt: time.Now().Add(ttl),
	}
	o := orm.NewOrm()
	if _, err := o.Insert(t); err != nil {
		return nil, err
	}
	return t, nil
}

func ConsumePasswordResetToken(token string) (string, error) {
	o := orm.NewOrm()
	t := PasswordResetToken{Token: token}
	if err := o.Read(&t); err != nil {
		return "", err
	}
	if t.UsedAt != nil {
		return "", errors.New("token already used")
	}
	if time.Now().After(t.ExpiresAt) {
		return "", errors.New("token expired")
	}
	now := time.Now()
	t.UsedAt = &now
	if _, err := o.Update(&t, "UsedAt"); err != nil {
		return "", err
	}
	return t.Email, nil
}
