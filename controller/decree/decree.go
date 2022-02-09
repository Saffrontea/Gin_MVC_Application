package decree

import (
	"Gin_MVC/model/law"
	"encoding/xml"
	"log"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

func DisplayDecree(c *gin.Context) {
	var l law.Law
	var open *os.File
	var err error
	defer func(open *os.File) {
		err := open.Close()
		if err != nil {
			open.Close()
		}
	}(open)
	str := c.Query("decree")
	if !strings.Contains(str, ".") && len(str) > 0 {
		open, err = os.Open(path.Join("resource", "decree", str, str+".xml"))

		if err != nil {
			c.Error(err)
		}
		log.Println(open.Name())
	} else {
		open, err = os.Open("model/decree/337M50000800043_20220401_504M60000800002.xml")
	}
	//open, err := os.Open("model/decree/417CO0000000018_20220101_503CO0000000186.xml")
	//open, err := os.Open("model/decree/414CO0000000407_20230401_503CO0000000229.xml")
	if err != nil {
		return
	}
	d := xml.NewDecoder(open)
	err = d.Decode(&l)
	c.HTML(202, "decree.html", gin.H{
		"Law": l,
	})
}
