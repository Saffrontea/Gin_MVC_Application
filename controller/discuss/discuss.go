package discuss

import (
	"Gin_MVC/controller/header"
	"Gin_MVC/controller/login"
	"Gin_MVC/controller/notify"
	"Gin_MVC/model/database"
	"Gin_MVC/model/decree"
	"Gin_MVC/model/discuss"
	"Gin_MVC/model/user"
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"log"
	"net/http"
	"sort"
	"strconv"
)

type Content []struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Hash  string `json:"hash"`
	//作成ユーザーのID
	CreateUser user.User `json:"createUser"`
	IsMine     bool
	Body       string `json:"body"`
	//メンションユーザーのID
	MentionTo []user.User `json:"mentionTo"`
}

type DiscussList []struct {
	Id              int
	Decree          decree.Decree
	Title           string
	CreateUser      user.User
	DiscussType     int
	Opened          int
	CommentQuantity int
}

func AddDiscuss(c *gin.Context) {
	usr, loginState, err := login.GetLoginUser(c)
	if err != nil {
		session := sessions.Default(c)
		session.Set("err", "議論にコメントするにはログインが必要です")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/login")
		return
		//errorMsg = err.Error()
	}
	discussId := c.Query("id")
	id, err := strconv.Atoi(discussId)
	if err != nil {
		log.Println("cannot parse ID :", discussId)
		c.Error(err)
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"headerUser": header.GetHeaderUser(usr),
			//ログイン状態
			"loginState": loginState,
		})
		c.Abort()
		return
	}
	dis, err := discuss.GetDiscuss(id)
	log.Println(dis.Decree.Name)
	title := c.PostForm("title")
	content := c.PostForm("content")
	var j discuss.ContentJSON
	err = json.Unmarshal([]byte(dis.Content), &j)
	j = *j.UpdateContent(title, usr.Id, content)
	contentJSON, _ := json.Marshal(j)
	dis.Content = string(contentJSON)
	log.Println(dis.Decree.Name)
	database.DB.Preload(clause.Associations).Updates(&dis)
	go notify.SendNotify(dis)
	c.Redirect(http.StatusSeeOther, "discuss?id="+strconv.Itoa(id))
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
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"headerUser": header.GetHeaderUser(usr),
			//ログイン状態
			"loginState": loginState,
			"errorMsg":   errorMsg,
		})
		c.Abort()
		return
	}
	//議論取得
	dis, err := discuss.GetDiscuss(id)
	if err != nil || dis.Id == 0 {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"headerUser": header.GetHeaderUser(usr),
			//ログイン状態
			"loginState": loginState,
			"errorMsg":   errorMsg,
		})
		c.Abort()
		return
	}
	log.Println(dis.Decree.Id)
	decname := dis.Decree.Name
	disType := dis.Discuss_Type
	disTitle := dis.Title
	//作成ユーザー
	var u user.User
	er := database.DB.Find(&u, dis.Create_User, "id = ?").Error //Select * from user where user.id = Create_User
	if er != nil {
		u = user.User{
			Id:     -1,
			Name:   "退会ユーザー",
			UserID: "退会ユーザー",
		}
	}
	var contentJSON discuss.ContentJSON
	e := json.Unmarshal([]byte(dis.Content), &contentJSON)
	if e != nil {
	}
	var content Content
	for _, s := range contentJSON {
		u, err := user.GetUserByID(s.Create_User)
		if err != nil {
			u = user.User{
				Id:     -1,
				Name:   "退会ユーザー",
				UserID: "退会ユーザー",
			}
		}
		var mentionUser []user.User
		for _, j := range s.MentionTo {
			m, err := user.GetUserByID(j)
			if err != nil {
				mentionUser = append(mentionUser, m)
			}
		}
		content = append(content, Content{{
			//Title:      s.Title,
			Hash:       s.Hash,
			CreateUser: u,
			IsMine:     loginState && u.Id == usr.Id,
			Body:       s.Body,
			MentionTo:  mentionUser,
		}}[0])
	}
	c.HTML(200, "discuss.html", gin.H{
		"headerUser": header.GetHeaderUser(usr),
		//ログイン状態
		"loginState": loginState,
		"errorMsg":   errorMsg,
		"id":         id,
		"decid":      dis.Decree.Id,
		//議論している法令名
		"decree": decname,
		//議論のタイトル
		"discuss": disTitle,
		//議論のカテゴリ(疑問、議論、相談)
		"discussType": disType,
		"opened":      dis.Opened,
		//作成ユーザー
		"createUser": u,
		"IsMine":     u.Id == usr.Id,
		"content":    content,
	})
}

