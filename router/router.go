package router

import (
	"Gin_MVC/controller/decree"
	"Gin_MVC/controller/discuss"
	"Gin_MVC/controller/index"
	"Gin_MVC/controller/login"
	"Gin_MVC/controller/news"
	"Gin_MVC/controller/notify"
	"Gin_MVC/controller/profile"
	"Gin_MVC/controller/register"
	"Gin_MVC/controller/search"
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
	router.StaticFile("/favicon.ico", "resource/img/favicon.ico")
	router.StaticFile("/icon.png", "resource/img/icon.png")
	router.StaticFile("/icon-dark.png", "resource/img/icon-dark.png")
	router.Static("resource", "./resource")
	router.GET("/news", news.Display)
	router.GET("/", index.Display)
	router.GET("/discuss", discuss.Display)
	router.POST("/createDiscuss", discuss.CreateDiscuss)
	router.POST("/addDiscuss", discuss.AddDiscuss)
	router.POST("/closeDiscuss", discuss.CloseDiscuss)
	router.GET("/discussList", discuss.GetDiscussList)
	router.GET("/login", login.Display)
	router.GET("/logout", login.LogoutUser)
	router.GET("/login/passwordForget", login.DisplayForgetPage)
	router.POST("/login/sendResetMail", login.SendForgetPasswordMail)
	router.GET("/login/reset", login.DisplayResetPassword)
	router.POST("/login/doResetPassword", login.DoResetPassword)
	router.POST("/doAuth", login.DoAuth)
	router.GET("/notify", notify.GetNotify)
	router.GET("/profile", profile.DisplayMyProfile)
	router.GET("/userprofile", profile.DisplayUserProfile)
	router.GET("/editProfile", profile.DisplayEditProfile)
	router.POST("/updateProfile", profile.UpdateProfile)
	router.GET("/decree", decree.Display)
	router.GET("/decreeList", decree.DisplayDecreeList)
	router.GET("/addStar", decree.AddStar)
	router.GET("/removeStar", decree.RemoveStar)
	router.GET("/register", register.DisplaySelectUserAccount)
	router.GET("/registerInput", register.Display)
	router.GET("/register/sendMail", register.SendMail)
	router.POST("/regUser", register.RegisterUser)
	router.POST("/deleteUser", register.DeleteUser)
	router.GET("/register/complete", register.Complete)
	router.GET("/search", search.Search)
	router.GET("/notifyList", notify.Display)
	return router
}
