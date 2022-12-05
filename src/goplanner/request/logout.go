package goplanner

import (
	"github.com/labstack/echo/v4"

	db "github.com/Lexographics/goplanner/src/goplanner/db"
	views "github.com/Lexographics/goplanner/src/goplanner/views"
)

// GET
func LogoutRequest(c echo.Context) error {
	session, err := c.Cookie("sessionId")

	if err != nil {
		return views.RedirectToAuthView(c)
	}
	db.InvalidateSession(session.Value)

	return views.RedirectToAuthView(c)
}
