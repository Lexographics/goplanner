package goplanner

import (
	_ "database/sql"
	_ "time"

	_ "github.com/go-sql-driver/mysql"

	_ "github.com/golang-jwt/jwt/v4"

	"github.com/labstack/echo/v4"
)



func AuthView(c echo.Context) error {
	return c.File("./static/auth.html")
}