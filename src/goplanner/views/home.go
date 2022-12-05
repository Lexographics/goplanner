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

	fmt.Printf("Cookie: %s", cookie.Value)

	id, success, err := db.ValidateSession(cookie.Value)
	if err != nil {
		fmt.Printf("HomeView error 2: %s\n", err)
		return RedirectToAuthView(c)
	}

	if success {
		res, err := db.Database.Query("SELECT `id`, `user_id`, `plan`, 'state', `end` from `plans` WHERE `user_id` = ?", id)
		defer res.Close()
		if err != nil {
			fmt.Printf("HomeView error 3: %s\n", err)
			return RedirectToAuthView(c)
		}
		
		plans := []db.Plan{}

		for res.Next() {
			var plan db.Plan
			err := res.Scan(&plan.Id, &plan.UserID, &plan.Plan, &plan.State, &plan.End)

			if err != nil {
				fmt.Printf("HomeView error 4: %s\n", err)
				return RedirectToAuthView(c)
			}
			
			plans = append(plans, plan)
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