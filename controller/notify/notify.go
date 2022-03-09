package notify

import (
	"Gin_MVC/controller/header"
	"Gin_MVC/controller/login"
	"Gin_MVC/model/database"
	"Gin_MVC/model/discuss"
	"Gin_MVC/model/notify"
	"Gin_MVC/model/user"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

//GetNotify 通知取得
func GetNotify(c *gin.Context) {
	user, _, err := login.GetLoginUser(c)
	log.Println(user.UserID)

	log.Println("notify", err)
	if err != nil {
		c.JSON(403, "")
	} else {
		notifies := user.Notify
		log.Println(user.Notify)
		c.JSON(http.StatusAccepted, notifies.Notify)
	}
}

func SendNotify(d discuss.Discuss) {

	var content discuss.ContentJSON
	json.Unmarshal([]byte(d.Content), &content)
	for _, i := range content[len(content)-1].MentionTo {
		id, err := user.GetUserByID(i)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		cusr, err := user.GetUserByID(content[len(content)-1].Create_User)
		if id.Id == cusr.Id {
			continue
		}
		if err != nil {
			cusr = user.User{
				Name:    "退会ユーザー",
				UserID:  "退会ユーザー",
				Birth:   nil,
				Profile: "",
				Image:   "",
			}
		}
		id.Notify = id.Notify.AddNotify(notify.NotifyJSON{
			{
				DiscussID: d.Id,
				Hash:      uuid.NewString(),
				Time:      time.Now().Format("2006-01-02"),
				Level:     0,
				Comment:   d.Title + "に" + cusr.GetUserName() + "さんがコメントしました",
			},
		})
		id.Notify.Id = id.Id
		database.DB.Save(&id.Notify)
		log.Println("notify", id.Notify)
	}
}

func Display(c *gin.Context) {
	errorMsg := ""
	usr, loginState, err := login.GetLoginUser(c)
	if err != nil {
		errorMsg = err.Error()
	}
	c.HTML(http.StatusOK, "notify.html", gin.H{
		"headerUser": header.GetHeaderUser(usr),
		//ログイン状態
		"loginState": loginState,
		"errorMsg":   errorMsg,
	})
}
