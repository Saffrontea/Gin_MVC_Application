package main

import (
	"Gin_MVC/model/database"
	"Gin_MVC/model/decree"
	"Gin_MVC/model/discuss"
	"Gin_MVC/model/notify"
	"Gin_MVC/model/priority"
	"Gin_MVC/model/user"
	"log"
	"testing"
)

func TestDB(t *testing.T) {
	er := database.DBConnection()
	if er != nil {
		log.Fatal(er)
	}
	database.Migrator([]interface{}{&user.User{},&decree.Decree{}, &notify.Notify{}, &priority.Priority{},&discuss.Discuss{}})
	d := []decree.Decree{
		{
		Id:               1,
		Decree_Reference: "",
		Name:             "aaa",
		Last_update:      nil,
		},
		{
			Id: 2,
			Name: "bbb",
		},
		{
			Id:3,
			Name: "ccc",
		},
	}
	database.DB.Create(&d)
}