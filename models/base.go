package models

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
)

func InitDB() error {
	driver := web.AppConfig.DefaultString("db::driver", "sqlite3")
	dsn := web.AppConfig.DefaultString("db::dsn", "file:goconda_dev.db?cache=shared&_fk=1")

	switch driver {
	case "mysql":
		orm.RegisterDriver("mysql", orm.DRMySQL)
	case "sqlite3":
		orm.RegisterDriver("sqlite3", orm.DRSqlite)
	default:
		orm.RegisterDriver(driver, orm.DRSqlite)
	}

	if err := orm.RegisterDataBase("default", driver, dsn); err != nil {
		return fmt.Errorf("register db: %w", err)
	}

	orm.RegisterModel(
		new(User),
		new(RevokedToken),
		new(EmailVerificationToken),
		new(VerifiedUser),
		new(Role),
		new(UserRole),
		new(Permission),
		new(PasswordResetToken),
	)
	return nil
}
