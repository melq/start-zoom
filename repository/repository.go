package repository

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"start-zoom/common"
)

func GetSameNames(meets []common.Meet, name string) []int {
	var ids []int
	for i, meet := range meets {
		if meet.Name == name {
			ids = append(ids, i)
		}
	}
	return ids
}

// 同ディレクトリにファイルの存在を確認する関数
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// LoadConfig jsonファイルを読み込んで構造体を返す関数
func LoadConfig(filename string) common.Config {
	config := common.NewConfig()
	if !fileExists(filename) {
		if _, err := os.Create(filename); err != nil {
			log.Fatal(err)
		}
	}
	bytes, err := ioutil.ReadFile(filename) // json読み込み
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
func SaveConfig(config *common.Config, filename string) {
	configJson, err := json.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}
	fp, err := os.OpenFile(filename, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
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
