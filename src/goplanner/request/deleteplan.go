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
func DeletePlanRequest(c echo.Context) error {
	planid := c.FormValue("id")

	type State struct {
		Status string
	}

	cookie, err := c.Cookie("sessionId")
	if err != nil {
		fmt.Printf("DeletePlanRequest error 1: %s\n", err)
		return c.JSON(http.StatusUnauthorized, State{
			Status: "Error",
		})
	}


	id, success, err := db.ValidateSession(cookie.Value)
	if err != nil {
		fmt.Printf("DeletePlanRequest Error 2\n")
		return c.JSON(http.StatusUnauthorized, State{
			Status: "Unauthorized",
		})
	}

	if success {
		var plan db.Plan
		res := db.Database.Find(&plan, "id = ? AND user_id = ?", planid, id).Delete(&plan)
		if res.Error != nil {
			fmt.Printf("DeletePlanRequest Error 3: %s\n", res.Error)
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
