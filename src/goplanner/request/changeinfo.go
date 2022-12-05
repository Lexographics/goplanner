package goplanner

import (
	"net/http"

	"github.com/labstack/echo/v4"

	db "github.com/Lexographics/goplanner/src/goplanner/db"
)

// GET
func ChangeInfoRequest(c echo.Context) error {
	newUsername := c.FormValue("username")
	info := c.FormValue("info")

	type State struct {
		Status string
	}

	if info != "username" && info != "email" && info != "password" {
		return c.JSON(http.StatusUnauthorized, State{
			Status: "Error",
		})
	}

	cookie, err := c.Cookie("sessionId")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, State{
			Status: "Error",
		})
	}


	id, success, err := db.ValidateSession(cookie.Value)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, State{
			Status: "Unauthorized",
		})
	}

	if success {
		user := db.User{
			ID: uint(id),
		}
		res := db.Database.Model(&user).Update(info, newUsername)
		if res.Error != nil {
			return c.JSON(http.StatusUnauthorized, State{
				Status: "Unauthorized",
			})
		}

		return c.JSON(http.StatusOK, State{
			Status: "Success",
		})
	}

	return c.JSON(http.StatusUnauthorized, State{
		Status: "Error",
	})
}
