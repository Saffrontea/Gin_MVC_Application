package decree

import (
	"Gin_MVC/controller/header"
	"Gin_MVC/controller/login"
	"Gin_MVC/model/database"
	"Gin_MVC/model/decree"
	"Gin_MVC/model/law"
	"Gin_MVC/model/user"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5/plumbing"
	"gorm.io/gorm/clause"
	"log"
	"net/http"
	"strconv"
)

func Display(c *gin.Context) {
	errorMsg := ""
	usr, loginState, err := login.GetLoginUser(c)
	if err != nil {
		errorMsg = err.Error()
	}

	decreeId := c.Query("decid")
	rev := c.Query("rev")
	id, err := strconv.Atoi(decreeId)
	if err != nil {
		log.Println("cannot parse ID :", decreeId)
		c.Error(err)
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"headerUser": header.GetHeaderUser(usr),
			//ログイン状態
			"loginState": loginState,
			"errorMsg":   errorMsg,
		})
	}
	dec, err := decree.GetDecree(id)
	starflg := false
	if loginState {
		for _, star := range usr.Star {
			starflg = starflg || star.Star == id
		}
	}
	if err != nil {
		c.Error(err)
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"headerUser": header.GetHeaderUser(usr),
			//ログイン状態
			"loginState": loginState,
			"errorMsg":   errorMsg,
		})
	}
	decname := dec.Name
	declastupdate := dec.LastUpdate
	var l *law.Law
	if rev == "" {
		l, _ = dec.GetLaw()
	} else {
		l, _ = dec.GetOldLaw(plumbing.NewHash(rev))
	}
	n, _ := dec.GetRevisions()
	msg, _ := c.Get("msg")
	c.Set("msg", "")
	c.HTML(200, "decree.html", gin.H{
		"headerUser": header.GetHeaderUser(usr),
		//ログイン状態
		"loginState": loginState,
		"errorMsg":   errorMsg,
		"id":         id,
		"Law":        l,
		"msg":        msg,
		"star":       starflg,
		"revisions":  n,
		"rev":        rev,
		"decid":      dec.Id,
		//法令名
		"decreeName": decname,
		//Last Commit Time
		"decreeLastUpdate": declastupdate,
	})
}

func RemoveStar(c *gin.Context) {
	errorMsg := ""
	usr, loginState, err := login.GetLoginUser(c)
	if err != nil {
		errorMsg = "お気に入り削除にはログインが必要です"
		s := sessions.Default(c)
		s.Set("err", errorMsg)
		s.Save()
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	decreeId := c.Query("decid")
	id, err := strconv.Atoi(decreeId)
	if err != nil {
		log.Println("cannot parse ID :", decreeId)
		c.Error(err)
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"headerUser": header.GetHeaderUser(usr),
			//ログイン状態
			"loginState": loginState,
			"errorMsg":   errorMsg,
		})
		return
	}
	getDecree, err := decree.GetDecree(id)
	if err != nil {
		c.AbortWithStatus(502)
		return
	}
	usr.RemoveStar(getDecree.Id)
	log.Println(usr.Star)
	database.DB.Save(usr)
	c.Redirect(http.StatusSeeOther, "/decree?decid="+decreeId)

}

func AddStar(c *gin.Context) {
	errorMsg := ""
	usr, loginState, err := login.GetLoginUser(c)
	if err != nil {
		errorMsg = "お気に入り登録にはログインが必要です"
		s := sessions.Default(c)
		s.Set("err", errorMsg)
		s.Save()
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	decreeId := c.Query("decid")
	id, err := strconv.Atoi(decreeId)
	if err != nil {
		log.Println("cannot parse ID :", decreeId)
		c.Error(err)
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"headerUser": header.GetHeaderUser(usr),
			//ログイン状態
			"loginState": loginState,
			"errorMsg":   errorMsg,
		})
		return
	}
	getDecree, err := decree.GetDecree(id)
	if err != nil {
		c.AbortWithStatus(502)
		return
	}
	flg := false
	for _, star := range usr.Star {
		flg = flg || star.Star == getDecree.Id
	}
	if !flg {
		usr.Star = append(usr.Star, user.Star{
			UserId: usr.Id,
			Star:   getDecree.Id,
		})
		database.DB.Save(usr)
		c.Set("msg", "お気に入り登録が完了しました")
		c.Redirect(http.StatusSeeOther, "decree?decid="+strconv.Itoa(id))
	} else {
		c.Set("msg", "既に登録済みです")
		c.Redirect(http.StatusSeeOther, "decree?decid="+strconv.Itoa(id))
	}

}

func DisplayDecreeList(c *gin.Context) {
	var d []decree.Decree
	var o clause.OrderByColumn
	if c.Query("sort") == "true" {
		o = clause.OrderByColumn{Column: clause.Column{Name: "id"}, Desc: false}
	} else {
		o = clause.OrderByColumn{Column: clause.Column{Name: "id"}, Desc: true}
	}
	database.DB.Model(&decree.Decree{}).Order(o).Where("name like ?", "%"+c.Query("search")+"%").Find(&d)
	c.HTML(200, "decreeList.html", gin.H{
		"decreeList": d,
	})
}
