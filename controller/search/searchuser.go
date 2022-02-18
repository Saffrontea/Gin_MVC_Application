package search

import (
	"Gin_MVC/model/database"
	"Gin_MVC/model/user"
	"github.com/gin-gonic/gin"
)

type SearchUserResult []struct {
	Name string `json:"name"`
	Image string `json:"image"`
	Username string `json:"username"`
}

func SearchUser(name string) *SearchUserResult{
	var u []user.User
	name = "%"+name+"%"
	err := database.DB.Joins("inner join  priorities on users.id = priorities.id").Where("name like ? and publish = true" ,name).Or("username like ?",name).Order("priority").Limit(20).Find(&u).Error
	if err != nil {
		return &SearchUserResult{}
	}
	su := make(SearchUserResult,len(u))
	for i := 0; i < len(su); i++ {
		su[i].Name = u[i].GetUserName()
		su[i].Image = u[i].Image.GetImage()
		su[i].Username = "@"+u[i].Username
	}
	return &su
}
func SearchUserAPI(c *gin.Context)  {
	name := c.Query("q")
	c.JSON(200,*SearchUser(name))
}
