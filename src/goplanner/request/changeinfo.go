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

// GET
func ChangeInfoRequest(c echo.Context) error {
	newUsername := c.FormValue("username")
	info := c.FormValue("info")

	type State struct {
		Status string
	}

	if info != "username" && info != "email" && info != "password" {
		fmt.Printf("ChangeInfoRequest error 0: No valid info : %s\n", info)
		return c.JSON(http.StatusUnauthorized, State{
			Status: "Error",
		})
	}

	cookie, err := c.Cookie("sessionId")
	if err != nil {
		fmt.Printf("ChangeInfoRequest error 1: %s\n", err)
		return c.JSON(http.StatusUnauthorized, State{
			Status: "Error",
		})
	}


	id, success, err := db.ValidateSession(cookie.Value)
	if err != nil {
		fmt.Printf("ChangeInfoRequest Error 2\n")
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
			fmt.Printf("ChangeInfoRequest Error 3: %s\n", res.Error)
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
