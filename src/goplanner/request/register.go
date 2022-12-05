package goplanner

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	db "github.com/Lexographics/goplanner/src/goplanner/db"
	views "github.com/Lexographics/goplanner/src/goplanner/views"
)

// POST
func RegisterRequest(c echo.Context) error {
	user := db.User{}
	
	user.Username = c.FormValue("name")
	user.Password = c.FormValue("password")
	user.Email = c.FormValue("mail")
	log.Printf("Register with name:%s, mail:%s", user.Username, user.Email)


	res := db.Database.Create(&user)
	if res.Error != nil {
		log.Printf("Error creating user " + res.Error.Error())
		return views.RedirectToHomeView(c)
	}

	affected := res.RowsAffected
	if affected == 1 {
		log.Printf("new user: %#v", user)
		
		cookie, err := db.CreateToken(int64(user.ID))
		if err != nil {
			return c.String(http.StatusInternalServerError, "something went wrong")
		}
	
		c.SetCookie(&cookie)
	} else {
		log.Printf("Unable to create new user")
	}


	return views.RedirectToHomeView(c)
}