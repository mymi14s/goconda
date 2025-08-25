package models

import (
    "crypto/rand"
    "encoding/hex"
    "errors"
    "time"

    "github.com/beego/beego/v2/client/orm"
)

type EmailVerificationToken struct {
    Token     string    `orm:"size(64);pk" json:"token"`
    Email     string    `orm:"size(191)" json:"email"`
    ExpiresAt time.Time `orm:"type(datetime)" json:"expires_at"`
    UsedAt    *time.Time `orm:"null;type(datetime)" json:"used_at"`
    CreatedAt time.Time `orm:"auto_now_add;type(datetime)" json:"created_at"`
}

func (t *EmailVerificationToken) TableName() string { return "email_verification_tokens" }

type VerifiedUser struct {
    Email      string    `orm:"size(191);pk" json:"email"`
    VerifiedAt time.Time `orm:"type(datetime)" json:"verified_at"`
}

func (v *VerifiedUser) TableName() string { return "verified_users" }

func randomToken(n int) (string, error) {
    b := make([]byte, n)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return hex.EncodeToString(b), nil
}

func CreateVerificationToken(email string, ttl time.Duration) (*EmailVerificationToken, error) {
    tok, err := randomToken(16)
    if err != nil {
        return nil, err
    }
    t := &EmailVerificationToken{
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

func ConsumeVerificationToken(token string) (string, error) {
    o := orm.NewOrm()
    t := EmailVerificationToken{Token: token}
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

func MarkUserVerified(email string) error {
    o := orm.NewOrm()
    v := &VerifiedUser{Email: email, VerifiedAt: time.Now()}
    // Upsert behavior: try insert; if exists, update time
    if _, err := o.Insert(v); err != nil {
        // try update
        if _, err2 := o.Update(v, "VerifiedAt"); err2 != nil {
            return err
        }
    }
    return nil
}

func IsUserVerified(email string) (bool, error) {
    o := orm.NewOrm()
    v := VerifiedUser{Email: email}
    if err := o.Read(&v); err != nil {
        if err == orm.ErrNoRows {
            return false, nil
        }
        return false, err
    }
    return true, nil
}
