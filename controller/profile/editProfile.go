package profile

import (
	"Gin_MVC/controller/header"
	"Gin_MVC/controller/login"
	"Gin_MVC/model/database"
	"Gin_MVC/model/location"
	"Gin_MVC/model/user"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
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
			//ユーザー名
			Name string
			//ユーザーID
			UserID string
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
			usr.UserID,
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

func UpdateProfile(c *gin.Context) {
	usr, b, err := login.GetLoginUser(c)
	if err != nil || !b {

	}
	img, err := c.FormFile("img")
	if err == nil {
		file, err := img.Open()
		if err != nil {
			log.Println(err)
		}
		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
			}
		}(file)
		if err == nil {
			decode, _, err := image.Decode(file)
			if err == nil {
				//TODO:Error Handle
				usr.Image = user.Image(usr.Image.SaveImage(decode))
				log.Println(usr.Image)
			} else {
				log.Println(err)
			}

		}
	} else {
		log.Println(err)
	}
	//TODO:Update Profile

	usr.Name = c.PostForm("username")
	session := sessions.Default(c)
	session.Set("User", usr.UserID)
	session.Save()
	log.Println("aaa")
	usr.Profile = c.PostForm("profile")
	atoi, err := strconv.Atoi(c.PostForm("location"))
	if err != nil {
		log.Println("parse error?")
		return
	}
	usr.Location = uint32(atoi)
	usr.Tel = c.PostForm("tel")
	usr.Publish = c.PostForm("publish") == "on"
	log.Println(usr)
	if database.DB.Save(&usr).Error != nil {
		c.AbortWithStatus(502)
	}
	c.Redirect(http.StatusSeeOther, "profile")
}
