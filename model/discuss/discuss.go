package discuss

import (
	"Gin_MVC/model/database"
	"Gin_MVC/model/decree"
	"encoding/json"
	"hash"

	"gorm.io/gorm/clause"
)

/*
	Discuss
	議論のルート構造体
*/
type Discuss struct {
	Id int `gorm:"primaryKey;autoIncrement"`
	//紐づく法令
	Decree decree.Decree `gorm:"primaryKey;foreignKey:Id"`
	//作成ユーザーのId
	Create_User int
	/*
		議論の種類
		0 : 疑問
		1 : 議論
		2 : 相談
	*/
	Discuss_Type int
	//タイトル
	Title   string
	Opened  int
	Content json.RawMessage `json:"content"`
}

/*
	ContentJSON
	議論の内容
*/
type ContentJSON []struct {
	Title string    `json:"title"`
	Hash  hash.Hash `json:"hash"`
	//作成ユーザーのID
	Create_User int    `json:"createUser"`
	Body        string `json:"body"`
	//メンションユーザーのID
	MentionTo []int `json:"mentionTo"`
}

/*
	GetDiscuss
	議論を取得する
*/
func GetDiscuss(id int) (Discuss, error) {
	var d Discuss
	err := database.DB.Find(&d, id, "id = ?").Error
	e := database.DB.Preload(clause.Associations).Find(&decree.Decree{}).Error
	if e != nil {
		err = e
	}
	return d, err
}

// func CreateDiscuss(discuss Discuss) error{
// 	var s string
// 	_ = json.NewEncoder().Encode(Discuss{})

// 	return database.DB.Create(&Discuss{
// 		Content: s,
// 	}).Error
// }