func CreateDiscuss(c *gin.Context) {
	//errorMsg := ""
	usr, loginState, err := login.GetLoginUser(c)
	if err != nil {
		session := sessions.Default(c)
		session.Set("err", "議論作成にはログインが必要です")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/login")
		return
		//errorMsg = err.Error()
	}
	discussType, err := strconv.Atoi(c.PostForm("type"))
	content := c.PostForm("content")
	atoi, err := strconv.Atoi(c.Query("decid"))
	if err != nil {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"headerUser": header.GetHeaderUser(usr),
			//ログイン状態
			"loginState": loginState,
		})
		c.Abort()
		return
	}
	dec, err := decree.GetDecree(atoi)
	var d = discuss.Discuss{
		DecreeId:     dec.Id,
		Decree:       dec,
		Create_User:  usr.Id,
		Discuss_Type: discussType,
		Title:        c.PostForm("title"),
		Opened:       1,
		//Content: ,
	}
	j := discuss.ContentJSON{}
	j = *j.UpdateContent(d.Title, d.Create_User, content)
	contentJSON, _ := json.Marshal(j)
	log.Println(j)
	d.Content = string(contentJSON)

	err = discuss.CreateDiscuss(&d)
	if err != nil {
		return
	}
	c.Redirect(http.StatusSeeOther, "discuss?id="+strconv.Itoa(d.Id))

}

func GetDiscussList(c *gin.Context) {

	errorMsg := ""
	usr, loginState, err := login.GetLoginUser(c)
	if err != nil {
		errorMsg = err.Error()
	}
	sortQuery := c.Query("sort")
	idstr := c.Query("decid")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.Error(err)
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"headerUser": header.GetHeaderUser(usr),
			//ログイン状態
			"loginState": loginState,
			"errorMsg":   errorMsg,
		})
		c.Abort()
		return
	}
	getDecree, err := decree.GetDecree(id)
	if err != nil {
		return
	}
	if getDecree.Id == 0 {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"headerUser": header.GetHeaderUser(usr),
			//ログイン状態
			"loginState": loginState,
			"errorMsg":   errorMsg,
		})
		c.Abort()
		return
	}
	fromDecree, err := discuss.GetDiscussFromDecree(getDecree)
	if sortQuery == "" {
		f := fromDecree
		ff := make([]discuss.Discuss, len(f))
		for i := len(f) - 1; i >= 0; i-- {
			ff[len(f)-i-1] = f[i]
		}
		fromDecree = ff
	}
	var discussList DiscussList
	for _, d := range fromDecree {
		u, err := user.GetUserByID(d.Create_User)
		if err != nil {
			u = user.User{
				Id:     -1,
				Name:   "退会ユーザー",
				UserID: "退会ユーザー",
			}
		}
		var c discuss.ContentJSON
		err = json.Unmarshal([]byte(d.Content), &c)
		if err != nil {
			c = discuss.ContentJSON{}
		}
		discussList = append(discussList, DiscussList{{
			Id:              d.Id,
			Decree:          d.Decree,
			Title:           d.Title,
			CreateUser:      u,
			DiscussType:     d.Discuss_Type,
			Opened:          d.Opened,
			CommentQuantity: len(c),
		}}[0])
	}
	if sortQuery == "comment" {
		sort.Slice(discussList, func(i, j int) bool {
			return discussList[i].CommentQuantity > discussList[j].CommentQuantity
		})
	} else if sortQuery == "commentRev" {
		sort.Slice(discussList, func(i, j int) bool {
			return discussList[i].CommentQuantity < discussList[j].CommentQuantity
		})
	}
	log.Println(len(fromDecree))
	if err != nil {
		c.Error(err)
		return
	}
	law, err := getDecree.GetLaw()
	if err != nil {
		return
	}
	c.HTML(200, "discussList.html", gin.H{
		"headerUser": header.GetHeaderUser(usr),
		//ログイン状態
		"loginState": loginState,
		"errorMsg":   errorMsg,
		"LawNum":     law.LawNum,
		"title":      getDecree.Name,
		"id":         getDecree.Id,
		"sort":       sortQuery,
		"list":       discussList,
	})
}

func CloseDiscuss(c *gin.Context) {
	usr, loginState, err := login.GetLoginUser(c)
	if err != nil || !loginState {
		session := sessions.Default(c)
		session.Set("err", "議論クローズにはログインが必要です")
		session.Save()
		c.Abort()
		c.Redirect(http.StatusSeeOther, "/login")
		return
		//errorMsg = err.Error()
	}

	idstr := c.PostForm("discussid")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		log.Println("parseerror")
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"headerUser": header.GetHeaderUser(usr),
			//ログイン状態
			"loginState": loginState,
		})
		return
	}
	d, err := discuss.GetDiscuss(id)
	if err != nil {
		log.Println("discuss error")
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"headerUser": header.GetHeaderUser(usr),
			//ログイン状態
			"loginState": loginState,
		})
		return
	}
	if d.Create_User == usr.Id {
		d.Opened = 0
		err = database.DB.Save(&d).Error
		if err != nil {
			c.Set("err", "クローズに失敗しました")
			log.Println(err.Error())
		}
		c.Redirect(http.StatusSeeOther, "discuss?id="+strconv.Itoa(id))
	}
}
