/*
	index.go
	indexページ処理
*/
package index

import (
	"Gin_MVC/controller/header"
	"Gin_MVC/controller/login"
	"github.com/gin-gonic/gin"
)

func Display(c *gin.Context) {
	errorMsg := ""

	usr, loginState, err := login.GetLoginUser(c)
	if err != nil {
		errorMsg = err.Error()
	}

	c.HTML(200, "index.html", gin.H{
		//ヘッダーのユーザー情報
		"headerUser": header.GetHeaderUser(usr),

		//ログイン状態
		"loginState": loginState,
		"errorMsg":   errorMsg,
	})
}
