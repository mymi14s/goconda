package controllers

import (
	"errors"
	"strings"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"

	"github.com/mymi14s/goconda/models"
	"github.com/mymi14s/goconda/utils"
	jwtutil "github.com/mymi14s/goconda/utils/jwt"
	"github.com/mymi14s/goconda/utils/response"
)

// Context keys
const ctxUserKey = "current_user"

type BaseController struct {
	web.Controller
}

// JSON helpers
func (c *BaseController) JSONOK(data interface{}) {
	response.JSONOK(c.Ctx, data)
}

func (c *BaseController) JSONError(code int, msg string) {
	response.JSONError(c.Ctx, code, msg)
}

// GetCurrentUser returns the authenticated user set by middleware.
// If not set, it attempts to resolve from the Authorization header.
func (c *BaseController) GetCurrentUser() (*models.User, error) {
	if c.Ctx == nil || c.Ctx.Input == nil {
		return nil, errors.New("no context")
	}
	if v := c.Ctx.Input.GetData(ctxUserKey); v != nil {
		if u, ok := v.(*models.User); ok {
			return u, nil
		}
	}

	// Fallback: try to parse from Authorization header
	auth := c.Ctx.Request.Header.Get("Authorization")
	if auth == "" {
		return nil, errors.New("no auth header")
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return nil, errors.New("invalid authorization header")
	}
	claims, err := jwtutil.Parse(parts[1])
	if err != nil {
		return nil, err
	}
	u, err := models.GetUserByEmail(claims.Email)
	if err != nil || u == nil {
		return nil, errors.New("user not found")
	}
	// cache in context for the remainder of the request
	c.Ctx.Input.SetData(ctxUserKey, u)
	return u, nil
}

// MustAuth aborts the request with 401 if user is not authenticated
func (c *BaseController) MustAuth() (*models.User, bool) {
	u, err := c.GetCurrentUser()
	if err != nil || u == nil {
		c.JSONError(401, "unauthorized")
		return nil, false
	}
	return u, true
}

func (c *BaseController) ParseJSON(v interface{}) error {
	return utils.ParseJSON(c.Ctx.Request, v)
}

// RequireAuth is used by middleware to protect arbitrary routes not bound to a controller instance.
func RequireAuth(ctx *context.Context) bool {
	if ctx == nil || ctx.Input == nil {
		return false
	}
	auth := ctx.Request.Header.Get("Authorization")
	if auth == "" {
		response.JSONError(ctx, 401, "missing Authorization header")
		return false
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
		response.JSONError(ctx, 401, "invalid Authorization format")
		return false
	}

	claims, err := jwtutil.Parse(parts[1])
	if err != nil {
		response.JSONError(ctx, 401, "invalid or expired token")
		return false
	}

	// Check revocation, if supported
	if revoked, err := models.IsTokenRevoked(claims.ID); err == nil && revoked {
		response.JSONError(ctx, 401, "token revoked")
		return false
	}

	u, err := models.GetUserByEmail(claims.Email)
	if err != nil || u == nil {
		response.JSONError(ctx, 401, "account not found")
		return false
	}

	// Store the user on the context for controllers to read later
	ctx.Input.SetData(ctxUserKey, u)
	return true
}

// superuser bypass
func (c *BaseController) IsEmailVerified(email string) bool {
    u, _ := c.GetCurrentUser()
    if u != nil && (u.IsSuperuser) {
        return true
    }
    ok, _ := models.IsUserVerified(email)
    return ok
}

func (c *BaseController) RequirePermission(resource, action string) bool {
	u, err := c.GetCurrentUser()
	if err != nil || u == nil {
		response.JSONError(c.Ctx, 401, "unauthorized")
		return false
	}
	// bypass if superuser or has Superuser role
	if u.IsSuperuser {
		return true
	}
	if ok, _ := models.HasRole(u.Email, "Superuser"); ok {
		return true
	}
	if err := models.RequirePermission(u.Email, resource, action); err != nil {
		response.JSONError(c.Ctx, 403, "forbidden")
		return false
	}
	return true
}
