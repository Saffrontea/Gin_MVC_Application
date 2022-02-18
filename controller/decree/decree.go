package decree

import (
	"Gin_MVC/controller/header"
	"Gin_MVC/controller/login"
	"Gin_MVC/model/decree"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func Display(c *gin.Context){
	errorMsg := ""
	usr, loginState, err := login.GetLoginUser(c)
	if err != nil {
		errorMsg = err.Error()
	}

	decreeId := c.Query("id")
	id, err := strconv.Atoi(decreeId)
	if err != nil {
		log.Println("cannot parse ID :", decreeId)
		c.Error(err)
	}
	dec ,err := decree.GetDecree(id)
	decname := dec.Name
	declastupdate := dec.LastUpdate
	l,_:= dec.GetLaw()

	c.HTML(200,"decree.html" ,gin.H{
		"headerUser": header.GetHeaderUser(usr),
		//ログイン状態
		"loginState": loginState,
		"errorMsg":   errorMsg,

		"Law":        l,
		//法令名
		"decreeName" : decname,
		//Last Commit Time
		"decreeLastUpdate" : declastupdate,
	})

}
