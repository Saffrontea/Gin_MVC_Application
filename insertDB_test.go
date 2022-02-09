package main

import (
	"Gin_MVC/model/database"
	"Gin_MVC/model/decree"
	"Gin_MVC/model/discuss"
	"Gin_MVC/model/notify"
	"Gin_MVC/model/priority"
	"Gin_MVC/model/user"
	"os"
	"path"
	"testing"
	"time"
)

func TestInsertDB(t *testing.T) {
	_ = database.DBConnection()
	database.Migrator([]interface{}{&user.User{}, &decree.Decree{}, &notify.Notify{}, &priority.Priority{}, &discuss.Discuss{}})
	dir, err := os.ReadDir("resource/decree")
	if err != nil {
		panic(err)
	}
	for _, entry := range dir {
		open, err := os.Open(path.Join("resource", "decree", entry.Name(), entry.Name()+".xml"))
		if err != nil {
			// panic(err)
			return
		}
		law, err := decree.CreateLaw(open)
		if err != nil {
			panic(err)
		}
		t.Log(law.LawBody.LawTitle.Text)
		err = decree.CreateDecree(
			decree.Decree{
				DecreeReference: entry.Name(),
				Name:            law.LawBody.LawTitle.Text,
				LastUpdate:      time.Date(2022, 1, 32, 0, 0, 0, 0, time.Local),
			},
		)
		if err != nil {
			panic(err)
		}
	}

}
