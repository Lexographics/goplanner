package goplanner

import (
	"io"
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	db "github.com/Lexographics/goplanner/src/goplanner/db"
	request "github.com/Lexographics/goplanner/src/goplanner/request"
	views "github.com/Lexographics/goplanner/src/goplanner/views"
)


type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}


func Init()  {
	db.InitDatabase()
	
	// Echo
	e := echo.New()
	
	t := &Template{
		templates: template.Must(template.ParseGlob("./static/*.html")),
	}
	e.Renderer = t

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}] ${status} ${method} ${host}${path} ${latency_human}` + "\n",
	}))
	e.Use(middleware.Recover())
	e.Static("/", "./static")
	
	authGroup := e.Group("/home")
	jwtMiddleware := middleware.JWTWithConfig(
		middleware.JWTConfig{
			SigningMethod: "HS512",
			SigningKey: []byte("secret"), // ! temporary !
			TokenLookup: "cookie:sessionId",
			ErrorHandlerWithContext: func(err error, c echo.Context) error {
				return c.Redirect(http.StatusOK, "/auth")
			},
		},
	)
	authGroup.Use(jwtMiddleware)

	authGroup.GET("/", views.HomeView)

	e.GET("/", views.HomeView)
	e.GET("/profile", views.ProfileView)
	e.GET("/auth", views.AuthView)
	
	e.POST("/login", request.LoginRequest)
	e.POST("/register", request.RegisterRequest)
	e.GET("/changeinfo", request.ChangeInfoRequest)
	e.GET("/logout", request.LogoutRequest)
	e.GET("/newplan", request.NewPlanRequest)
	e.GET("/deleteplan", request.DeletePlanRequest)
	e.GET("/renameplan", request.RenamePlanRequest)
	e.GET("/setplanstate", request.SetPlanStateRequest)


	e.Logger.Fatal(e.Start(":8000"))
}