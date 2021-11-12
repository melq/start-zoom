package repository

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

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
	Id 		int		`json:"Id"`
	Name    string	`json:"Name"`
	Weekday string	`json:"Weekday"`
	Date 	string	`json:"Date"`
	Start   string	`json:"Start"`
	End     string	`json:"End"`
	Url     string	`json:"Url"`
	ZoomId	string	`json:"ZoomId"`
	Pass	string	`json:"Pass"`
}
type Config struct {
	Meets 		[]Meet 		`json:"Meets"`
	SumId 		int    		`json:"SumId"`
	TimeMargin	int			`json:"TimeMargin"`
	IsAsk		bool		`json:"IsAsk"`
}

func NewMeet() Meet {
	return Meet{}
}

func NewConfig() Config {
	return Config{
		Meets: make([]Meet, 0),
		SumId: 0,
		TimeMargin: 20,
		IsAsk: false,
	}
}

func (meet *Meet) IsNotEmpty() bool {
	return len(meet.Name) > 0
}

func GetSameNames(meets []Meet, name string) []int {
	var ids []int
	for i, meet := range meets {
		if meet.Name == name {
			ids = append(ids, i)
		}
	}
	return ids
}

/*func MakeBatchIfNotExist() {  // D:直下にバッチを作成する機能を廃止
	_, err := os.Stat("D:/myzoom.bat")
	if err == nil {
		return
	}

	var bytes []byte
	if fileExists("myzoom.bat") {
		bytes, err = ioutil.ReadFile("myzoom.bat")
		if err != nil {
			log.Fatalln(err)
		}
	}
	fp, err := os.OpenFile("D:/myzoom.bat", os.O_TRUNC | os.O_WRONLY | os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := fp.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if _, err = fp.Write(bytes); err != nil {
		log.Fatal(err)
	}
}*/

// 同ディレクトリにファイルの存在を確認する関数
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// LoadConfig jsonファイルを読み込んで構造体を返す関数
func LoadConfig(filename string) Config {
	config := NewConfig()
	if !fileExists(filename) {
		if _, err := os.Create(filename); err != nil {
			log.Fatal(err)
		}
	}
	bytes, err := ioutil.ReadFile(filename)	//json読み込み
	if err != nil {
		log.Fatal(err)
	}
	if len(bytes) != 0 {
		if err := json.Unmarshal(bytes, &config); err != nil {
			log.Fatal(err)
		}
	}
	return config
}

// SaveConfig jsonファイルに書き込む関数
func SaveConfig(config *Config, filename string) {
	configJson, err := json.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}
	fp, err := os.OpenFile(filename, os.O_TRUNC | os.O_WRONLY | os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := fp.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if _, err = fp.Write(configJson); err != nil {
		log.Fatal(err)
	}
}