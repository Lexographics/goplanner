package goplanner

import (
	_ "database/sql"
	"fmt"
	"net/http"
	_ "time"

	_ "github.com/go-sql-driver/mysql"

	_ "github.com/golang-jwt/jwt/v4"

	"github.com/labstack/echo/v4"

	db "github.com/Lexographics/goplanner/src/goplanner/db"
)


func ProfileView(c echo.Context) error {
	cookie, err := c.Cookie("sessionId")
	if err != nil {
		fmt.Printf("ProfileView error 1: %s\n", err)
		return RedirectToAuthView(c)
	}

	id, success, err := db.ValidateSession(cookie.Value)
	if err != nil {
		fmt.Printf("ProfileView error 2: %s\n", err)
		return RedirectToAuthView(c)
	}

	if success {
		var user db.User
		res := db.Database.First(&user, "id = ?", id)
		if res.Error != nil {
			fmt.Printf("Error 1 LoginRequest: %s\n", res.Error)
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
	

	fmt.Printf("ProfileView error 5: %s\n", err)
	return RedirectToAuthView(c)
}