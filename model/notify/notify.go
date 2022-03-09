package notify

import (
	"Gin_MVC/model/database"
	"encoding/json"
	"gorm.io/gorm"
	"log"
)

type Notify struct {
	Id     int    `gorm:"primaryKey"`
	Notify string `json:"notify"`
}

func CreateNotify(user int, tx *gorm.DB) *gorm.DB {
	var s string
	t, _ := json.Marshal(NotifyJSON{
		{},
	})
	s = string(t)
	return tx.Create(&Notify{
		Id:     user,
		Notify: s,
	})
}

func GetNotify(user int) (Notify, error) {
	var n Notify
	e := database.DB.First(&n, user).Error
	return n, e
}

func (n Notify) AddNotify(comment NotifyJSON) Notify {
	var j NotifyJSON
	err := json.Unmarshal([]byte(n.Notify), &j)
	log.Println(n.Notify)
	if err != nil {
		log.Println(err)
		j = NotifyJSON{}
	}
	j = append(j, comment...)
	b, _ := json.Marshal(j)
	n.Notify = string(b)
	log.Println(n.Notify)
	var r = Notify{}
	//tx := database.DB.First(&r, n.Id).Update("Notify", n.Notify)
	//tx := database.DB.Save(&n)
	database.DB.Find(&n, "id", r.Id)
	return n //, tx.Error
}

func (n Notify) RemoveNotify() error {
	var r = Notify{}
	b, _ := json.Marshal(NotifyJSON{})
	n.Notify = string(b)

	return database.DB.First(&r, n.Id).Update("Notify", n.Notify).Error
}
