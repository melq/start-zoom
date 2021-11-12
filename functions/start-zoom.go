package functions

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"start-zoom/repository"
	"strconv"
	"strings"
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
	fmt.Println(meet.Name, "の会議を開きます")
	url := ""
	if len(meet.Url) > 0 {
		url = meet.Url
	} else {
		url = "https://us02web.zoom.us/j/" + meet.ZoomId + "?pwd=" + meet.Pass
	}
	err := exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", url).Start()
	if err != nil {
		log.Fatalln(err)
	}
}

func StartMeet(config repository.Config) {
	var currentMeet repository.Meet
	var nextMeet repository.Meet

	now := time.Now()
	now.In(time.FixedZone("Asia/Tokyo", 9*60*60))

	fmt.Printf("現在時刻: %02d : %02d\n\n", now.Hour(), now.Minute())
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

func inputName(meets []repository.Meet /*mode == 0: 作成, 1: 編集*/) string {
	var name string
	flg := 0
	for flg == 0 {
		fmt.Print("\n会議名を入力:")
		name = read()

		ids := repository.GetSameNames(meets, name)
		if len(ids) > 0 {
			fmt.Println("以下の同名の会議があります。同名の会議は登録できません")
			for _, i := range ids {
				fmt.Println()
				showMeet(meets[i])
			}
		} else { flg = 1 }
		fmt.Println()
	}
	return name
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

	startTime := fmt.Sprintf("%02d:%02d", tmp / 100, tmp % 100)
	return startTime
}

func inputEndTime() string {
	fmt.Println()
	tmp := InputNum("終了時刻を入力(例：14:30 => 1430 (半角数字), 存在しない時刻は入力しないでください)")

	endTime := fmt.Sprintf("%02d:%02d", tmp / 100, tmp % 100)
	return endTime
}

func inputUrlOrId() (string, string, string) {
	fmt.Println("\n会議のURLまたはIDとパスワードを入力")
	url := ""; id := ""; pass := ""
	switch InputNum("1: URLで登録, 2: IDとパスワードで登録") {
	case 1:
		fmt.Println("URLを入力してください")
		url = read()
	case 2:
		fmt.Println("IDを入力してください")
		id = read()
		fmt.Println("パスワードを入力してください")
		pass = read()
	}
	return url, id, pass
}

func adjustYear(dateStr string) time.Time {
	now := time.Now()
	date, err := time.Parse("2006-1-2", strconv.Itoa(now.Year()) + "-" + dateStr)
	if err != nil {
		log.Fatalln(err)
	}
	if now.After(date) {
		date = date.AddDate(1, 0, 0)
	}
	return date
}

func makeSchtasks(meet repository.Meet) {
	// repository.MakeBatchIfNotExist() // D:直下にバッチを作成する機能を廃止

	var id string; var pass string
	if len(meet.Url) > 0 {
		idIndex := strings.Index(meet.Url, "j/") + 2
		qIndex := strings.Index(meet.Url, "?")
		id = meet.Url[idIndex:qIndex]

		pwdIndex := strings.Index(meet.Url, "pwd=") + 4
		pass = meet.Url[pwdIndex:]
	} else {
		id = meet.ZoomId
		pass = meet.Pass
	}
	stime, _ := time.Parse("15:04", meet.Start)
	stime = stime.Add(time.Duration(-5) * time.Minute)
	stimeStr := fmt.Sprintf("%02d:%02d", stime.Hour(), stime.Minute())
	dates := adjustYear(meet.Date)
	year, month, date := dates.Date()
	dateWithYear := fmt.Sprintf("%04d/%02d/%02d", year, month, date)

	fmt.Println("./settask.bat", meet.Name, id, pass, stimeStr, dateWithYear)
	_, err := exec.Command("settask.bat", meet.Name, id, pass, stimeStr, dateWithYear).Output()

	if err != nil {
		log.Fatalln("settask", err)
	} else {
		fmt.Println("登録しました")
	}
}

func MakeMeet(config *repository.Config, filename string) {
	fmt.Println("新しく会議を登録します")
	config.SumId++
	meet := repository.NewMeet()

	meet.Id = config.SumId
	meet.Name = inputName(config.Meets)
	meet.Weekday, meet.Date = inputWeekday()
	meet.Start = inputStartTime()
	meet.End = inputEndTime()
	meet.Url, meet.ZoomId, meet.Pass = inputUrlOrId()

	config.Meets = append(config.Meets, meet)
	repository.SaveConfig(config, filename)
	fmt.Println(meet.Name, "を作成しました")

	if len(meet.Date) > 0 {
		fmt.Println("\nこの予定をタスクスケジューラに登録しますか？(Zoomの場合のみ)")
		switch InputNum("1: はい, 2: いいえ") {
		case 1:
			makeSchtasks(meet)
		case 2:
		}
	}
}

func showMeet(meet repository.Meet) {
	fmt.Println(meet.Name)
	if len(meet.Url) > 0 {
		fmt.Println(" URL:", meet.Url)
	} else {
		fmt.Println(" ID:", meet.ZoomId, ",", "Pass:", meet.Pass)
	}
	if len(meet.Weekday) > 0 {
		fmt.Println(" 曜日:", meet.Weekday)
	} else {
		dates := adjustYear(meet.Date)
		year, month, date := dates.Date()
		fmt.Println(" 日時:", fmt.Sprintf("%04d-%02d-%02d", date, month, year))
	}
	fmt.Println(" 時刻:", meet.Start, "-", meet.End)
}

func ShowMeets(config repository.Config) {
	fmt.Println("登録されている会議を表示します")
	if len(config.Meets) == 0 {
		fmt.Println("登録会議なし")
	} else {
		for i, meet := range config.Meets {
			fmt.Print("\n", i + 1, ": ")
			showMeet(meet)
		}
	}
}

func editMeet(config *repository.Config, filename string) {
	fmt.Println("登録会議の編集をします")
	ShowMeets(*config)
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
			"0: 戻る, 1: 名前, 2: 曜日または日付, 3: 開始時刻, 4: 終了時刻, 5: URLまたはZoomIDとパスワード") {
		case 1: tmpMeet.Name = inputName(config.Meets)
		case 2: tmpMeet.Weekday, tmpMeet.Date = inputWeekday()
		case 3: tmpMeet.Start = inputStartTime()
		case 4: tmpMeet.End = inputEndTime()
		case 5: tmpMeet.Url, tmpMeet.ZoomId, tmpMeet.Pass = inputUrlOrId()
		default:
			fmt.Println("戻ります")
			return
		}
		config.Meets[meetNum] = tmpMeet

		repository.SaveConfig(config, filename)
		fmt.Println(config.Meets[meetNum].Name, "に編集しました")
	}
}

