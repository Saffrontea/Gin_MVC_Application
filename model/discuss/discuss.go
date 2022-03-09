package discuss

import (
	"Gin_MVC/model/database"
	"Gin_MVC/model/decree"
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
	"log"
	"time"
)

/*
	Discuss
	議論のルート構造体
*/
type Discuss struct {
	Id int `gorm:"primaryKey;autoIncrement"`
	//紐づく法令
	DecreeId int
	Decree   decree.Decree `gorm:"foreignKey:DecreeId"`
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
	Title     string
	Opened    int
	UpdatedAt time.Time
	Content   string `json:"content"`
}

/*
	ContentJSON
	議論の内容
*/
type ContentJSON []struct {
	Title string `json:"title"`
	Hash  string `json:"hash"`
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
	err := database.DB.Preload(clause.Associations).Find(&d, id, "id = ?").Error
	//e := database.DB.Preload(clause.Associations).Find(&d.Decree).Error
	//if e != nil {
	//	err = e
	//}
	return d, err
}

func GetDiscussFromDecree(dec decree.Decree) ([]Discuss, error) {
	var d []Discuss
	err := database.DB.Preload(clause.Associations).Where("decree_id = ?", dec.Id).Find(&d).Error
	return d, err
}

func CreateDiscuss(discuss *Discuss) error {
	err := database.DB.Preload(clause.Associations).Create(&discuss).Error
	database.DB.Find(&discuss)
	return err
}

func CreateDiscussJSON(content ContentJSON) (*string, error) {
	s, err := json.Marshal(content)
	str := string(s)
	if err != nil {
		return nil, err
	}
	return &str, nil
}

func (content *ContentJSON) UpdateContent(title string, userId int, Text string) *ContentJSON {
	c := ContentJSON{
		{
			Title:       title,
			Hash:        uuid.New().String(),
			Create_User: userId,
			Body:        Text,
			MentionTo:   []int{},
		},
	}

	if len(*content) > 0 {
		flg := false
		c[0].MentionTo = (*content)[len(*content)-1].MentionTo
		for _, i := range (*content)[len(*content)-1].MentionTo {
			flg = (*content)[len(*content)-1].Create_User == i
		}
		if !flg {
			c[0].MentionTo = append(c[0].MentionTo, (*content)[len(*content)-1].Create_User)
		} else {
			c[0].MentionTo = (*content)[len(*content)-1].MentionTo
		}
	}
	log.Println(c[0].MentionTo)
	//*content[len(content)-1]
	c = append(*content, c...)
	content = &c
	return content
}

func CloseOldDiscuss() {
	now := time.Now()
	time.Sleep(time.Until(time.Date(now.Year(), now.Month(), now.AddDate(0, 0, 1).Day(), 0, 0, 0, 0, now.Location())))
	for {
		var d []Discuss
		database.DB.Where("discusses.updated_at < ?", time.Now().AddDate(0, 0, -30)).Find(&d)
		log.Println(d)
		for _, dis := range d {
			dis.Opened = -1
		}
		database.DB.Updates(d)
		time.Sleep(time.Hour * 24)
	}
}
