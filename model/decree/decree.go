/*
	decree
	法令
*/
package decree

import (
	"Gin_MVC/model/database"
	"Gin_MVC/model/law"
	"encoding/xml"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/go-git/go-git/v5"
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
			GetDecreeFromAPI(s.LawId, x.ApplData.Date)
		}
	}
	return nil
}

func GetDecreeFromAPI(Id string, date string) {
	var l law.Law
	var d Decree
	dpath := path.Join(DecreeDir, Id)
	get, err := http.Get("https://elaws.e-gov.go.jp/api/1/lawdata/" + Id)
	if err != nil {
		return
	}
	err = xml.NewDecoder(get.Body).Decode(&l)
	if err != nil {
		return
	}

	d = Decree{
		DecreeReference: Id,
		Name:            l.LawBody.LawTitle.Line[0],
		LastUpdate:      time.Now(),
	}
	_, e := os.Open(dpath)
	ext := os.IsExist(e)
	if !ext {
		err = os.Mkdir(dpath, os.ModeDir)
		if err != nil {
			return
		}
	}
	var b []byte
	_, err = get.Body.Read(b)
	if err != nil {
		return
	}
	err = os.WriteFile(path.Join(dpath, Id+".xml"), b, 0655)
	if err != nil {
		return
	}
	if _, e := os.Open(dpath); os.IsExist(e) {
		open, err := git.PlainOpen(dpath)
		if err != nil {
			return
		}
		worktree, err := open.Worktree()
		if err != nil {
			return
		}
		_, err = worktree.Commit(date, &git.CommitOptions{})
		if err != nil {
			return
		}

		if err != nil {
			return
		}
		CreateDecree(d)
	} else {
		init, err := git.PlainInit(dpath, false)
		if err != nil {
			return
		}
		worktree, err := init.Worktree()
		if err != nil {
			return
		}
		_, err = worktree.Add(path.Join(dpath, Id+".xml"))
		if err != nil {
			return
		}
		_, err = worktree.Commit(date, &git.CommitOptions{})
		if err != nil {
			return
		}

		err = UpdateDecree(d)
		if err != nil {
			return
		}
	}

}

func (decree Decree) GetDecree() (*law.Law, error) {
	f, e := os.Open(path.Join(DecreeDir, decree.DecreeReference, decree.DecreeReference+".xml"))
	if e != nil {
		return nil, e
	}
	return law.CreateLaw(f)
}
