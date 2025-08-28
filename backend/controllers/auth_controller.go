package controllers

import (
	"net/http"
	"github.com/beego/beego/v2/server/web"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"

	"github.com/mymi14s/goconda/models"
	"github.com/mymi14s/goconda/utils/hash"
	jwtutil "github.com/mymi14s/goconda/utils/jwt"
	"github.com/mymi14s/goconda/utils/validators"
)

type AuthController struct {
	BaseController
}

type registerPayload struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (c *AuthController) Register() {
	var p registerPayload
	if err := c.ParseJSON(&p); err != nil {
		c.JSONError(400, err.Error())
		return
	}
	if err := validators.ValidateEmail(p.Email); err != nil {
		c.JSONError(400, "invalid email")
		return
	}
	if err := validators.RequireNonEmpty(map[string]string{
		"password":   p.Password,
		"first_name": p.FirstName,
		"last_name":  p.LastName,
	}); err != nil {
		c.JSONError(400, err.Error())
		return
	}

	// ensure email normalized
	email := strings.ToLower(strings.TrimSpace(p.Email))

	// Check exists
	if existing, _ := models.GetUserByEmail(email); existing != nil {
		c.JSONError(409, "user already exists")
		return
	}

	pwHash, err := hash.HashPassword(p.Password)
	if err != nil {
		c.JSONError(500, "failed to hash password")
		return
	}

	u := &models.User{
		Email:        email,
		FirstName:    p.FirstName,
		LastName:     p.LastName,
		PasswordHash: pwHash,
	}
	if _, err := orm.NewOrm().Insert(u); err != nil {
		c.JSONError(500, "failed to create user")
		return
	}

	// Issue token
	token, err := jwtutil.Generate(email)
	if err != nil {
		c.JSONError(500, "failed to generate token")
		return
	}

	

// --- BFF: set HttpOnly auth cookie instead of relying on sessionStorage ---
{
    // Use JWT as cookie value; browser won't see it (HttpOnly)
    // Cookie attrs depend on runmode for cross-site dev vs same-origin prod
    runMode := web.BConfig.RunMode
    secure := true // default
    sameSite := http.SameSiteLaxMode
    if runMode == "dev" {
        // Dev runs cross-site (frontend :3000 over http, backend :8080 over https)
        // Cookies for cross-site must be SameSite=None and Secure
        sameSite = http.SameSiteNoneMode
        secure = true
    }
    // Expiration: align with JWT exp, fallback to 60m
    expMin, _ := web.AppConfig.Int64("jwt::expiration_minutes")
    if expMin <= 0 { expMin = 60 }
    expires := time.Now().Add(time.Duration(expMin) * time.Minute)
    cookie := &http.Cookie{
        Name:     "goconda_auth",
        Value:    token,
        Path:     "/",
        HttpOnly: true,
        Secure:   secure,
        SameSite: sameSite,
        Expires:  expires,
        MaxAge:   int(time.Until(expires).Seconds()),
    }
    http.SetCookie(c.Ctx.ResponseWriter, cookie)
}
// --- end BFF cookie ---

c.JSONOK(map[string]any{
		"token":      token,
		"token_type": "Bearer",
		"expires_in": int(time.Hour.Seconds()) * 24, // ideally from jwtutil
		"user": map[string]any{
			"email":      u.Email,
			"first_name": u.FirstName,
			"last_name":  u.LastName,
		},
	})
}

type loginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *AuthController) Login() {
	var p loginPayload
	if err := c.ParseJSON(&p); err != nil {
		c.JSONError(400, err.Error())
		return
	}
	email := strings.ToLower(strings.TrimSpace(p.Email))
	if err := validators.ValidateEmail(email); err != nil {
		c.JSONError(400, "invalid email")
		return
	}
	if strings.TrimSpace(p.Password) == "" {
		c.JSONError(400, "password is required")
		return
	}

	u, err := models.GetUserByEmail(email)
	if err != nil || u == nil {
		c.JSONError(401, "invalid credentials")
		return
	}
	if !hash.CheckPassword(p.Password, u.PasswordHash) {
		c.JSONError(401, "invalid credentials")
		return
	}

	token, err := jwtutil.Generate(email)
	if err != nil {
		c.JSONError(500, "failed to generate token")
		return
	}

	c.JSONOK(map[string]any{
		"token":      token,
		"token_type": "Bearer",
		"email":      u.Email,
		"first_name": u.FirstName,
		"lastname":   u.LastName,
	})
}

type logoutPayload struct {
	JTI string `json:"jti"`
}

type LogoutController struct {
	BaseController
}

