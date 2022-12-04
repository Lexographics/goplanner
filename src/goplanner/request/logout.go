package goplanner

import (
	_ "database/sql"
	"fmt"
	_ "time"

	_ "github.com/go-sql-driver/mysql"

	_ "github.com/golang-jwt/jwt/v4"

	"github.com/labstack/echo/v4"

	db "github.com/Lexographics/goplanner/src/goplanner/db"
	views "github.com/Lexographics/goplanner/src/goplanner/views"
)

// GET
func LogoutRequest(c echo.Context) error {
	session, err := c.Cookie("sessionId")

	if err != nil {
		fmt.Printf("LogoutRequest error 1: %s", err)
		return views.RedirectToAuthView(c)
	}
	db.InvalidateSession(session.Value)

	return views.RedirectToAuthView(c)
}
