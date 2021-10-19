package main

import (
	"fmt"
	"log"
	"os/exec"
	"startZoom/repository"
	"strconv"
	"time"
)

// InputNum /*数値入力用関数*/
func InputNum (msg string) int {
	for {
		fmt.Println(msg)
		i, e := strconv.Atoi(read())
		if e != nil {
			continue
		}
		return i
	}
}

func checkTime(meet repository.Meet, timeMargin int) bool {
	now := time.Now()
	now.In(time.FixedZone("Asia/Tokyo", 9*60*60))

	nowTime, _ := time.Parse("15:4", strconv.Itoa(now.Hour()) + ":" + strconv.Itoa(now.Minute()))
	startTime, _ := time.Parse("15:04", meet.Start)
	startTime = startTime.Add(time.Duration(-1 * timeMargin) * time.Minute)
	endTime, _ := time.Parse("15:04", meet.End)

	if startTime.Before(nowTime) && endTime.After(nowTime) {
		return true
	}
	return false
}

func getEarlierMeet(meet1 repository.Meet, meet2 repository.Meet) repository.Meet {
	now := time.Now()
	now.In(time.FixedZone("Asia/Tokyo", 9*60*60))

	nowTime, _ := time.Parse("15:4", strconv.Itoa(now.Hour()) + ":" + strconv.Itoa(now.Minute()))
	time1, _ := time.Parse("15:04", meet1.Start)
	time2, _ := time.Parse("15:04", meet2.Start)

	if nowTime.After(time1) && nowTime.After(time2) {
		return repository.NewMeet()
	} else if nowTime.After(time1) {
		return meet2
	} else if nowTime.After(time2) {
		return meet1
	} else {
		if time1.Before(time2) {
			return meet1
		}
		return meet2
	}
}

func runMeet(meet repository.Meet) {
	fmt.Println(meet.Name, "のURLを開きます")
	err := exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", meet.Url).Start()
	if err != nil {
		log.Fatalln(err)
	}
}

func startMeet(config repository.Config) {
	var currentMeet repository.Meet
	var nextMeet repository.Meet

	now := time.Now()
	now.In(time.FixedZone("Asia/Tokyo", 9*60*60))

	fmt.Printf("現在時刻: %02d : %02d", now.Hour(), now.Minute())
	_, month, day := now.Date()
	todayStr := strconv.Itoa(int(month)) + "-" + strconv.Itoa(day)

	proc := func () {

	}

	for _, meet := range config.Meets {
		if meet.Date == todayStr {
			if checkTime(meet, config.TimeMargin) {
				currentMeet = meet
			}
			if nextMeet.Name != "" {
				nextMeet = getEarlierMeet(nextMeet, meet)
			} else {
				nextMeet = meet
			}
		}
	}
	for _, meet := range config.Meets {
		if meet.Weekday == now.Weekday().String() {
			if checkTime(meet, config.TimeMargin) {
				currentMeet = meet
			}
			if nextMeet.Name != "" {
				nextMeet = getEarlierMeet(nextMeet, meet)
			} else {
				nextMeet = meet
			}
		}
	}
	if currentMeet.IsNotEmpty() {
		if nextMeet.IsNotEmpty() && checkTime(nextMeet, config.TimeMargin) {
			runMeet(nextMeet)
		} else {
			runMeet(currentMeet)
		}
	} else {
		fmt.Println("現在または", config.TimeMargin, "分後に進行中の授業はありません")
		fmt.Print("\n")
		if nextMeet.IsNotEmpty() && config.IsAsk {
			msg := nextMeet.Start + " から " + nextMeet.Name + " が始まりますが、起動しますか？" +
				"\n1: はい, 2: いいえ"
			if InputNum(msg) == 1 {
				runMeet(nextMeet)
			} else {
				fmt.Println("起動せず戻ります")
			}
		}
	}
}

// StartZoomMain メイン関数
func StartZoomMain(opts Options) {
	filename := "config.json"
	config := repository.LoadConfig(filename)

	if len(opts.Start) != 0 {
		startZoom(config)
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
		switch InputNum("\n行いたい操作の番号を入力してください\n0: 終了, 1: 授業開始, 2: 授業登録, 3: 授業リスト, 4: 登録授業の編集・削除, 5: 選択して授業開始, 6: 設定") {
		case 0:
			fmt.Println("終了します")
			flg = 1
		case 1:
			startMeet(config)
		case 2:
			fmt.Println("新しく授業を登録します。")
			config.SumId++
			config.Classes = append(config.Classes, makeClass(config.SumId))
			saveConfig(config, filename)
		case 3:
			showClassList(config.Classes)
		case 4:
			config.Classes = editDeleteClasses(config.Classes)
			saveConfig(config, filename)
		case 5:
			anytimeStart(config.Classes)
		case 6:
			config = editConfig(config)
			saveConfig(config, filename)
		default:
		}
	}
}

