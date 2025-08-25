package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type User struct {
	Email        string    `orm:"size(191);pk" json:"email"`
	FirstName    string    `orm:"size(100)" json:"first_name"`
	LastName     string    `orm:"size(100)" json:"last_name"`
	PasswordHash string    `orm:"size(255)" json:"-"`
	CreatedAt    time.Time `orm:"auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt    time.Time `orm:"auto_now;type(datetime)" json:"updated_at"`
	IsSuperuser bool      `orm:"default(false)" json:"is_superuser"`
}

func (u *User) TableName() string { return "users" }

func GetUserByEmail(email string) (*User, error) {
	o := orm.NewOrm()
	u := User{Email: email}
	err := o.Read(&u)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
