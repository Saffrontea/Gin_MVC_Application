package news

import (
	"Gin_MVC/controller/header"
	"Gin_MVC/controller/login"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func Display(c *gin.Context) {
	errorMsg := ""
	usr, loginState, err := login.GetLoginUser(c)
	if err != nil {
		errorMsg = err.Error()
	}
	str := c.Query("q")
	if str == "" {
		c.HTML(http.StatusOK, "newsList.html", gin.H{
			"headerUser": header.GetHeaderUser(usr),
			//ログイン状態
			"loginState": loginState,
			"errorMsg":   errorMsg,
		})
		return
	} else {
		_, err = os.Open("view/news/news" + str + ".html")
	}

	if os.IsNotExist(err) {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"headerUser": header.GetHeaderUser(usr),
			//ログイン状態
			"loginState": loginState,
			"errorMsg":   errorMsg,
		})
		return
	} else {
		c.HTML(http.StatusOK, "news"+str+".html", gin.H{
			"headerUser": header.GetHeaderUser(usr),
			//ログイン状態
			"loginState": loginState,
			"errorMsg":   errorMsg,
		})
	}

}
