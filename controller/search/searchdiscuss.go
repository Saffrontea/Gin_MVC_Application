package search

import (
	"Gin_MVC/model/database"
	"Gin_MVC/model/decree"
	"Gin_MVC/model/discuss"
	"gorm.io/gorm/clause"
)

type SearchDiscussResult []struct {
	Id          int
	DiscussType int
	Decree      decree.Decree
	Title       string
}

func SearchDisucuss(name string) (*SearchDiscussResult, int) {
	var d []discuss.Discuss
	name = "%" + name + "%"
	database.DB.Preload(clause.Associations).Where("discusses.Title like ?", name).Find(&d)
	if len(d) == 0 {
		return nil, 0
	}
	count := len(d)
	if count > 500 {
		d = d[:500]
	}
	dr := make(SearchDiscussResult, len(d))
	for i := 0; i < len(dr); i++ {
		dr[i].Id = d[i].Id
		dr[i].DiscussType = d[i].Discuss_Type
		dr[i].Title = d[i].Title
		dr[i].Decree = d[i].Decree
	}
	return &dr, count
}

func GetTopDiscusses() *SearchDiscussResult {
	var d []discuss.Discuss
	database.DB.Preload(clause.Associations).Order(clause.OrderByColumn{
		Column: clause.Column{Name: "updated_at"},
		Desc:   true,
	}).Limit(4).Find(&d)
	dr := make(SearchDiscussResult, len(d))
	for i := 0; i < len(dr); i++ {
		dr[i].Id = d[i].Id
		dr[i].DiscussType = d[i].Discuss_Type
		dr[i].Title = d[i].Title
		dr[i].Decree = d[i].Decree
	}
	return &dr
}
