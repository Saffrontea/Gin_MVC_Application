package main

import (
	"Gin_MVC/model/database"
	"Gin_MVC/model/decree"
	"Gin_MVC/model/discuss"
	"Gin_MVC/model/notify"
	"Gin_MVC/model/priority"
	"Gin_MVC/model/user"
	"github.com/google/uuid"
	"log"
	"testing"
)

func TestDB(t *testing.T) {
	er := database.DBConnection()
	if er != nil {
		log.Fatal(er)
	}
	database.Migrator([]interface{}{&user.User{}, &decree.Decree{}, &notify.Notify{}, &priority.Priority{}, &discuss.Discuss{}})
	d := []decree.Decree{
		{
			Id:              1,
			DecreeReference: "",
			Name:            "aaa",
			//LastUpdate:      nil,
		},
		{
			Id:   2,
			Name: "bbb",
		},
		{
			Id:   3,
			Name: "ccc",
		},
	}
	database.DB.Create(&d)
}

func TestAddnotify(t *testing.T) {
	er := database.DBConnection()
	if er != nil {
		log.Fatal(er)
	}
	database.Migrator([]interface{}{&user.User{}, &decree.Decree{}, &notify.Notify{}, &priority.Priority{}, &discuss.Discuss{}})
	usr, err := user.GetUserByID(4)
	if err != nil {
		return
	}
	usr.Notify.Id = 4
	usr.Notify.AddNotify(notify.NotifyJSON{
		{
			DiscussID: 0,
			Hash:      uuid.NewString(),
			Level:     1,
			Comment:   "Aさんがコメントしました",
		},
	})
	database.DB.Save(&usr)
}
