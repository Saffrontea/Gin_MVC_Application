package user

import (
	"Gin_MVC/model/database"
	"Gin_MVC/model/discuss"
	"Gin_MVC/model/notify"
	"Gin_MVC/model/priority"
	"encoding/json"
	"log"
	"testing"

	"gorm.io/gorm"
)

func TestUser(t *testing.T) {
	er := database.DBConnection()
	if er != nil {
		log.Fatal(er)
	}
	database.Migrator([]interface{}{&User{}, &Star{}, &notify.Notify{}, &priority.Priority{}, &discuss.Discuss{}})
	s, erro := GetUser("saitou")
	log.Print(s, erro)
	ts, _ := json.Marshal(notify.NotifyJSON{
		struct {
			DiscussID int    "json:\"discussID\""
			Hash      string "json:\"hash\""
			Level     int    "json:\"level\""
			Comment   string "json:\"comment\""
		}{},
	})
	var u = User{
		//UserId:       0,
		Name:     "斉藤",
		Ruby:     "サイトウ",
		UserID:   "saito",
		Password: "asdfgdf",
		Tel:      "03000000000",
		Location: 0,
		Publish:  false,
		Notify:   notify.Notify{Notify: string(ts)},
		Priority: priority.Priority{Priority: 100},
		Discuss:  []discuss.Discuss{},
		Star:     []Star{},
	}

	if err := database.Transaction(func(tx *gorm.DB) error {
		if err := CreateUser(&u, tx).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatal(err.Error())
	}

	log.Println(u.Star)
	database.DB.Save(&u)
}
