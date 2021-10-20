package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"startZoom/repository"
	"strconv"
	"time"
)

var sc = bufio.NewScanner(os.Stdin)

/*入力読み込み用関数*/
func read() string {
	sc.Scan()
	return sc.Text()
}

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
		fmt.Println("現在または", config.TimeMargin, "分後に進行中の会議はありません")
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

func inputName() string {
	fmt.Print("\n授業名を入力:")
	return read()
}

func inputWeekday() (string, string) {
	fmt.Println("\nZoomが開催される曜日を指定。毎週開催されるものでなくある日程のみの会議の場合は、日付のみの指定も可能です")
	n := InputNum("曜日(または日付指定)を選択: 0: 日付で指定する, 1: Sunday, 2: Monday, 3: Tuesday, 4: Wednesday, 5: Thursday, 6: Friday, 7: Saturday")
	var weekday string; var date string
	switch n {
	case 0:
		tmp := InputNum("日付を入力(例：1月2日 => 0102 (半角数字))")
		date = strconv.Itoa(tmp / 100) + "-" + strconv.Itoa(tmp % 100)
	case 1, 2, 3, 4, 5, 6, 7:
		weekday = repository.WeekdayString[n - 1]
	default: weekday, date = inputWeekday()
	}
	return weekday, date
}

func inputStartTime() string {
	fmt.Println()
	tmp := InputNum("開始時刻を入力(例：14:30 => 1430 (半角数字), 存在しない時刻は入力しないでください)")

	startTime := strconv.Itoa(tmp / 100) + ":" + strconv.Itoa(tmp % 100)
	if tmp % 100 == 0 { startTime += "0" }
	return startTime
}

func inputEndTime() string {
	fmt.Println()
	tmp := InputNum("終了時刻を入力(例：14:30 => 1430 (半角数字), 存在しない時刻は入力しないでください)")

	endTime := strconv.Itoa(tmp / 100) + ":" + strconv.Itoa(tmp % 100)
	if tmp % 100 == 0 { endTime += "0" }
	return endTime
}

func inputUrl() string {
	fmt.Println("\n会議のURLを入力")
	return read()
}

func makeMeet(config repository.Config, filename string) {
	fmt.Println("新しく会議を登録します")
	config.SumId++
	meet := repository.NewMeet()

	meet.Id = config.SumId
	meet.Name = inputName()
	meet.Weekday, meet.Date = inputWeekday()
	meet.Start = inputStartTime()
	meet.End = inputEndTime()
	meet.Url = inputUrl()

	config.Meets = append(config.Meets, meet)
	repository.SaveConfig(config, filename)

	fmt.Println(meet.Name, "を作成しました")
}

func showMeet(meet repository.Meet) {
	fmt.Println(meet.Name)
	fmt.Println(" URL:", meet.Url)
	if len(meet.Weekday) > 0 {
		fmt.Println(" 曜日:", meet.Weekday)
	} else {
		meetDate, err := time.Parse("01-02", meet.Date)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(" 日時:", meetDate)
	}
	fmt.Println(" 時刻:", meet.Start, "-", meet.End)
}

func showMeets(config repository.Config) {
	fmt.Println("登録されている会議を表示します")
	fmt.Println()
	if len(config.Meets) == 0 {
		fmt.Println("登録会議なし")
	} else {
		for i, meet := range config.Meets {
			fmt.Print("\n", i + 1, ": ")
			showMeet(meet)
		}
	}
}

func editMeet(config repository.Config, filename string) {
	fmt.Println("登録会議の編集をします")
	showMeets(config)
	fmt.Println()

	meetNum := InputNum("編集を行う会議の番号を入力してください(編集せず戻る場合「0」)")
	if meetNum == 0 {
		fmt.Println("戻ります")
		return
	}
	meetNum--
	if meetNum >= len(config.Meets) || meetNum < 0 {
		fmt.Println("番号が不正です")
		return
	} else {
		tmpMeet := config.Meets[meetNum]
		switch InputNum(tmpMeet.Name + "の何を編集しますか？\n" +
			"0: 戻る, 1: 名前, 2: 曜日または日付, 3: 開始時刻, 4: 終了時刻, 5: URL") {
		case 1: tmpMeet.Name = inputName()
		case 2: tmpMeet.Weekday, tmpMeet.Date = inputWeekday()
		case 3: tmpMeet.Start = inputStartTime()
		case 4: tmpMeet.End = inputEndTime()
		case 5: tmpMeet.Url = inputUrl()
		default:
			fmt.Println("戻ります")
			return
		}
		config.Meets[meetNum] = tmpMeet

		repository.SaveConfig(config, filename)
		fmt.Println(config.Meets[meetNum].Name, "に編集しました")
	}
}

func deleteMeet(config repository.Config) {
}

func editOrDeleteMeet(config repository.Config) {
	fmt.Println("\n登録会議の編集・削除を行います")
	if len(config.Meets) == 0 {
		fmt.Println("登録会議なし")
		return
	}
	switch InputNum("0: 戻る, 1: 編集, 2: 削除") {
	case 1: editMeet(config)
	case 2: deleteMeet(config)
	default: return
	}
}

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
		switch InputNum("\n行いたい操作の番号を入力してください\n0: 終了, 1: 会議開始, 2: 会議登録, 3: 授業リスト, 4: 登録会議の編集・削除, 5: 選択して会議開始, 6: 設定") {
		case 0:
			fmt.Println("終了します")
			flg = 1
		case 1:
			startMeet(config)
		case 2:
			makeMeet(config, filename)
		case 3:
			showMeets(config)
		case 4:
			editMeet(config, filename)
		case 5:
			anytimeStart(config.Classes)
		case 6:
			config = editConfig(config)
			saveConfig(config, filename)
		default:
		}
	}
}

