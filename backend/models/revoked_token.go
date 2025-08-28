// models/revoked_token.go
package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type RevokedToken struct {
	JTI       string    `orm:"pk;size(191);column(jti)" json:"jti"`
	ExpiresAt time.Time `orm:"type(datetime);index" json:"expires_at"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)" json:"created_at"`
}

func (r *RevokedToken) TableName() string { return "revoked_token" }

func RevokeToken(jti string, exp time.Time) error {
	if jti == "" {
		return nil
	}
	o := orm.NewOrm()
	rt := &RevokedToken{
		JTI:       jti,
		ExpiresAt: exp,
	}
	// Insert or ignore if exists
	_, err := o.Insert(rt)
	if err != nil {
		// try update
		_, _ = o.Update(rt)
	}
	return nil
}

func IsTokenRevoked(jti string) (bool, error) {
	if jti == "" {
		return false, nil
	}
	o := orm.NewOrm()
	rt := RevokedToken{JTI: jti}
	err := o.Read(&rt)
	if err == orm.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}
