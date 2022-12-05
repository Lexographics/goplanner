package goplanner

import (
	"github.com/labstack/echo/v4"

	db "github.com/Lexographics/goplanner/src/goplanner/db"
	views "github.com/Lexographics/goplanner/src/goplanner/views"
)

// POST
func LoginRequest(c echo.Context) error {
	name := c.FormValue("name")
	password := c.FormValue("password")

	var user db.User
	res := db.Database.First(&user, "username = ? AND password = ?", name, password)
	if res.Error != nil {
		return views.RedirectToHomeView(c)
	}
	
	cookie, err := db.CreateToken(int64(user.ID))
	if err != nil {
		return views.RedirectToHomeView(c)
	}
	
	c.SetCookie(&cookie)
	return views.RedirectToHomeView(c)
}
