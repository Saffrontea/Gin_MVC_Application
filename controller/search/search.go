package search

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func ISearch(c *gin.Context)  {
	str := c.Query("q")
	if strings.HasPrefix(str,"@") {
		
	}
}