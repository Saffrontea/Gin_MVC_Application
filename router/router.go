package router

import (
	"Gin_MVC/controller/decree"
	"Gin_MVC/controller/discuss"
	"Gin_MVC/controller/index"
	"Gin_MVC/controller/login"
	"Gin_MVC/controller/notify"
	"Gin_MVC/controller/profile"
	"Gin_MVC/controller/register"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func GetRouter() *gin.Engine {
	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("graduate", store))
	//router.LoadHTMLGlob("view/*.html")
	router.LoadHTMLGlob("view/*/*.html")
	router.Static("resource", "./resource")

	router.GET("/", index.Display)
	router.GET("/discuss", discuss.Display)
	router.GET("/login", login.Display)
	router.POST("/doAuth", login.DoAuth)
	router.GET("/notify", notify.GetNotify)
	router.GET("/profile", profile.Display)
	router.GET("/editProfile", profile.DisplayEditProfile)
	router.GET("/decree", decree.Display)
	router.GET("/register", register.Display)
	return router
}
