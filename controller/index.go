package controller

import (
	"Gin_MVC/model/user"
	"github.com/gin-gonic/gin"
)

func IndexDisplayAction(c *gin.Context) {
	user, err := user.GetUser("saitou")

	if err != nil {

	}
	//userprofile ,er :=

	c.HTML(200, "index.html", gin.H{

		"username": user.Name,
		//"userprofile": userprofile,
		"str": "Index Page",
	})
}
