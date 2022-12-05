package goplanner

import (
	"net/http"
	"strconv"
	"time"

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
		dateInt, err := strconv.Atoi(expire)
		date := time.Unix(int64(dateInt), 0)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, State{
				Status: "Error",
			})
		}

		plan := db.Plan{
			UserID: int(id),
			Plan: text,
			End: date,
		}
		res := db.Database.Create(&plan)
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
