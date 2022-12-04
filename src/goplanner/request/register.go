package goplanner

import (
	_ "database/sql"
	"log"
	"net/http"
	_ "time"

	_ "github.com/go-sql-driver/mysql"

	_ "github.com/golang-jwt/jwt/v4"

	"github.com/labstack/echo/v4"

	db "github.com/Lexographics/goplanner/src/goplanner/db"
	views "github.com/Lexographics/goplanner/src/goplanner/views"
)

// POST
func RegisterRequest(c echo.Context) error {
	user := db.User{}
	
	user.Name = c.FormValue("name")
	user.Password = c.FormValue("password")
	user.Email = c.FormValue("mail")

	log.Printf("Register with name:%s, pwd:%s, mail:%s", user.Name, user.Password, user.Email)

	res, err := db.Database.Exec("INSERT INTO `users`(`username`, `email`, `password`) VALUES (?, ?, ?)", user.Name, user.Email, user.Password)
	if err != nil {
		log.Printf("Error creating user " + err.Error())
		return views.RedirectToHomeView(c)
	}

	affected, err := res.RowsAffected()
	if affected == 1 {
		log.Printf("new user: %#v", user)
		id, err := res.LastInsertId()
		if err != nil {
			log.Println("Error LastInsertId: ", err)
			return c.String(http.StatusInternalServerError, "something went wrong")
		}

		cookie, err := db.CreateToken(id)
		if err != nil {
			log.Println("Error creating jwt token: ", err)
			return c.String(http.StatusInternalServerError, "something went wrong")
		}
	
		c.SetCookie(&cookie)
	} else {
		log.Printf("Unable to create new user")
	}


	return views.RedirectToHomeView(c)
}