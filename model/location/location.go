package location

import (
	_ "embed"
	"encoding/json"
)

//go:embed LocationList.json
var location []byte

func GetLocationList() []string {
	var locationList []string
	_ = json.Unmarshal(location, &locationList)
	return locationList
}
