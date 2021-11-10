package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
	"start-zoom/functions"
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
	filename := "config.json"
	config := repository.LoadConfig(filename)

	if len(opts.Start) != 0 {
		functions.StartMeet(config)
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
		switch functions.InputNum("\n行いたい操作の番号を入力してください\n0: 終了, 1: 会議開始, 2: 会議登録, 3: 会議リスト, 4: 登録会議の編集・削除, 5: 選択して会議開始, 6: 設定") {
		case 0:
			fmt.Println("終了します")
			flg = 1
		case 1:
			functions.StartMeet(config)
		case 2:
			functions.MakeMeet(&config, filename)
		case 3:
			functions.ShowMeets(config)
		case 4:
			functions.EditOrDeleteMeet(&config, filename)
		case 5:
			functions.StartSpecifiedMeet(config)
		case 6:
			functions.EditSetting(&config, filename)
		default:
		}
	}
}

