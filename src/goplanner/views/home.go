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


func HomeView(c echo.Context) error {

	cookie, err := c.Cookie("sessionId")
	if err != nil {
		fmt.Printf("HomeView error 1: %s\n", err)
		return RedirectToAuthView(c)
	}

	id, success, err := db.ValidateSession(cookie.Value)
	if err != nil {
		fmt.Printf("HomeView error 2: %s\n", err)
		return RedirectToAuthView(c)
	}

	if success {
		var plans []db.Plan
		res := db.Database.Find(&plans, "user_id = ?", id)
		if res.Error != nil {
			fmt.Printf("HomeView error 3: %s\n", res.Error)
			return RedirectToAuthView(c)
		}

		type Page struct {
			Id      string
			Plans   []db.Plan
		}
		page := Page{Id: string(rune(id)), Plans: plans}

		return c.Render(http.StatusOK, "PlansPage", page)
	}
	

	fmt.Printf("HomeView error 5: %s\n", err)
	return RedirectToAuthView(c)
}



func RedirectToHomeView(c echo.Context) error {
	return c.File("./static/toHome.html")
}

func RedirectToAuthView(c echo.Context) error {
	return c.File("./static/toAuth.html")
}