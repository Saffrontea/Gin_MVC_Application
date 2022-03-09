package search

import (
	"Gin_MVC/model/database"
	"Gin_MVC/model/decree"
	"gorm.io/gorm/clause"
)

type SearchDecreeResult []struct {
	Id   int
	Name string
}

func SearchDecree(name string) (*SearchDecreeResult, int) {
	var d []decree.Decree
	name = "%" + name + "%"

	database.DB.Preload(clause.Associations).Where("decrees.Name like ?", name).Order("id desc").Find(&d)
	if len(d) == 0 {
		return nil, 0
	}
	count := len(d)
	if count > 500 {
		d = d[:500]
	}
	dr := make(SearchDecreeResult, len(d))
	for i := 0; i < len(dr); i++ {
		dr[i].Id = d[i].Id
		dr[i].Name = d[i].Name
	}
	return &dr, count
}

//func SearchDecreeAPI(c *gin.Context) {
//	name := c.Query("q")
//	c.JSON(200, *SearchDecree(name))
//}

func GetTopDecrees() *SearchDecreeResult {
	var d []decree.Decree

	database.DB.Preload(clause.Associations).Order("last_update desc").Order("id desc").Limit(4).Find(&d)
	if len(d) == 0 {
		return nil
	}
	dr := make(SearchDecreeResult, len(d))
	for i := 0; i < len(dr); i++ {
		dr[i].Id = d[i].Id
		dr[i].Name = d[i].Name
	}
	return &dr
}
