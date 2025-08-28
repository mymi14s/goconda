package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"

	"github.com/mymi14s/goconda/models"
	_ "github.com/mymi14s/goconda/routers"
	"github.com/mymi14s/goconda/utils/hash"
)

func mustLoadConfig() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	confPath := filepath.Join("conf", fmt.Sprintf("app.%s.conf", env))
	if _, err := os.Stat(confPath); err != nil {
		log.Fatalf("config file not found: %s (APP_ENV=%s)", confPath, env)
	}
	if err := web.LoadAppConfig("ini", confPath); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
}

func setupSessionsAndStatic() {
	web.BConfig.WebConfig.Session.SessionOn = true
	web.BConfig.WebConfig.Session.SessionProvider = "file"
	web.BConfig.WebConfig.Session.SessionProviderConfig = "./.sessions"
	if err := os.MkdirAll("./.sessions", 0o755); err != nil {
		log.Printf("warn: could not create session dir: %v", err)
	}

	web.SetStaticPath("/static", "static")
}

func bootstrapAdmin() error {
	// Read admin from config
	adminEmail := web.AppConfig.DefaultString("admin::email", "")
	adminPass := web.AppConfig.DefaultString("admin::password", "")
	if adminEmail == "" || adminPass == "" {
		return nil
	}
	// Ensure Superuser role exists
	_ = models.EnsureRole("Superuser")

	// Create user if not exists
	o := orm.NewOrm()
	existing, _ := models.GetUserByEmail(adminEmail)
	if existing == nil {
		_hash, _ := hash.Make(adminPass)
		u := &models.User{Email: adminEmail, FirstName: "Admin", LastName: "User", PasswordHash: _hash, IsSuperuser: true}
		if _, err := o.Insert(u); err != nil {
			return err
		}
	} else {
		if !existing.IsSuperuser {
			existing.IsSuperuser = true
			_, _ = o.Update(existing, "IsSuperuser")
		}
	}

	// assign role
	_ = models.AssignRole(adminEmail, "Superuser")

	return nil
}

func main() {
	// Enable sessions
	web.BConfig.WebConfig.Session.SessionOn = true
	web.BConfig.WebConfig.Session.SessionName = "bffsid"
	web.BConfig.WebConfig.Session.SessionCookieLifeTime = 86400 // 1 day
	mustLoadConfig()
	setupSessionsAndStatic()

	if err := models.InitDB(); err != nil {
		log.Fatalf("DB init failed: %v", err)
	}
	if err := orm.RunSyncdb("default", false, true); err != nil {
		log.Fatalf("RunSyncdb error: %v", err)
	}
	if err := bootstrapAdmin(); err != nil {
		log.Printf("bootstrap admin: %v", err)
	}

	port, _ := web.AppConfig.Int("httpport")
	appname := web.AppConfig.DefaultString("appname", "goconda")
	log.Printf("%s starting on :%d", appname, port)
	web.InsertFilter("*", web.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true, // if you need cookies/JWT via cookie
		MaxAge:           600,
	}))

	web.BConfig.Listen.HTTPPort = port
	web.SetStaticPath("/static", "static")

	web.Run()
}
