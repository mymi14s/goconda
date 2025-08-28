// middleware/auth.go
package middleware

import (
	"strings"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"

	"github.com/mymi14s/goconda/controllers"
)

// SetupCookieAuthBridge maps the HttpOnly auth cookie to Authorization header
// so existing JWT-based middleware keeps working without exposing tokens to JS.
func SetupCookieAuthBridge() {
	web.InsertFilter("/api/v1/*", web.BeforeRouter, func(ctx *context.Context) {
		// If Authorization already present, respect it
		if ctx.Request.Header.Get("Authorization") == "" {
			if c, err := ctx.Request.Cookie("goconda_auth"); err == nil && c != nil && c.Value != "" {
				ctx.Request.Header.Set("Authorization", "Bearer "+c.Value)
			}
		}
		// In dev, also allow preflight with credentials by echoing Origin
		if ctx.Input.IsOptions() {
			origin := ctx.Request.Header.Get("Origin")
			if origin != "" {
				ctx.Output.Header("Access-Control-Allow-Origin", origin)
				ctx.Output.Header("Vary", "Origin")
			}
		}
	})
}

func Protect(path string) {
	web.InsertFilter(path, web.BeforeRouter, func(ctx *context.Context) {
		if !controllers.RequireAuth(ctx) { // return bool to indicate pass/fail
			return // IMPORTANT: stop here; Output.JSON marks response as started
		}

		if ctx.Input.IsOptions() {
			ctx.Output.SetStatus(200)
			return
		}

	})
}

func ProtectMany(paths ...string) {
	for _, p := range paths {
		Protect(p)

		if !strings.HasSuffix(p, "/*") && !strings.HasSuffix(p, "/") {
			Protect(p + "/")
		}
	}
}
