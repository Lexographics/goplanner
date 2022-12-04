package goplanner

import (
	_ "database/sql"
	"fmt"
	_ "time"

	_ "github.com/go-sql-driver/mysql"

	_ "github.com/golang-jwt/jwt/v4"

	"github.com/labstack/echo/v4"

	db "github.com/Lexographics/goplanner/src/goplanner/db"
	views "github.com/Lexographics/goplanner/src/goplanner/views"
)

// POST
func LoginRequest(c echo.Context) error {
	name := c.FormValue("name")
	password := c.FormValue("password")

	
	res, err := db.Database.Query("SELECT `id` FROM `users` WHERE username=? AND password=?", name, password)
	if err != nil {
		fmt.Printf("Error 1 \n")
		return views.RedirectToHomeView(c)
	}

	defer res.Close()

	var id int64
	user := 0
	for res.Next() {
		err := res.Scan(&id)
		user += 1

		if err != nil {
			fmt.Printf("Error 2 \n")
			return views.RedirectToHomeView(c)
		}
	}

	if user == 1 {
		
		cookie, err := db.CreateToken(id)
		if err != nil {
			fmt.Printf("Error 3 \n")
			return views.RedirectToHomeView(c)
		}
		
		c.SetCookie(&cookie)
		return views.RedirectToHomeView(c)
	} else {
		fmt.Printf("Error 4 \n")
		return views.RedirectToHomeView(c)
	}
}
