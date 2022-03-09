package search

import (
	"Gin_MVC/controller/header"
	"Gin_MVC/controller/login"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Search(c *gin.Context) {
	errorMsg := ""
	usr, loginState, err := login.GetLoginUser(c)
	if err != nil {
		errorMsg = err.Error()
	}

	str := c.Query("q")
	var u *SearchUserResult
	var userCount int
	var d *SearchDecreeResult
	var decreeCount int
	var dis *SearchDiscussResult
	var discussCount int
	if strings.HasPrefix(str, "@") {
		u, userCount = SearchUser(str[1:])
	} else if strings.HasPrefix(str, "") { //TODO:?
		u, userCount = SearchUser(str)
		d, decreeCount = SearchDecree(str)
		dis, discussCount = SearchDisucuss(str)
	}
	c.HTML(http.StatusOK, "search.html", gin.H{
		"headerUser": header.GetHeaderUser(usr),
		//ログイン状態
		"loginState":       loginState,
		"errorMsg":         errorMsg,
		"query":            str,
		"SearchUserRes":    u,
		"SearchDecreeRes":  d,
		"userCount":        userCount,
		"DiscussCount":     discussCount,
		"DecreeCount":      decreeCount,
		"SearchDiscussRes": dis,
	})

}
