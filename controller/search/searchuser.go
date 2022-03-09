package search

import (
	"Gin_MVC/model/database"
	"Gin_MVC/model/user"
)

type SearchUserResult []struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Image   string `json:"image"`
	UserID  string `json:"userid"`
	Profile string `json:"profile"`
}

func SearchUser(name string) (*SearchUserResult, int) {
	var u []user.User
	name = "%" + name + "%"
	err := database.DB.Joins("inner join  priorities on users.id = priorities.id").Where("name like ? and name not like '退会ユーザー'", name).Or("user_id like ? and name not like '退会ユーザー'", name).Order("priority").Limit(20).Find(&u).Error
	if err != nil {
		return nil, 0
	}
	if len(u) == 0 {
		return nil, 0
	}
	count := len(u)
	if count > 500 {
		u = u[:500]
	}
	su := make(SearchUserResult, len(u))
	for i := 0; i < len(su); i++ {
		su[i].Id = u[i].Id
		su[i].Name = u[i].GetUserName()
		su[i].Image = u[i].Image.GetImage()
		su[i].UserID = "@" + u[i].UserID
		su[i].Profile = u[i].Profile
	}
	return &su, count
}

//func SearchUserAPI(c *gin.Context) {
//	name := c.Query("q")
//	c.JSON(200, *SearchUser(name))
//}
