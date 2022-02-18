/*
	decree
	法令
*/
package decree

import (
	"Gin_MVC/model/database"
	"Gin_MVC/model/law"
	"encoding/xml"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const DecreeDir = "resource/decree"

type Decree struct {
	Id int `gorm:"primaryKey;autoIncrement"`
	//gitディレクトリのPath
	DecreeReference string
	//法令名
	Name string
	//Last Commit Time
	LastUpdate time.Time
}

/*
	法令取得
*/
func GetDecree(id int) (Decree, error) {
	decree := Decree{}
	err := database.DB.Find(&decree, "Id", id).Error
	return decree, err
}

/*
	法令作成

*/
func CreateDecree(decree Decree) error {
	//TODO: バッチ処理でこれを呼び出す
	return database.DB.Create(&decree).Error
}

func UpdateDecree(decree Decree) error {
	return database.DB.Model(&Decree{}).Where("decree_reference = ?", decree.DecreeReference).Update("last_update", decree.LastUpdate).Error
}
func GetTodoyUpdate() error {
	get, err := http.Get("https://elaws.e-gov.go.jp/api/1/updatelawlists/" + time.Now().Format("20060102"))
	if err != nil {
		return err
	}
	r := xml.NewDecoder(get.Body)
	var x law.RequestXML
	err = r.Decode(&x)
	if err != nil {
		return err
	}
	if x.Result.Code == "0" {
		for _, s := range x.ApplData.LawNameListInfo {
			log.Println(s.LawId)
			GetDecreeFromAPI(s.LawId, x.ApplData.Date)
			time.Sleep(10 * time.Millisecond)
		}
	}
	return nil
}

func GetDecreeFromAPI(Id string, date string) {
	var l law.Law
	var r law.RequestRespXML
	var d Decree
	dpath := path.Join(DecreeDir, Id)
	get, err := http.Get("https://elaws.e-gov.go.jp/api/1/lawdata/" + Id)
	if err != nil {
		return
	}
	log.Println(get.StatusCode)
	err = xml.NewDecoder(get.Body).Decode(&r)
	log.Println(r.ApplData.LawFullText)
	// err = xml.Unmarshal([]byte(r.ApplData.LawFullText.Law), &l)
	l = r.ApplData.LawFullText.Law
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(l.LawBody.LawTitle)
	d = Decree{
		DecreeReference: Id,
		Name:            l.LawBody.LawTitle.Text,
		LastUpdate:      time.Now(),
	}
	_, e := os.Open(dpath)
	ext := false
	if e != nil {
		ext = os.IsExist(e)
	}
	if !ext {
		// dpath = "../../" + dpath
		err = os.Mkdir(dpath, os.ModePerm)
		if err != nil {
			log.Println("mkdir")
			log.Println(err)
			return
		}
	}
	var b []byte
	b, _ = xml.Marshal(l)
	err = os.WriteFile(path.Join(dpath, Id+".xml"), b, os.ModePerm)
	if err != nil {
		log.Println("write")
		log.Println(err)
		return
	}
	if _, e := os.Open(dpath); os.IsExist(e) {
		open, err := git.PlainOpen(dpath)
		if err != nil {
			log.Println("open")
			log.Println(err)
			return
		}
		worktree, err := open.Worktree()
		if err != nil {
			log.Println("worktree")
			log.Println(err)
			return
		}
		worktree.Add(path.Join(dpath, Id+".xml"))
		_, err = worktree.Commit(date, &git.CommitOptions{
			Author: &object.Signature{
				Name:  "GitLaw",
				Email: "",
			},
		})
		if err != nil {
			log.Println("commit")
			log.Println(err)
			return
		}

		if err != nil {
			log.Println(err)
			return
		}
		CreateDecree(d)
	} else {
		init, err := git.PlainInit(dpath, false)
		if err != nil {
			log.Println(err)
			return
		}
		worktree, err := init.Worktree()
		if err != nil {
			log.Println(err)
			return
		}
		_, err = worktree.Add(path.Join(dpath, Id+".xml"))
		if err != nil {
			log.Println(err)
			return
		}
		_, err = worktree.Commit(date, &git.CommitOptions{})
		if err != nil {
			log.Println(err)
			return
		}

		err = UpdateDecree(d)
		if err != nil {
			log.Println(err)
			return
		}
	}

}

func (decree Decree) GetLaw() (*law.Law, error) {
	f, e := os.Open(path.Join(DecreeDir, decree.DecreeReference, decree.DecreeReference+".xml"))
	if e != nil {
		return nil, e
	}
	return law.CreateLaw(f)
}


