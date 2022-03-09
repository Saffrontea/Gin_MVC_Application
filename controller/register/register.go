package register

import (
	"Gin_MVC/controller/login"
	"Gin_MVC/model/database"
	"Gin_MVC/model/discuss"
	"Gin_MVC/model/location"
	"Gin_MVC/model/notify"
	"Gin_MVC/model/priority"
	"Gin_MVC/model/user"
	"bytes"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

func Display(c *gin.Context) {
	session := sessions.Default(c)
	errorMsg := ""
	e := session.Get("err")
	session.Delete("err")
	session.Save()
	_, loginState, err := login.GetLoginUser(c)
	if loginState {
		session.Set("err", "すでにログインされています")
		session.Save()
		c.Redirect(302, "/")
	}
	if err != nil {
		errorMsg = err.Error()
	}
	c.HTML(200, "registerInput.html", gin.H{
		//ヘッダーのユーザー情報
		"headerUser": nil,
		//ログイン状態
		"loginState":   false,
		"e":            e,
		"errorMsg":     errorMsg,
		"LocationList": location.GetLocationList(),
	})
}

func DisplaySelectUserAccount(c *gin.Context) {
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
	e := session.Get("err")
	if e != nil && e.(string) != "" {
		errorMsg = e.(string)
	}
	c.HTML(200, "register.html", gin.H{
		"err":      e,
		"errorMsg": errorMsg,
	})
}

func RegisterUser(c *gin.Context) {
	if c.PostForm("name") == "" || c.PostForm("ruby") == "" || c.PostForm("username") == "" || c.PostForm("password") == "" || c.PostForm("tel") == "" {
		c.Set("err", "未入力の項目があります")
		c.Abort()
		c.Redirect(http.StatusSeeOther, "/register")
	}

	l, e := strconv.Atoi(c.PostForm("location"))
	p := c.PostForm("publish") == "on"
	if e != nil {
		l = 50
	}
	//var u user.User
	//c.Bind(&u)
	birth, err := time.Parse(time.RFC3339, c.PostForm("birth")+"T00:00:00+09:00")
	if err != nil {
		log.Println(err)
		return
	}
	var u = user.User{
		//UserId:       0,
		Name:     c.PostForm("username"),
		UserID:   c.PostForm("userid"),
		Password: c.PostForm("password"),
		Tel:      c.PostForm("tel"),
		Mail:     c.PostForm("email"),
		Authed:   false,
		Birth:    &birth,
		Location: uint32(l),
		Publish:  p,
		Notify: notify.Notify{
			Notify: "[]",
		},
		Priority: priority.Priority{Priority: 100},
		Discuss:  []discuss.Discuss{},
		Star:     []user.Star{},
	}
	log.Println(u)
	if u.Name == "" && u.UserID == "" && u.Password == "" {
		log.Println("error")
		return
	}
	err = database.Transaction(func(tx *gorm.DB) error {
		e := user.CreateUser(&u, tx).Error
		return e
	})
	if err != nil {
		s := sessions.Default(c)
		str := ""
		if strings.Contains(err.Error(), "mail") {
			str = str + " 使用済みメールアドレスです。"
		}
		if strings.Contains(err.Error(), "user") {
			str = str + " 使用済みユーザーIDです。"
		}
		s.Set("err", "会員登録に失敗しました"+str)
		s.Save()
		c.Redirect(http.StatusSeeOther, "/registerInput")
	} else {
		key := strings.Replace(uuid.NewString(), "-", "", -1)
		user.WaitAuthUsers[key] = u.Id
		c.Redirect(302, "/register/sendMail?q="+key)
	}
}

func SendMail(c *gin.Context) {
	//key, exist := c.Get("mail")
	//log.Println(key)
	//c.Set("mail", "")
	key := c.Query("q")
	if key == "" {
		c.HTML(403, "403.html", gin.H{})
		return
	}
	tmpl, err := template.ParseFiles(path.Join("view", "mail", "register.tmpl"))
	if err != nil {
		log.Println(err)
	}
	var b bytes.Buffer
	err = tmpl.Execute(&b, key)
	if err != nil {
		log.Println(err)
	}
	log.Println(b.String())
	c.HTML(200, "sendMail.html", gin.H{})
}
func Complete(c *gin.Context) {
	q := c.Query("q")
	log.Println(q)
	query, exist := user.WaitAuthUsers[q]
	log.Println(query)
	if !exist {
		c.HTML(404, "404.html", gin.H{})
		return
	}
	disabledUser, err := user.GetDisabledUser(query)
	log.Println(disabledUser.UserID)
	if err != nil {
		c.HTML(404, "404.html", gin.H{})
		return
	}
	disabledUser.Authed = true
	if err = database.DB.Save(&disabledUser).Error; err != nil {
		c.AbortWithStatus(502)
		return
	}
	delete(user.WaitAuthUsers, q)
	c.HTML(200, "complete.html", gin.H{})
}

func DeleteUser(c *gin.Context) {
	session := sessions.Default(c)

	usr, loginState, err := login.GetLoginUser(c)
	if !loginState || err != nil {
		session.Set("err", "ログインしていません")
		session.Save()
		c.Redirect(302, "/")
	}

	if c.PostForm("agree") == "on" {
		t := time.UnixMilli(0)
		u := user.User{
			Id:       usr.Id,
			Name:     "退会ユーザー",
			UserID:   uuid.NewString(),
			Password: "",
			Mail:     uuid.NewString(),
			Authed:   false,
			Birth:    &t,
			Tel:      "",
			Location: 0,
			Publish:  false,
			Profile:  "",
			Image:    "",
		}
		database.DB.Save(&u)
		c.Redirect(http.StatusSeeOther, "/logout")
	}
}
