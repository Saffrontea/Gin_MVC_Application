package register

import (
	"Gin_MVC/controller/login"
	"Gin_MVC/model/database"
	"Gin_MVC/model/discuss"
	"Gin_MVC/model/notify"
	"Gin_MVC/model/priority"
	"Gin_MVC/model/user"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func DisplayRegister(c *gin.Context) {
	session := sessions.Default(c)
	errorMsg := ""
	_, loginState, err := login.GetLoginUser(c)
	if loginState {
		session.Set("err", "すでにログインされています")
		session.Save()
		c.Redirect(302, "/")
	}
	if err != nil {
		errorMsg = err.Error()
	}
	e, _ := c.Get("err")
	if e != "" {
		errorMsg = e.(string)
	}
	c.HTML(200, "register.html", gin.H{
		//ヘッダーのユーザー情報
		"headerUser": nil,
		//ログイン状態
		"loginState": false,
		"errorMsg":   errorMsg,
	})
}

func RegisterUser(c *gin.Context) {
	l, e := strconv.Atoi(c.PostForm("location"))
	p := c.PostForm("publish") == "on"
	if e != nil {
		l = 50
	}
	var u = user.User{
		//UserId:       0,
		Name:     c.PostForm("name"),
		Ruby:     c.PostForm("ruby"),
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
		Tel:      c.PostForm("tel"),
		Location: uint32(l),
		Publish:  p,
		Notify:   notify.Notify{},
		Priority: priority.Priority{Priority: 100},
		Discuss:  []discuss.Discuss{},
		Star:     []user.Star{},
	}
	err := database.Transaction(func(tx *gorm.DB) error {
		e := user.CreateUser(&u, tx).Error
		return e
	})
	if err != nil {
		c.Set("err", "会員登録に失敗しました")
		c.Redirect(http.StatusTemporaryRedirect, "/register")
	}
}
