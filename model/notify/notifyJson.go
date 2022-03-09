package notify

type NotifyJSON []struct {
	DiscussID int    `json:"discussID"`
	Hash      string `json:"hash"`
	Time      string `json:"time"`
	Level     int    `json:"level"`
	Comment   string `json:"comment"`
}
