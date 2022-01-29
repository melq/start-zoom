package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
	"start-zoom"
	"start-zoom/repository"
)

type Options struct {
	Start []bool `short:"s" long:"start" description:"Get starting zoom"`
}

var opts Options

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}
	filename := "files/config.json"
	config := repository.LoadConfig(filename)

	if len(opts.Start) != 0 {
		start_zoom.StartMeet(config)
		return
	}

	fmt.Println("\n" +
		"---------------------------------------------\n" +
		"----------------- StartZoom -----------------\n" +
		"----------- (made by RikuTsuzuki) -----------\n" +
		"---------------------------------------------")
	fmt.Print("\n")

	flg := 0
	for flg == 0 {
		switch start_zoom.InputNum("\n行いたい操作の番号を入力してください\n0: 終了, 1: 会議開始, 2: 会議登録, 3: 会議リスト, 4: 登録会議の編集・削除, 5: 選択して会議開始, 6: 設定") {
		case 0:
			fmt.Println("終了します")
			flg = 1
		case 1:
			start_zoom.StartMeet(config)
		case 2:
			start_zoom.MakeMeet(&config, filename)
		case 3:
			start_zoom.ShowMeets(config)
		case 4:
			start_zoom.EditOrDeleteMeet(&config, filename)
		case 5:
			start_zoom.StartSpecifiedMeet(config)
		case 6:
			start_zoom.EditSetting(&config, filename)
		default:
		}
	}
}
