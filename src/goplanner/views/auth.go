package goplanner

import (
	"github.com/labstack/echo/v4"
)



func AuthView(c echo.Context) error {
	return c.File("./static/auth.html")
}