package goplanner

import (
	"fmt"
	"net/http"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/golang-jwt/jwt"
)


type User struct {
	ID            uint      `gorm:"primaryKey"`
	Username      string    `json:"name"`
	Password      string    `json:"password"`
	Email         string    `json:"email"`
}

type Plan struct {
	Id     int
	UserID int
	Plan   string
	State  string
	End    time.Time
}

type Session struct {
	Session string   `gorm:"primaryKey"`
	UserId  uint
	Start   time.Time
	Expire  time.Time
}

var Database *gorm.DB

func InitDatabase() {
	var err error

	Database, err = gorm.Open(mysql.Open("root:@/goplanner?parseTime=true"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	sqldb, err := Database.DB()
	sqldb.SetConnMaxLifetime(time.Minute * 3)
	sqldb.SetMaxOpenConns(10)
	sqldb.SetMaxIdleConns(10)

	Database.AutoMigrate(&Session{})
	Database.AutoMigrate(&Plan{})
	Database.AutoMigrate(&User{})
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

	
	// res, err := Database.Exec("INSERT INTO `sessions`(`session`, `user_id`, `start`, `expire`) VALUES (?, ?, ?, ?)", token, id, now, expires)
	session := Session{
		Session: token,
		UserId: uint(id),
		Start: now,
		Expire: expires,
	}
	res := Database.Create(&session)
	
	if err != nil {
		fmt.Printf("CreateToken error 2: %s", err)
		return http.Cookie{}, err
	}

	affected := res.RowsAffected
	if affected != 1 {
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
	var sessions []Session
	res := Database.Find(&sessions, "session = ?", token) // where session = token
	if res.Error != nil {
		fmt.Printf("ValidateSession error 1: %s\n", res.Error)
		return 0, false, res.Error
	}

	for _, session := range sessions {
		if session.Session == token {
			if time.Now().Unix() > session.Start.Unix() && time.Now().Unix() < session.Expire.Unix() {
				return int64(session.UserId), true, nil
			}
		}
	}
	fmt.Printf("ValidateSession error 2: %s\n", res.Error)
	return 0, false, nil
}


// Invalidate given tokens session
func InvalidateSession(token string) error {
	var session Session
	res := Database.Find(&session, "session = ?", token).Delete(&session)
	if res.Error != nil {
		return res.Error
	}
	return nil
}