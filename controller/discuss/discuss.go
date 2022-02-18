package discuss

import (
	"Gin_MVC/controller/header"
	"Gin_MVC/controller/login"
	"Gin_MVC/model/database"
	"Gin_MVC/model/discuss"
	"Gin_MVC/model/user"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func addDiscuss(c *gin.Context) {
	decree := c.Param("decree")
	if decree == "" {
		return
	}
}

func Display(c *gin.Context) {

	errorMsg := ""
	usr, loginState, err := login.GetLoginUser(c)
	if err != nil {
		errorMsg = err.Error()
	}

	discussId := c.Query("id")
	id, err := strconv.Atoi(discussId)
	if err != nil {
		log.Println("cannot parse ID :", discussId)
		c.Error(err)
	}
	//議論取得
	dis, err := discuss.GetDiscuss(id)

	decname := dis.Decree.Name
	disType := dis.Discuss_Type
	disTitle := dis.Title
	//作成ユーザー
	var u user.User
	er := database.DB.Find(&u, dis.Create_User, "id = ?").Error //Select * from user where user.id = Create_User
	if er != nil {
		u = user.User{
			Id:   -1,
			Name: "退会ユーザー",
		}
	}
	var content discuss.ContentJSON
	e := json.Unmarshal(dis.Content, &content)
	if e != nil {
	}
	c.HTML(200, "discuss.html", gin.H{
		"headerUser": header.GetHeaderUser(usr),
		//ログイン状態
		"loginState": loginState,
		"errorMsg":   errorMsg,

		//議論している法令名
		"decree": decname,
		//議論のタイトル
		"discuss": disTitle,
		//議論のカテゴリ(疑問、議論、相談)
		"discussType": disType,
		//作成ユーザー
		"createUser": u,
		"content":    content,
	})
}

func createDiscuss(c *gin.Context) {
	//errorMsg := ""
	usr, _, err := login.GetLoginUser(c)
	if err != nil {
		//errorMsg = err.Error()
	}
	discussType,err  := strconv.Atoi(c.PostForm("type"))
	content := c.PostForm("content")
	var d = discuss.Discuss{
		//Decree:       ,
		Create_User: usr.Id,
		Discuss_Type: discussType,
		Title:       c.PostForm("title"),
		Opened:      1,
		//Content: ,
	}
	j := discuss.ContentJSON{}
	j.UpdateContent(d.Title,d.Create_User,content)
	contentJSON,_ := json.Marshal(j)
	d.Content = contentJSON

	err = discuss.CreateDiscuss(d)
	if err != nil {
		return
	}

}

