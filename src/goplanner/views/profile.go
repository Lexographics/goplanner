package goplanner

import (
	"net/http"

	"github.com/labstack/echo/v4"

	db "github.com/Lexographics/goplanner/src/goplanner/db"
)


func ProfileView(c echo.Context) error {
	cookie, err := c.Cookie("sessionId")
	if err != nil {
		return RedirectToAuthView(c)
	}

	id, success, err := db.ValidateSession(cookie.Value)
	if err != nil {
		return RedirectToAuthView(c)
	}

	if success {
		var user db.User
		res := db.Database.First(&user, "id = ?", id)
		if res.Error != nil {
			return RedirectToHomeView(c)
		}


		type ProfilePage struct {
			Username string
			Email    string
		}
		page := ProfilePage{
			Username: user.Username,
			Email: user.Email,
		}
		return c.Render(http.StatusOK, "ProfilePage", page)
	}
	
	return RedirectToAuthView(c)
}