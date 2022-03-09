package profile

import (
	"Gin_MVC/controller/header"
	"Gin_MVC/controller/login"
	"Gin_MVC/model/decree"
	"Gin_MVC/model/location"
	"Gin_MVC/model/user"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func DisplayMyProfile(c *gin.Context) {
	var decName []decree.Decree
	errorMsg := ""
	usr, loginState, err := login.GetLoginUser(c)
	if err != nil || !loginState {
		errorMsg = err.Error()
		c.Redirect(http.StatusSeeOther, "login")
		return
	}
	for _, star := range usr.Star {
		getDecree, err := decree.GetDecree(star.Star)
		if err == nil {
			decName = append(decName, getDecree)
		}
	}
	if usr == nil {
		//header := &user.User{}
	}
	img := usr.Image.GetImage()
	e, _ := c.Get("err")
	c.HTML(200, "profile.html", gin.H{
		"headerUser": header.GetHeaderUser(usr),
		"myProfile":  true,
		"err":        e,
		"user":       usr,
		//ユーザー名
		"username": usr.GetUserName(),
		"userid":   usr.UserID,
		//自己紹介
		"userprofile": usr.Profile,
		//お気に入り
		"stars":        decName,
		"LocationList": location.GetLocationList(),
		//アイコン
		"img": img,
		//ログイン状態
		"loginState": loginState,
		"errorMsg":   errorMsg,
	})
}

func DisplayUserProfile(c *gin.Context) {
	var decName []decree.Decree
	errorMsg := ""
	usr, loginState, err := login.GetLoginUser(c)
	if err != nil {
		errorMsg = err.Error()
	}

	if usr == nil {
		//header := &user.User{}
	}
	idstr := c.Query("userid")

	id, err := strconv.Atoi(idstr)
	if idstr == "" || (id == usr.Id && loginState) {
		c.Redirect(http.StatusSeeOther, "profile")
		return
	}
	if err != nil {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"headerUser": header.GetHeaderUser(usr),
			//ログイン状態
			"loginState": loginState,
			"errorMsg":   errorMsg,
		})
		return
	}
	userProfile, err := user.GetUserByID(id)
	if err != nil {
		c.HTML(404, "404.html", gin.H{})
		return
	}
	img := userProfile.Image.GetImage()
	for _, star := range usr.Star {
		getDecree, err := decree.GetDecree(star.Star)
		if err == nil {
			decName = append(decName, getDecree)
		}
	}
	e, _ := c.Get("err")
	c.HTML(200, "profile.html", gin.H{
		"headerUser": header.GetHeaderUser(usr),
		"myProfile":  false,
		"err":        e,
		"user":       userProfile,
		"userid":     userProfile.UserID,
		//ユーザー名
		"username": userProfile.GetUserName(),
		//自己紹介
		"userprofile": userProfile.Profile,
		//お気に入り
		"stars": decName,
		//アイコン
		"img":          img,
		"LocationList": location.GetLocationList(),
		//ログイン状態
		"loginState": loginState,
		"errorMsg":   errorMsg,
	})
}
