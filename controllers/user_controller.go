package controllers

import "fmt"

type UserController struct {
	BaseController
}

func (c *UserController) Me() {

	fmt.Println("hello world\npero")
	u, err := c.GetCurrentUser()
	fmt.Println(u)
	if err != nil || u == nil {
		c.JSONError(401, "unauthorized")
		return
	}
	c.JSONOK(map[string]interface{}{
		"email":      u.Email,
		"first_name": u.FirstName,
		"last_name":  u.LastName,
		"created_at": u.CreatedAt,
		"updated_at": u.UpdatedAt,
	})
}
