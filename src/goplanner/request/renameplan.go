package goplanner

import (
	_ "database/sql"
	"fmt"
	"net/http"
	"strconv"
	_ "time"

	_ "github.com/go-sql-driver/mysql"

	_ "github.com/golang-jwt/jwt/v4"

	"github.com/labstack/echo/v4"

	db "github.com/Lexographics/goplanner/src/goplanner/db"
)

// GET
func RenamePlanRequest(c echo.Context) error {
	planid := c.FormValue("id")
	newText := c.FormValue("newtext")

	type State struct {
		Status string
	}

	cookie, err := c.Cookie("sessionId")
	if err != nil {
		fmt.Printf("RenamePlanRequest error 1: %s\n", err)
		return c.JSON(http.StatusUnauthorized, State{
			Status: "Error",
		})
	}


	id, success, err := db.ValidateSession(cookie.Value)
	if err != nil {
		fmt.Printf("RenamePlanRequest Error 2\n")
		return c.JSON(http.StatusUnauthorized, State{
			Status: "Unauthorized",
		})
	}

	if success {
		planidInt, err := strconv.Atoi(planid)
		if err != nil {
			fmt.Printf("RenamePlanRequest Error 3\n")
			return c.JSON(http.StatusUnauthorized, State{
				Status: "Unauthorized",
			})
		}

		plan := db.Plan{
			Id: planidInt,
			UserID: int(id),
		}
		res := db.Database.Model(&plan).Update("plan", newText)
		if res.Error != nil {
			fmt.Printf("RenamePlanRequest Error 4: %s\n", res.Error)
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
