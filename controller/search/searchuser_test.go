package search_test

import (
	"Gin_MVC/controller/search"
	"Gin_MVC/model/database"
	"Gin_MVC/model/priority"
	"Gin_MVC/model/user"
	"testing"
)

func TestSearchUser(t *testing.T) {
	database.DBConnection()
	database.Migrator([]interface{}{&user.User{},&priority.Priority{}})
	name := "sa"
	t.Log(*search.SearchUser(name))
}