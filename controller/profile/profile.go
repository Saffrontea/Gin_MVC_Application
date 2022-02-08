package profile

import (
	"Gin_MVC/controller/header"
	"Gin_MVC/controller/login"
	"Gin_MVC/model/decree"
	"github.com/gin-gonic/gin"
)

func ProfileDisplay(c *gin.Context) {
	var decName []decree.Decree
	errorMsg := ""
	usr, loginState, err := login.GetLoginUser(c)
	if err != nil {
		errorMsg = err.Error()
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
	c.HTML(200, "profile.html", gin.H{
		"headerUser": header.GetHeaderUser(usr),
		//ユーザー名
		"username": usr.GetUserName(),
		//自己紹介
		"userprofile": usr.Profile,
		//お気に入り
		"stars": decName,
		//アイコン
		"img": img,
		//ログイン状態
		"loginState": loginState,
		"errorMsg":   errorMsg,
	})
}
