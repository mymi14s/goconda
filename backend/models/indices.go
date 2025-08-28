package models

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
)

// EnsureUserRoleUniqueIndex creates a composite unique index on (email, role)
// to prevent duplicate role names for the same user.
func EnsureUserRoleUniqueIndex() error {
	driver := web.AppConfig.DefaultString("db::driver", "sqlite3")
	o := orm.NewOrm()
	var stmt string
	switch driver {
	case "mysql":
		stmt = "ALTER TABLE user_roles ADD UNIQUE KEY idx_user_role_unique (email, role)"
	default:
		// sqlite and others
		stmt = "CREATE UNIQUE INDEX IF NOT EXISTS idx_user_role_unique ON user_roles (email, role)"
	}
	// Attempt execution; ignore specific errors when index already exists.
	if _, err := o.Raw(stmt).Exec(); err != nil {
		// For MySQL, error 1061 means duplicate key name; 1062 duplicate entry; we just log not fatal
		return fmt.Errorf("create unique index failed: %w", err)
	}
	return nil
}
