package decree

import (
	"encoding/xml"
	"log"
	"os"
	"testing"
)

func TestLaw(t *testing.T) {
	var l Law
	open, err := os.Open("417CO0000000018_20220101_503CO0000000186.xml")
	if err != nil {
		return
	}
	d := xml.NewDecoder(open)
	err = d.Decode(&l)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(l.LawBody.MainProvision.Article[0].Paragraph[0].ParagraphSentence)
}
