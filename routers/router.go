package routers

import (
	"github.com/beego/beego/v2/server/web"

	items "github.com/mymi14s/goconda/apps/items/controllers"
	"github.com/mymi14s/goconda/controllers"
	"github.com/mymi14s/goconda/middleware"
)

func init() {

	middleware.ProtectMany(
		"/api/v1/users/me",
		"/api/v1/items",
		"/api/v1/items/*", // covers /items/:id paths
		"/api/v1/upload",
	)

	ns := web.NewNamespace("/api/v1",
		web.NSNamespace("/auth",
			web.NSRouter("/register", &controllers.AuthController{}, "post:Register"),
			web.NSRouter("/login", &controllers.AuthController{}, "post:Login"),
			web.NSRouter("/logout", &controllers.LogoutController{}, "post:Logout"),
			web.NSRouter("/forgot-password", &controllers.AuthController{}, "post:ForgotPassword"),
			web.NSRouter("/reset-password", &controllers.AuthController{}, "post:ResetPassword"),
			web.NSRouter("/change-password", &controllers.AuthController{}, "post:ChangePassword"),
			web.NSRouter("/change-email", &controllers.AuthController{}, "post:ChangeEmail"),
			web.NSRouter("/send-verification", &controllers.AuthController{}, "post:SendVerification"),
			web.NSRouter("/verify", &controllers.AuthController{}, "get:VerifyEmail"),
		),
		web.NSRouter("/users/me", &controllers.UserController{}, "get:Me"),
		web.NSRouter("/items", &items.ItemController{}, "get:List;post:Create"),
		web.NSRouter("/items/:id", &items.ItemController{}, "get:GetOne;put:Update;delete:Delete"),
		web.NSRouter("/upload", &controllers.UploadController{}, "post:Upload"),
	)
	web.AddNamespace(ns)

}
