package common

var WeekdayString = [7]string{
	"Sunday",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
}

// Meet 授業の情報を格納する構造体
type Meet struct {
	Id      int    `json:"Id"`
	Name    string `json:"Name"`
	Weekday string `json:"Weekday"`
	Date    string `json:"Date"`
	Start   string `json:"Start"`
	End     string `json:"End"`
	Url     string `json:"Url"`
	ZoomId  string `json:"ZoomId"`
	Pass    string `json:"Pass"`
}

func NewMeet() Meet {
	return Meet{}
}

func (meet *Meet) IsNotEmpty() bool {
	return len(meet.Name) > 0
}
