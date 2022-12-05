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
func SetPlanStateRequest(c echo.Context) error {
	planid := c.FormValue("id")
	state := c.FormValue("state")

	type State struct {
		Status string
	}

	cookie, err := c.Cookie("sessionId")
	if err != nil {
		fmt.Printf("SetPlanStateRequest error 1: %s\n", err)
		return c.JSON(http.StatusUnauthorized, State{
			Status: "Error",
		})
	}


	id, success, err := db.ValidateSession(cookie.Value)
	if err != nil {
		fmt.Printf("SetPlanStateRequest Error 2\n")
		return c.JSON(http.StatusUnauthorized, State{
			Status: "Unauthorized",
		})
	}

	if success {

		_, err = db.Database.Exec("UPDATE `plans` SET `state`=? WHERE id=? AND user_id=?", state, planid, id)
		if err != nil {
			fmt.Printf("SetPlanStateRequest Error 3: %s\n", err)
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
