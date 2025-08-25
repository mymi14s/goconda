package models

import (
	"errors"

	"github.com/beego/beego/v2/client/orm"
)

type Role struct {
	Name string `orm:"size(100);pk" json:"name"`
}

func (r *Role) TableName() string { return "roles" }

type UserRole struct {
	ID    int64  `orm:"auto;column(id)" json:"id"`
	Email string `orm:"size(191)" json:"email"`
	Role  string `orm:"size(100)" json:"role"`
}

func (ur *UserRole) TableName() string { return "user_roles" }

type Permission struct {
	ID       int64  `orm:"auto;column(id)" json:"id"`
	Role     string `orm:"size(100)" json:"role"`
	Resource string `orm:"size(191)" json:"resource"`
	Action   string `orm:"size(50)" json:"action"` // read, create, update, delete
}

func (p *Permission) TableName() string { return "permissions" }

func EnsureRole(name string) error {
	o := orm.NewOrm()
	if _, err := o.Insert(&Role{Name: name}); err != nil {
		// ignore duplicate
		return nil
	}
	return nil
}

func AssignRole(email, role string) error {
	o := orm.NewOrm()
	_, err := o.Insert(&UserRole{Email: email, Role: role})
	return err
}

func Grant(role, resource, action string) error {
	o := orm.NewOrm()
	_, err := o.Insert(&Permission{Role: role, Resource: resource, Action: action})
	return err
}

func HasRole(email, role string) (bool, error) {
	o := orm.NewOrm()
	cnt, err := o.QueryTable(new(UserRole)).Filter("Email", email).Filter("Role", role).Count()
	return cnt > 0, err
}

func HasPermission(email, resource, action string) (bool, error) {
	o := orm.NewOrm()
	// superuser role bypass happens at controller level via user flag/role
	var roles []UserRole
	_, err := o.QueryTable(new(UserRole)).Filter("Email", email).All(&roles)
	if err != nil {
		return false, err
	}
	if len(roles) == 0 {
		return false, nil
	}
	for _, ur := range roles {
		cnt, err := o.QueryTable(new(Permission)).Filter("Role", ur.Role).Filter("Resource", resource).Filter("Action", action).Count()
		if err != nil {
			return false, err
		}
		if cnt > 0 {
			return true, nil
		}
	}
	return false, nil
}

func RequirePermission(email, resource, action string) error {
	ok, err := HasPermission(email, resource, action)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("forbidden")
	}
	return nil
}

func MigrateUserEmail(oldEmail, newEmail string) error {
	o := orm.NewOrm()
	// update user_roles
	_, err := o.QueryTable(new(UserRole)).Filter("Email", oldEmail).Update(orm.Params{"Email": newEmail})
	if err != nil {
		return err
	}
	return nil
}
