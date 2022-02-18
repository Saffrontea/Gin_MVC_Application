package profile

import (
	"Gin_MVC/controller/header"
	"Gin_MVC/controller/login"
	"Gin_MVC/model/location"
	"Gin_MVC/model/user"
	"github.com/gin-gonic/gin"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
)

func DisplayEditProfile(c *gin.Context) {
	errorMsg := ""
	usr, loginState, err := login.GetLoginUser(c)
	locList := location.GetLocationList()
	if err != nil {
		errorMsg = err.Error()
	}

	c.HTML(200, "editProfile.html", gin.H{
		"headerUser": header.GetHeaderUser(usr),
		"user": struct {
			//本名
			Name string
			//ユーザー名
			UserName string
			//自己紹介
			Profile string
			//居住地
			Location uint32
			//公開設定
			Publish bool
			//電話番号
			Tel string
		}{
			usr.Name,
			usr.Username,
			usr.Profile,
			usr.Location,
			usr.Publish,
			usr.Tel,
		}, //パスワード等を秘匿
		//アイコン
		"img": usr.Image.GetImage(),
		//ログイン状態
		"loginState":   loginState,
		"errorMsg":     errorMsg,
		"LocationList": locList,
	})
}


func UpdateProfile(c *gin.Context){
	usr, b, err := login.GetLoginUser(c)
	if err != nil || !b {

	}
	img, err := c.FormFile("img")
	file,err := img.Open()
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)
	if err != nil {
	}

	decode, _, err := image.Decode(file)
	if err != nil {
		//TODO:Error Handle
	}
	usr.Image = user.Image(user.SaveImage(decode))
}