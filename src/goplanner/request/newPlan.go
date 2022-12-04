package goplanner

import (
	_ "database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"
	_ "time"

	_ "github.com/go-sql-driver/mysql"

	_ "github.com/golang-jwt/jwt/v4"

	"github.com/labstack/echo/v4"

	db "github.com/Lexographics/goplanner/src/goplanner/db"
)

// GET
func NewPlanRequest(c echo.Context) error {
	text := c.FormValue("text")
	expire := c.FormValue("expire")

	type State struct {
		Status string
	}

	cookie, err := c.Cookie("sessionId")
	if err != nil {
		fmt.Printf("NewPlanRequest error 1: %s\n", err)
		return c.JSON(http.StatusUnauthorized, State{
			Status: "Error",
		})
	}


	id, success, err := db.ValidateSession(cookie.Value)
	if err != nil {
		fmt.Printf("NewPlanRequest Error 2\n")
		return c.JSON(http.StatusUnauthorized, State{
			Status: "Unauthorized",
		})
	}

	if success {
		dateInt, err := strconv.Atoi(expire)
		date := time.Unix(int64(dateInt), 0)
		if err != nil {
			fmt.Printf("Invalid date value: %s\n", err)
			return c.JSON(http.StatusUnauthorized, State{
				Status: "Error",
			})
		}

		_, err = db.Database.Exec("INSERT INTO `plans`(`user_id`, `plan`, `end`) VALUES (?, ?, ?)", id, text, date)
		if err != nil {
			fmt.Printf("NewPlanRequest Error 3: %s\n", err)
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
