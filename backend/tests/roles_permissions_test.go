package tests

import (
    "testing"

    "github.com/beego/beego/v2/client/orm"
    "github.com/beego/beego/v2/server/web"
    "github.com/mymi14s/goconda/models"
)

func init() {
    web.AppConfigPath = "conf/app.dev.conf"
    _ = models.InitDB()
    orm.RunSyncdb("default", true, true)
}

func TestRolesPermissions(t *testing.T) {
    o := orm.NewOrm()
    u := &models.User{Email: "bob@example.com", FirstName: "Bob", LastName: "B", PasswordHash: "x"}
    if _, err := o.Insert(u); err != nil {
        t.Fatalf("insert user: %v", err)
    }
    if err := models.EnsureRole("Reader"); err != nil { t.Fatal(err) }
    if err := models.AssignRole(u.Email, "Reader"); err != nil { t.Fatal(err) }
    if err := models.Grant("Reader", "items", "read"); err != nil { t.Fatal(err) }

    ok, err := models.HasPermission(u.Email, "items", "read")
    if err != nil { t.Fatal(err) }
    if !ok { t.Fatalf("expected permission") }

    ok, err = models.HasPermission(u.Email, "items", "delete")
    if err != nil { t.Fatal(err) }
    if ok { t.Fatalf("should not have delete") }
}
