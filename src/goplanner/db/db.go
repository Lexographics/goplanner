package goplanner

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)


type User struct {
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
}

type JwtClaims struct {
	Name     string  `json:"name"`
	jwt.StandardClaims
}

type Plan struct {
	Id     int
	UserID int
	Plan   string
	Start  time.Time
	End    time.Time
}

var Database *sql.DB

func InitDatabase() {
	var err error
	Database, err = sql.Open("mysql", "root:@/goplanner?parseTime=true")

	if err != nil {
		panic(err)
	}

	Database.SetConnMaxLifetime(time.Minute * 3)
	Database.SetMaxOpenConns(10)
	Database.SetMaxIdleConns(10)
}

func CreateToken(id int64) (http.Cookie, error) {
	now := time.Now()
	expires := time.Now().Add(24 * time.Hour)

	claims := jwt.StandardClaims{
		Id: "user",
		ExpiresAt: expires.Unix(),
	}
	
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token, err := rawToken.SignedString([]byte("secret")) // ! temporary !
	if err != nil {
		fmt.Println("CreateToken error 1")
		return http.Cookie{}, err
	}

	res, err := Database.Exec("INSERT INTO `sessions`(`session`, `user_id`, `start`, `expire`) VALUES (?, ?, ?, ?)", token, id, now, expires)
	if err != nil {
		fmt.Printf("CreateToken error 2: %s", err)
		return http.Cookie{}, err
	}

	_, err = res.RowsAffected()
	if err != nil {
		fmt.Println("CreateToken error 3")
		return http.Cookie{}, err
	}

	cookie := http.Cookie{
		Name: "sessionId",
		Value: token,
		Expires: expires,
	}
	return cookie, nil
}


// Returns true if a session is valid
func ValidateSession(token string) (int64, bool, error){
	res, err := Database.Query("SELECT `user_id`, `start`, `expire` FROM `sessions` WHERE session=?", token)
	if err != nil {
		fmt.Printf("ValidateSession error 1: %s\n", err)
		return 0, false, err
	}

	var id int64
	var start time.Time
	var expire time.Time
	for res.Next() {
		res := res.Scan(&id, &start, &expire)
		if res != nil {
			fmt.Printf("ValidateSession error 2: %s\n", err)
			return 0, false, err
		}

		if time.Now().Unix() > start.Unix() && time.Now().Unix() < expire.Unix() {
			return id, true, nil
		}

	}
	
	return 0, false, nil
}


// Invalidate given tokens session
func InvalidateSession(token string) error {
	_, err := Database.Exec("DELETE FROM `sessions` WHERE session=?", token)
	if err != nil {
		return err
	}
	return nil
}