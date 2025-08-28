package controllers

import (
	"fmt"

	base_controller "github.com/mymi14s/goconda/controllers"
	"github.com/mymi14s/goconda/models"
	"github.com/mymi14s/goconda/utils"
	"github.com/mymi14s/goconda/utils/mailer"
)

type FrontendController struct {
	base_controller.BaseController
}

// @router / [get]
func (c *FrontendController) Index() {

	// Point to your template file (adjust the path/ext to match your views)
	c.TplName = "frontend/index.html" // e.g. views/frontend/index.html

	var ERROR_TITLE string = "Email Sender"
	var error string = "smtp not configured (host/user/pass/from)"
	utils.LogError(map[string]any{
		"title":  ERROR_TITLE,
		"source": "Mailer",
		"error":  error,
	})

	// Render now (explicit) — or rely on AutoRender if enabled
	if err := c.Render(); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.WriteString("template render error")
	}
}

// @router /frontend/api/data [get]
func (c *FrontendController) GetInfo() {

	ss := models.SiteSetting{}
	data, _ := ss.Get()
	c.JSONOK(&data)
}

type ContactFormData struct {
	Email   string
	Subject string
	Message string
}

func (c *FrontendController) ContactForm() {
	var form ContactFormData
	if err := c.ParseJSON(&form); err != nil {
		fmt.Println((err))
		c.JSONError(400, err.Error())
		return
	}

	// fetch email from settings
	ss := models.SiteSetting{}
	data, _ := ss.Get()

	if err := mailer.SendEmail(form.Message, []string{data.Email}, form.Subject); err != nil {
		c.JSONError(400, err.Error())
	}
	c.JSONOK(map[string]any{"status": true})
}