func deleteSchtasks(name string) {
	_, err := exec.Command("deletetask.bat", name).Output()

	if err != nil {
		//log.Fatalln("deletetask", err)
	} else {
		fmt.Println(name, "をタスクスケジューラから削除しました")
	}
}

func deleteMeet(config *repository.Config, filename string) {
	fmt.Println("登録会議の削除をします")
	ShowMeets(*config)
	fmt.Println()
	meetNum := InputNum("削除したい会議の番号を入力してください(すべて削除する場合は「-1」)(削除せず戻る場合は「0」)")
	if meetNum == 0 {
		fmt.Println("削除せずに戻ります")
		return
	}
	if meetNum == -1 {
		fmt.Println("すべての会議データを削除します.よろしいですか？")
		switch InputNum("1: はい, 2: いいえ") {
		case 1:
			for _, meet := range config.Meets {
				deleteSchtasks(meet.Name)
			}
			config.Meets = []repository.Meet{}
			repository.SaveConfig(config, filename)
			fmt.Println("すべてのデータを削除しました")
		default:
			fmt.Println("削除せず戻ります")
		}
		return
	} else {
		meetNum--
		if meetNum >= len(config.Meets) || meetNum < 0 {
			fmt.Println("番号が不正です")
			return
		} else {
			fmt.Println(config.Meets[meetNum].Name, "のデータを削除します.よろしいですか？")
			switch InputNum("1: はい, 2: いいえ") {
			case 1:
				fmt.Println(config.Meets[meetNum].Name, "のデータを削除します")
				var tmpMeets []repository.Meet
				for i, meet := range config.Meets {
					if i == meetNum { continue }
					tmpMeets = append(tmpMeets, meet)
				}
				deleteSchtasks(config.Meets[meetNum].Name)
				config.Meets = tmpMeets
				repository.SaveConfig(config, filename)
				fmt.Println("\n削除しました")
			default:
				fmt.Println("削除せず戻ります")
			}
			return
		}
	}
}

func EditOrDeleteMeet(config *repository.Config, filename string) {
	fmt.Println("\n登録会議の編集・削除を行います")
	if len(config.Meets) == 0 {
		fmt.Println("登録会議なし")
		return
	}
	switch InputNum("0: 戻る, 1: 編集, 2: 削除") {
	case 1: editMeet(config, filename)
	case 2: deleteMeet(config, filename)
	default: return
	}
}

func StartSpecifiedMeet(config repository.Config) {
	fmt.Println("指定された会議を開きます")
	ShowMeets(config)
	fmt.Println()
	meetNum := InputNum("開く会議の番号を入力(戻る場合「0」)")
	if meetNum == 0 {
		fmt.Println("戻ります")
		return
	}
	runMeet(config.Meets[meetNum])
}

func EditSetting(config *repository.Config, filename string) {
	fmt.Println("設定の変更をします")
	switch InputNum("0: 戻る, 1: 会議開始前の余裕時間, 2: 該当会議がない場合の質問") {
	case 0:
		fmt.Println("戻ります")
		return
	case 1:
		fmt.Println("\n会議開始時刻の何分前から開くようにするか設定します(現在は", config.TimeMargin, "分")
		config.TimeMargin = InputNum("何分前から起動可能に設定しますか？: ")
	case 2:
		fmt.Println("授業開始を選択した際に、開始時刻に該当する会議がなかったときに、同じ日のなかで" +
			"最も開始時刻の近いものを開くかどうかの質問の有無を設定します")
		switch InputNum("1: 聞く, 2: 聞かない") {
		case 1: config.IsAsk = true
		case 2: config.IsAsk = false
		}
	}
	repository.SaveConfig(config, filename)
	fmt.Println("設定を変更しました")
}