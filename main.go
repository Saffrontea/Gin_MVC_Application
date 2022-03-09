package main

import (
	"Gin_MVC/model/database"
	"Gin_MVC/model/decree"
	"Gin_MVC/model/discuss"
	"Gin_MVC/model/notify"
	"Gin_MVC/model/priority"
	"Gin_MVC/model/user"
	"Gin_MVC/router"
)

//go:generate  air

func main() {
	_ = database.DBConnection()
	user.WaitAuthUsers = map[string]int{}
	user.PasswordForgetUsers = map[string]int{}
	database.Migrator([]interface{}{&user.User{}, &decree.Decree{}, &notify.Notify{}, &priority.Priority{}, &discuss.Discuss{}, &user.Star{}})
	r := router.GetRouter()
	go discuss.CloseOldDiscuss()
	r.Run(":3000")
}
