package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/mymi14s/goconda/apps/items/models"
	base_controller "github.com/mymi14s/goconda/controllers"

	"github.com/beego/beego/v2/client/orm"
)

type ItemController struct {
	base_controller.BaseController
}

func (c *ItemController) Prepare() {
	c.MustAuth()
}

type itemReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// @router /api/v1/items [get]
func (c *ItemController) List() {
	user, ok := c.MustAuth()
	if !ok {
		return
	}
	limit, _ := c.GetInt64("limit", 20)
	offset, _ := c.GetInt64("offset", 0)
	items, total, err := models.ListItemsByOwner(user.Email, limit, offset)
	if err != nil {
		c.JSONError(500, "failed to list items")
		return
	}
	c.JSONOK(map[string]interface{}{
		"total": total,
		"items": items,
	})
}

// @router /api/v1/items [post]
func (c *ItemController) Create() {
	user, ok := c.MustAuth()
	if !ok {
		return
	}

	var req itemReq
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.JSONError(400, "invalid json")
		return
	}
	if req.Name == "" {
		c.JSONError(400, "name is required")
		return
	}
	it := models.Item{
		Name:        req.Name,
		Description: req.Description,
		Owner:       user,
	}
	o := orm.NewOrm()
	if _, err := o.Insert(&it); err != nil {
		c.JSONError(500, "failed to create item")
		return
	}
	c.JSONOK(it)
}

// @router /api/v1/items/:id [get]
func (c *ItemController) GetOne() {
	user, ok := c.MustAuth()
	if !ok {
		return
	}
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	o := orm.NewOrm()
	it := models.Item{ID: id}
	if err := o.Read(&it); err != nil {
		c.JSONError(404, "not found")
		return
	}
	o.LoadRelated(&it, "Owner")
	if it.Owner == nil || it.Owner.Email != user.Email {
		c.JSONError(403, "forbidden")
		return
	}
	c.JSONOK(it)
}

// @router /api/v1/items/:id [put]
func (c *ItemController) Update() {
	user, ok := c.MustAuth()
	if !ok {
		return
	}
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	o := orm.NewOrm()
	it := models.Item{ID: id}
	if err := o.Read(&it); err != nil {
		c.JSONError(404, "not found")
		return
	}
	o.LoadRelated(&it, "Owner")
	if it.Owner == nil || it.Owner.Email != user.Email {
		c.JSONError(403, "forbidden")
		return
	}
	var req itemReq
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.JSONError(400, "invalid json")
		return
	}
	if req.Name != "" {
		it.Name = req.Name
	}
	if req.Description != "" {
		it.Description = req.Description
	}
	if _, err := o.Update(&it); err != nil {
		c.JSONError(500, "failed to update item")
		return
	}
	c.JSONOK(it)
}

// @router /api/v1/items/:id [delete]
func (c *ItemController) Delete() {
	user, ok := c.MustAuth()
	if !ok {
		return
	}
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	o := orm.NewOrm()
	it := models.Item{ID: id}
	if err := o.Read(&it); err != nil {
		c.JSONError(404, "not found")
		return
	}
	o.LoadRelated(&it, "Owner")
	if it.Owner == nil || it.Owner.Email != user.Email {
		c.JSONError(403, "forbidden")
		return
	}
	if _, err := o.Delete(&it); err != nil {
		c.JSONError(500, "failed to delete")
		return
	}
	c.JSONOK(map[string]any{"deleted": id})
}
