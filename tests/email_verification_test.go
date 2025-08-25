package tests

import (
    "testing"
    "time"

    "github.com/beego/beego/v2/client/orm"
    "github.com/beego/beego/v2/server/web"

    "github.com/mymi14s/goconda/models"
)

func init() {
    web.AppConfigPath = "conf/app.dev.conf"
    _ = models.InitDB()
    orm.RunSyncdb("default", true, true)
}

func TestEmailVerificationFlow(t *testing.T) {
    // create user
    o := orm.NewOrm()
    u := &models.User{Email: "alice@example.com", FirstName: "Alice", LastName: "Liddell", PasswordHash: "x"}
    if _, err := o.Insert(u); err != nil {
        t.Fatalf("insert user: %v", err)
    }

    // initially not verified
    ok, err := models.IsUserVerified(u.Email)
    if err != nil {
        t.Fatalf("IsUserVerified error: %v", err)
    }
    if ok {
        t.Fatalf("expected not verified")
    }

    // create token
    tok, err := models.CreateVerificationToken(u.Email, time.Hour)
    if err != nil {
        t.Fatalf("CreateVerificationToken: %v", err)
    }
    if tok.Email != u.Email || tok.Token == "" {
        t.Fatalf("invalid token %+v", tok)
    }

    // consume
    email, err := models.ConsumeVerificationToken(tok.Token)
    if err != nil {
        t.Fatalf("ConsumeVerificationToken: %v", err)
    }
    if email != u.Email {
        t.Fatalf("expected email match")
    }

    // mark verified
    if err := models.MarkUserVerified(email); err != nil {
        t.Fatalf("MarkUserVerified: %v", err)
    }

    // now verified
    ok, err = models.IsUserVerified(u.Email)
    if err != nil {
        t.Fatalf("IsUserVerified 2: %v", err)
    }
    if !ok {
        t.Fatalf("expected verified true")
    }
}
