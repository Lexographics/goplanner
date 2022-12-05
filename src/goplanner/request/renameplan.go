package goplanner

import (
	"net/http"
	"strconv"

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
		planidInt, err := strconv.Atoi(planid)
		if err != nil {
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
