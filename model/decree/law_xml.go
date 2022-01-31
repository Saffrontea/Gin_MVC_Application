package decree

import (
	"encoding/xml"
	"os"
)

type Law struct {
	XMLName         xml.Name `xml:"Law"`
	Text            string   `xml:",chardata"`
	Era             string   `xml:"Era,attr"`
	Lang            string   `xml:"Lang,attr"`
	LawType         string   `xml:"LawType,attr"`
	Num             string   `xml:"Num,attr"`
	Year            string   `xml:"Year,attr"`
	PromulgateMonth string   `xml:"PromulgateMonth,attr"`
	PromulgateDay   string   `xml:"PromulgateDay,attr"`
	LawNum          string   `xml:"LawNum"`
	LawBody         struct {
		Text     string `xml:",chardata"`
		LawTitle struct {
			Text       string `xml:",chardata"`
			Kana       string `xml:"Kana,attr"`
			Abbrev     string `xml:"Abbrev,attr"`
			AbbrevKana string `xml:"AbbrevKana,attr"`
		} `xml:"LawTitle"`
		EnactStatement string `xml:"EnactStatement"`
		MainProvision  struct {
			Text    string `xml:",chardata"`
			Article []struct {
				Text           string `xml:",chardata"`
				Num            string `xml:"Num,attr"`
				ArticleCaption string `xml:"ArticleCaption"`
				ArticleTitle   string `xml:"ArticleTitle"`
				Paragraph      []struct {
					Text              string `xml:",chardata"`
					Num               string `xml:"Num,attr"`
					ParagraphNum      string `xml:"ParagraphNum"`
					ParagraphSentence struct {
						Text     string `xml:",chardata"`
						Sentence []struct {
							Text     string `xml:",chardata"`
							Num      string `xml:"Num,attr"`
							Function string `xml:"Function,attr"`
						} `xml:"Sentence"`
					} `xml:"ParagraphSentence"`
					Item []struct {
						Text         string `xml:",chardata"`
						Num          string `xml:"Num,attr"`
						ItemTitle    string `xml:"ItemTitle"`
						ItemSentence struct {
							Text     string `xml:",chardata"`
							Sentence string `xml:"Sentence"`
							Column   []struct {
								Text     string `xml:",chardata"`
								Num      string `xml:"Num,attr"`
								Sentence string `xml:"Sentence"`
							} `xml:"Column"`
						} `xml:"ItemSentence"`
						Subitem1 []struct {
							Text             string `xml:",chardata"`
							Num              string `xml:"Num,attr"`
							Subitem1Title    string `xml:"Subitem1Title"`
							Subitem1Sentence struct {
								Text     string `xml:",chardata"`
								Sentence string `xml:"Sentence"`
								Column   []struct {
									Text     string `xml:",chardata"`
									Num      string `xml:"Num,attr"`
									Sentence string `xml:"Sentence"`
								} `xml:"Column"`
							} `xml:"Subitem1Sentence"`
						} `xml:"Subitem1"`
					} `xml:"Item"`
				} `xml:"Paragraph"`
			} `xml:"Article"`
		} `xml:"MainProvision"`
		SupplProvision []struct {
			Text                string `xml:",chardata"`
			Extract             string `xml:"Extract,attr"`
			AmendLawNum         string `xml:"AmendLawNum,attr"`
			Type                string `xml:"Type,attr"`
			SupplProvisionLabel string `xml:"SupplProvisionLabel"`
			Article             []struct {
				Text           string `xml:",chardata"`
				Num            string `xml:"Num,attr"`
				ArticleCaption string `xml:"ArticleCaption"`
				ArticleTitle   string `xml:"ArticleTitle"`
				Paragraph      []struct {
					Text              string `xml:",chardata"`
					Num               string `xml:"Num,attr"`
					ParagraphNum      string `xml:"ParagraphNum"`
					ParagraphSentence struct {
						Text     string `xml:",chardata"`
						Sentence string `xml:"Sentence"`
					} `xml:"ParagraphSentence"`
					Item []struct {
						Text         string `xml:",chardata"`
						Num          string `xml:"Num,attr"`
						ItemTitle    string `xml:"ItemTitle"`
						ItemSentence struct {
							Text     string `xml:",chardata"`
							Sentence string `xml:"Sentence"`
						} `xml:"ItemSentence"`
					} `xml:"Item"`
					TableStruct struct {
						Text  string `xml:",chardata"`
						Table struct {
							Text     string `xml:",chardata"`
							TableRow []struct {
								Text        string `xml:",chardata"`
								TableColumn []struct {
									Text     string `xml:",chardata"`
									Rowspan  string `xml:"rowspan,attr"`
									Sentence string `xml:"Sentence"`
								} `xml:"TableColumn"`
							} `xml:"TableRow"`
						} `xml:"Table"`
					} `xml:"TableStruct"`
				} `xml:"Paragraph"`
			} `xml:"Article"`
			Paragraph []struct {
				Text              string `xml:",chardata"`
				Num               string `xml:"Num,attr"`
				ParagraphNum      string `xml:"ParagraphNum"`
				ParagraphSentence struct {
					Text     string `xml:",chardata"`
					Sentence []struct {
						Text     string `xml:",chardata"`
						Function string `xml:"Function,attr"`
						Num      string `xml:"Num,attr"`
					} `xml:"Sentence"`
				} `xml:"ParagraphSentence"`
				ParagraphCaption string `xml:"ParagraphCaption"`
			} `xml:"Paragraph"`
		} `xml:"SupplProvision"`
	} `xml:"LawBody"`
}

func CreateLaw(file *os.File) (*Law, error) {
	var l Law
	var r []byte
	file.Read(r)
	err := xml.Unmarshal(r, &l)
	if err != nil {
		return nil, err
	}
	return &l, nil
}
