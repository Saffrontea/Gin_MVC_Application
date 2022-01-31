package decree

import "encoding/xml"

type RequestXML struct {
	XMLName xml.Name `xml:"RequestXML"`
	Text    string   `xml:",chardata"`
	Result  struct {
		Text    string `xml:",chardata"`
		Code    string `xml:"Code"`
		Message string `xml:"Message"`
	} `xml:"Result"`
	ApplData struct {
		Text            string `xml:",chardata"`
		Date            string `xml:"Date"`
		LawNameListInfo []struct {
			Text                  string `xml:",chardata"`
			LawTypeName           string `xml:"LawTypeName"`
			LawNo                 string `xml:"LawNo"`
			LawName               string `xml:"LawName"`
			LawNameKana           string `xml:"LawNameKana"`
			OldLawName            string `xml:"OldLawName"`
			PromulgationDate      string `xml:"PromulgationDate"`
			AmendName             string `xml:"AmendName"`
			AmendNo               string `xml:"AmendNo"`
			AmendPromulgationDate string `xml:"AmendPromulgationDate"`
			EnforcementDate       string `xml:"EnforcementDate"`
			EnforcementComment    string `xml:"EnforcementComment"`
			LawId                 string `xml:"LawId"`
			LawUrl                string `xml:"LawUrl"`
			EnforcementFlg        string `xml:"EnforcementFlg"`
			AuthFlg               string `xml:"AuthFlg"`
		} `xml:"LawNameListInfo"`
	} `xml:"ApplData"`
}
