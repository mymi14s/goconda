package models

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

type BaseModel struct{}

// Generic: fill any struct from map using `json` tags (or field names, case-insensitive),
// then insert with Beego ORM. Works for any model type.
func (BaseModel) Create(dst any, data map[string]any) (int64, error) {
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:          "json", // uses your `json:"..."` tags
		WeaklyTypedInput: true,   // "123" -> 123 etc.
		Result:           dst,
		MatchName: func(mapKey, fieldName string) bool {
			// allow case-insensitive match on field names when no tag present
			return strings.EqualFold(mapKey, fieldName)
		},
	})
	if err != nil {
		return 0, err
	}
	if err := dec.Decode(data); err != nil {
		return 0, err
	}

	// Optional: if model has a string `ID` field and it's empty, populate a UUID.
	autoUUIDIfStringID(dst)

	o := orm.NewOrm()
	return o.Insert(dst)
}

func autoUUIDIfStringID(dst any) {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return
	}
	e := v.Elem()
	f := e.FieldByName("ID")
	if f.IsValid() && f.CanSet() && f.Kind() == reflect.String && f.Len() == 0 {
		f.SetString(uuid.NewString())
	}
}

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
		new(ErrorLog),
	)
	return nil
}
