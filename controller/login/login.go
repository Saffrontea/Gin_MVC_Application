/*
	login.go
	ログイン関連処理
*/

package login

import (
	"Gin_MVC/model/database"
	"Gin_MVC/model/user"
	"bytes"
	"errors"
	"github.com/google/uuid"
	"log"
	"net/http"
	_ "net/http"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Display ログインフォーム表示
func Display(c *gin.Context) {
	session := sessions.Default(c)
	e := session.Get("err")
	if e != nil {
		session.Delete("err")
		session.Save()
	}

	c.HTML(200, "login.html", gin.H{
		"err": e,
	})
}

//DoAuth
/*
	ログイン処理
 	成功時 : /にリダイレクト,sessionにUserを保存
 	失敗時 : /loginにリダイレクト,sessionにerrを保存
*/
func DoAuth(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("err")
	if session.Get("user") != nil {
		c.Next()
	}
	postedUser := c.PostForm("username")
	user, _ := user.GetUser(postedUser)
	formPass := c.PostForm("password")
	log.Println(formPass + " : " + user.Password)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(formPass)); err != nil || user.UserID == "" || user.Password == "" {
		log.Println("Cannot login")
		c.Abort()
		session.Set("err", "ログインに失敗しました")
		session.Save()
		c.Redirect(302, "/login")
	} else {
		session.Set("User", user.UserID)
		session.Set("LoginTime", time.Now().Unix())
		session.Save()
		c.Next()
		//session.Set("length")
		log.Println("Logined User: " + user.UserID)
		c.Redirect(302, "/")
	}
}

//GetLoginUser ログインしているユーザーを返す関数
func GetLoginUser(c *gin.Context) (*user.User, bool, error) {
	session := sessions.Default(c)
	sessionUser := session.Get("User") //will be nil
	usr := user.User{}
	loginState := false
	errorMsg := ""
	if sessionUser != nil {
		var err error
		usr, err = user.GetUser(sessionUser.(string))
		loginState = true
		if err != nil {
			log.Println("error")
			usr = user.User{}
			loginState = false
			errorMsg = "不明なエラーが発生しました"
			session.Clear()
			return &usr, loginState, errors.New(errorMsg)
		}
	} else {
		return &user.User{}, loginState, errors.New("")
	}
	return &usr, loginState, nil
}

func LogoutUser(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	//session.Delete("User")
	//session.Delete("LoginTime")
	err := session.Save()
	if err != nil {
		c.Error(err)
		return
	}
	c.Redirect(http.StatusSeeOther, "/")
}

func DisplayForgetPage(c *gin.Context) {
	s := sessions.Default(c)
	c.HTML(http.StatusOK, "forgetPassword.html", gin.H{
		"err": s.Get("err"),
	})
}

func SendForgetPasswordMail(c *gin.Context) {
	mail := c.PostForm("email")
	usr, err := user.GetUserFromEmail(mail)
	if err != nil || usr.Id == 0 {
		s := sessions.Default(c)
		s.Set("err", "メールアドレスが存在しません")
		s.Save()
		c.Redirect(http.StatusSeeOther, "/login/passwordForget")
		return
	}
	key := strings.Replace(uuid.NewString(), "-", "", -1)
	user.PasswordForgetUsers[key] = usr.Id
	tmpl, err := template.ParseFiles(path.Join("view", "mail", "password.tmpl"))
	if err != nil {
		log.Println(err)
	}
	var b bytes.Buffer
	err = tmpl.Execute(&b, key)
	if err != nil {
		log.Println(err)
	}
	log.Println(b.String())
	c.HTML(http.StatusOK, "sendPasswordForgetMail.html", gin.H{})
}

func DisplayResetPassword(c *gin.Context) {
	q := c.Query("q")
	uid, exist := user.PasswordForgetUsers[q]
	if !exist {
		c.HTML(http.StatusForbidden, "403.html", gin.H{})
		return
	}
	usr, err := user.GetUserByID(uid)
	if err != nil {
		c.HTML(http.StatusForbidden, "403.html", gin.H{})
		return
	}
	c.HTML(http.StatusOK, "resetPassword.html", gin.H{
		"usr": usr,
		"q":   q,
	})
}
func DoResetPassword(c *gin.Context) {
	q := c.Query("q")
	uid, exist := user.PasswordForgetUsers[q]
	if !exist {
		c.HTML(http.StatusForbidden, "403.html", gin.H{})
		return
	}
	usr, err := user.GetUserByID(uid)
	if err != nil {
		c.HTML(http.StatusForbidden, "403.html", gin.H{})
		return
	}
	usr.Password = c.PostForm("password")
	err = usr.ChangePassword()
	if err != nil {
		c.AbortWithStatus(502)
		return
	}
	database.DB.Save(&usr)
	delete(user.PasswordForgetUsers, q)
	c.HTML(http.StatusOK, "resetPasswordComplete.html", gin.H{})
}
