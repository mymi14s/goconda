// middleware/auth.go
package middleware

import (
	"strings"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"

	"github.com/mymi14s/goconda/controllers"
)

func Protect(path string) {
	web.InsertFilter(path, web.BeforeRouter, func(ctx *context.Context) {
		if !controllers.RequireAuth(ctx) { // return bool to indicate pass/fail
			return // IMPORTANT: stop here; Output.JSON marks response as started
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
