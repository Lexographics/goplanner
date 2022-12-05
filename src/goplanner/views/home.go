package goplanner

import (
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/labstack/echo/v4"

	db "github.com/Lexographics/goplanner/src/goplanner/db"
)


func HomeView(c echo.Context) error {

	cookie, err := c.Cookie("sessionId")
	if err != nil {
		return RedirectToAuthView(c)
	}

	id, success, err := db.ValidateSession(cookie.Value)
	if err != nil {
		return RedirectToAuthView(c)
	}

	if success {
		var plans []db.Plan
		var res *gorm.DB

		start := c.QueryParam("start")
		end := c.QueryParam("end")

		startDate := time.Now()
		endDate := time.Now()

		startInt, err := strconv.Atoi(start)
		if err == nil {
			startDate = time.Unix(int64(startInt), 0)
		}

		endInt, err := strconv.Atoi(end)
		if err == nil {
			endDate = time.Unix(int64(endInt), 0)
		}
		

		if start == "" && end == "" {
			res = db.Database.Order("end asc").Find(&plans, "user_id = ?", id)
		} else {
			if start != "" && end == "" { // has start
				res = db.Database.Order("end asc").Find(&plans, "user_id = ? AND end >= ?", id, startDate)
			} else if start == "" && end != "" { // has end
				res = db.Database.Order("end asc").Find(&plans, "user_id = ? AND end <= ?", id, endDate)
			} else { // has both
				res = db.Database.Order("end asc").Find(&plans, "user_id = ? AND end >= ? AND end <= ?", id, startDate, endDate)
			}
		}

		if res.Error != nil {
			return RedirectToAuthView(c)
		}

		type Page struct {
			Id      string
			Plans   []db.Plan
		}
		page := Page{Id: string(rune(id)), Plans: plans}

		return c.Render(http.StatusOK, "PlansPage", page)
	}
	
	return RedirectToAuthView(c)
}



func RedirectToHomeView(c echo.Context) error {
	return c.File("./static/toHome.html")
}

func RedirectToAuthView(c echo.Context) error {
	return c.File("./static/toAuth.html")
}