// A basic logout that revokes by JTI if provided in claims
func (c *LogoutController) Logout() {
	u, ok := c.MustAuth()
	if !ok || u == nil {
		return
	}
	// The token's JTI is provided by jwtutil via Parse() but we don't get it here.
	// In a real app, you'd extract the bearer and revoke by claims.ID.
	// We'll best-effort parse the header.
	auth := c.Ctx.Request.Header.Get("Authorization")
	if auth == "" {
		

// Clear auth cookie on logout
http.SetCookie(c.Ctx.ResponseWriter, &http.Cookie{
    Name:     "goconda_auth",
    Value:    "",
    Path:     "/",
    HttpOnly: true,
    Secure:   true,
    SameSite: http.SameSiteLaxMode,
    MaxAge:   -1,
    Expires:  time.Unix(0,0),
})

c.JSONOK(map[string]any{"revoked": false})
		return
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 {
		c.JSONOK(map[string]any{"revoked": false})
		return
	}
	claims, err := jwtutil.Parse(parts[1])
	if err != nil || claims == nil || claims.ID == "" {
		c.JSONOK(map[string]any{"revoked": false})
		return
	}
	// Store revocation
	_ = models.RevokeToken(claims.ID, claims.ExpiresAt.Time)
	c.JSONOK(map[string]any{"revoked": true})
}

func (c *AuthController) SendVerification() {
	email := strings.TrimSpace(c.GetString("email"))
	if email == "" {
		c.JSONError(400, "email is required")
		return
	}
	// Ensure user exists
	u, _ := models.GetUserByEmail(email)
	if u == nil {
		c.JSONError(404, "account not found")
		return
	}
	// Create token (valid 24h)
	t, err := models.CreateVerificationToken(email, 24*time.Hour)
	if err != nil {
		c.JSONError(500, "could not create token")
		return
	}
	// "Send" email (for now, log to server; in production wire SMTP)
	// Example message a real mailer would send:
	// Verify link: {BASE_URL}/api/v1/auth/verify?token=t.Token
	c.Ctx.WriteString("Verification token: " + t.Token)
}

func (c *AuthController) VerifyEmail() {
	token := strings.TrimSpace(c.GetString("token"))
	if token == "" {
		c.JSONError(400, "token is required")
		return
	}
	email, err := models.ConsumeVerificationToken(token)
	if err != nil {
		c.JSONError(400, err.Error())
		return
	}
	if err := models.MarkUserVerified(email); err != nil {
		c.JSONError(500, "could not mark verified")
		return
	}
	c.JSONOK(map[string]any{"email": email, "verified": true})
}

func (c *AuthController) ForgotPassword() {
	email := strings.TrimSpace(c.GetString("email"))
	if email == "" {
		c.JSONError(400, "email is required")
		return
	}
	u, _ := models.GetUserByEmail(email)
	if u == nil {
		// don't reveal account existence
		c.JSONOK(map[string]any{"ok": true})
		return
	}
	t, err := models.CreatePasswordResetToken(email, 1*time.Hour)
	if err != nil {
		c.JSONError(500, "could not create token")
		return
	}
	// In production, send via email. For now, write token to body.
	c.Ctx.WriteString("Password reset token: " + t.Token)
}

func (c *AuthController) ResetPassword() {
	token := strings.TrimSpace(c.GetString("token"))
	newPass := strings.TrimSpace(c.GetString("new_password"))
	if token == "" || len(newPass) < 6 {
		c.JSONError(400, "invalid token or password too short")
		return
	}
	email, err := models.ConsumePasswordResetToken(token)
	if err != nil {
		c.JSONError(400, err.Error())
		return
	}
	o := orm.NewOrm()
	u, _ := models.GetUserByEmail(email)
	if u == nil {
		c.JSONError(404, "account not found")
		return
	}
	hv, _ := hash.Make(newPass)
	u.PasswordHash = hv
	if _, err := o.Update(u, "PasswordHash"); err != nil {
		c.JSONError(500, "failed to update password")
		return
	}
	c.JSONOK(map[string]any{"reset": true})
}

func (c *AuthController) ChangePassword() {
	u, err := c.GetCurrentUser()
	if err != nil || u == nil {
		c.JSONError(401, "unauthorized")
		return
	}
	current := strings.TrimSpace(c.GetString("current_password"))
	newPass := strings.TrimSpace(c.GetString("new_password"))
	if len(newPass) < 6 {
		c.JSONError(400, "password too short")
		return
	}
	if !hash.Check(current, u.PasswordHash) {
		c.JSONError(400, "current password incorrect")
		return
	}
	hv, _ := hash.Make(newPass)
	u.PasswordHash = hv
	if _, err := orm.NewOrm().Update(u, "PasswordHash"); err != nil {
		c.JSONError(500, "failed to change password")
		return
	}
	c.JSONOK(map[string]any{"changed": true})
}

func (c *AuthController) ChangeEmail() {
	u, err := c.GetCurrentUser()
	if err != nil || u == nil {
		c.JSONError(401, "unauthorized")
		return
	}
	pwd := strings.TrimSpace(c.GetString("password"))
	newEmail := strings.TrimSpace(c.GetString("new_email"))
	if newEmail == "" || !validators.IsEmailValid(newEmail) {
		c.JSONError(400, "invalid email")
		return
	}
	if !hash.Check(pwd, u.PasswordHash) && !u.IsSuperuser {
		c.JSONError(400, "password incorrect")
		return
	}
	// prevent collision
	if existing, _ := models.GetUserByEmail(newEmail); existing != nil {
		c.JSONError(400, "email already in use")
		return
	}
	old := u.Email
	u.Email = newEmail
	if _, err := orm.NewOrm().Update(u, "Email"); err != nil {
		c.JSONError(500, "failed to change email")
		return
	}
	// migrate roles to new email
	_ = models.MigrateUserEmail(old, newEmail)
	c.JSONOK(map[string]any{"email": newEmail})
}
