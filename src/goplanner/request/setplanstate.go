package goplanner

import (
	"net/http"

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
		var plan db.Plan
		res := db.Database.Find(&plan, "id = ? AND user_id = ?", planid, id)
		if res.Error != nil {
			return c.JSON(http.StatusUnauthorized, State{
				Status: "Unauthorized",
			})
		}
		plan.State = state

		res = db.Database.Save(&plan)
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
