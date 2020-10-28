package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

/*授業の情報を格納する構造体*/
type classData struct {
	Name    string `json:"Name"`
	Weekday string `json:"Weekday"`
	Start   string `json:"Start"`
	End     string `json:"End"`
	Url     string `json:"Url"`
}

var sc = bufio.NewScanner(os.Stdin)

/*入力読み込み用関数*/
func read() string {
	sc.Scan()
	return sc.Text()
}
/*数値入力用関数*/
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
/*ファイルの存在を確認する関数*/
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
/*jsonファイルを読み込んで構造体の配列を返す関数*/
func loadClasses(filename string) (classes []classData) {
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
		if err := json.Unmarshal(bytes, &classes); err != nil {
			log.Fatal(err)
		}
	}
	return
}
/*jsonファイルに書き込む関数*/
func saveClasses(classes []classData, filename string) {
	classJson, err := json.Marshal(classes)
	if err != nil {
		log.Fatal(err)
	}
	fp, err := os.OpenFile(filename, os.O_WRONLY | os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	_, err = fp.Write(classJson)
	if err != nil {
		panic(err)
	}
}
/*授業の名前を入力する関数*/
func inputName() (name string) {
	fmt.Print("授業名を入力:")
	name = read()
	return
}
/*授業の曜日を入力する関数*/
func inputWeekday() (weekday string) {
	switch InputNum("曜日を選択(開始時の曜日): 1: Sunday, 2: Monday, 3: Tuesday, 4: Wednesday, 5: Thursday, 6: Friday, 7: Saturday") {
	case 1: weekday = "Sunday"
	case 2: weekday = "Monday"
	case 3: weekday = "Tuesday"
	case 4: weekday = "Wednesday"
	case 5: weekday = "Thursday"
	case 6: weekday = "Friday"
	case 7: weekday = "Saturday"
	default: weekday = inputWeekday()
	}
	return
}
/*授業の開始時刻を入力する関数*/
func inputStartTime() (startTime string) {
	tmp := InputNum("開始時刻を入力(例：14:30 => 1430 (半角数字))")
	startTime = strconv.Itoa(tmp / 100) + ":" + strconv.Itoa(tmp % 100)
	if tmp % 100 == 0 { startTime += "0" }
	return
}
/*授業の終了時刻を入力する関数*/
func inputEndTime() (endTime string) {
	tmp := InputNum("終了時刻を入力")
	endTime = strconv.Itoa(tmp / 100) + ":" + strconv.Itoa(tmp % 100)
	if tmp % 100 == 0 { endTime += "0" }
	return
}
/*授業のURLを入力する関数*/
func inputUrl() (url string) {
	fmt.Print("ZoomURLを入力:")
	url = read()
	return
}
/*新規登録する授業の構造体を作成する関数*/
func makeClass() (cd classData) {
	cd.Name = inputName()
	cd.Weekday = inputWeekday()
	cd.Start = inputStartTime()
	cd.End = inputEndTime()
	cd.Url = inputUrl()
	return
}
/*ブラウザでZoomを開く関数*/
func startZoom(classes []classData) {
	trueNow := time.Now()
	fmt.Println("現在時刻:", trueNow.Hour(), ":", trueNow.Minute())
	for _, class := range classes {
		if class.Weekday == trueNow.Weekday().String() {
			now, _ := time.Parse("15:04", strconv.Itoa(trueNow.Hour())+ ":" +strconv.Itoa(trueNow.Minute()))
			startTime, _ := time.Parse("15:04", class.Start)
			startTime = startTime.Add(-10 * time.Minute)
			endTime, _ := time.Parse("15:04", class.End)
			if startTime.Before(now) && endTime.After(now) {
				fmt.Println(class.Name, "のZoomを開きます")
				err := exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", class.Url).Start()
				if err != nil {
					panic(err)
				}
				return
			}
		}
	}
}
/*授業単体の情報を表示する関数*/
func showClassData(cd classData) {
	fmt.Println(cd.Name)
	fmt.Println("", cd.Weekday, cd.Start, "~", cd.End)
	fmt.Println("", cd.Url)
}
/*登録授業のリストを表示する関数*/
func showClassList(classes []classData) {
	fmt.Println("\n登録されている授業を表示します.")
	fmt.Print("\n")
	if len(classes) == 0 {
		fmt.Println("登録授業なし")
	} else {
		for i, cd := range classes {
			fmt.Print(i+1, ": ")
			showClassData(cd)
		}
	}
}
/*登録授業単体を編集する関数*/
func editClassData(cd classData) (editedCd classData) {
	editedCd = cd
	switch InputNum(editedCd.Name + "の何を編集しますか？\n" +
					"1: 名前, 2: 曜日, 3: 開始時刻, 4: 終了時刻, 5: URL, 6: すべて") {
	case 1: editedCd.Name = inputName()
	case 2: editedCd.Weekday = inputWeekday()
	case 3: editedCd.Start = inputStartTime()
	case 4: editedCd.End = inputEndTime()
	case 5: editedCd.Url = inputUrl()
	case 6:
		fmt.Println("すべて編集します")
		editedCd = makeClass()
	default:
		editedCd = editClassData(cd)
	}
	return editedCd
}
/*登録授業リストを編集する関数*/
func editClasses(classes []classData) (editedClasses []classData) {
	showClassList(classes)
	fmt.Println("\n登録授業の編集をします")
	classNum := InputNum("編集したい授業の番号を入力してください(編集せず戻る場合は0)")
	if classNum == 0 {
		return classes
	} else {
		classNum -= 1
		if classNum >= len(classes) || classNum < 0 {
			fmt.Println("授業の番号が不正です")
			return classes
		} else {
			editedClasses = classes
			editedClasses[classNum] = editClassData(classes[classNum])
			fmt.Println("\n編集が正常に終了しました")
			fmt.Print("\n")
			showClassData(editedClasses[classNum])
		}
	}
	return
}
/*登録授業単体の削除を行う関数*/
func deleteClassData(classes []classData, index int) (editedClasses []classData) {
	for i, cd := range classes {
		if i == index { continue }
		editedClasses = append(editedClasses, cd)
	}
	return
}
/*登録授業の削除を行う関数*/
func deleteClasses(classes []classData) (editedClasses []classData) {
	showClassList(classes)
	fmt.Println("\n登録授業の削除をします")
	classNum := InputNum("削除したい授業の番号を入力してください(削除せず戻る場合は0)")
	if classNum == 0 {
		return classes
	} else {
		classNum -= 1
		if classNum >= len(classes) || classNum < 0 {
			fmt.Println("授業の番号が不正です")
			return classes
		} else {
			editedClasses = classes
			editedClasses = deleteClassData(classes, classNum)
			fmt.Println("\n削除が正常に終了しました")
		}
	}
	return
}
/*登録授業を編集・削除する関数*/
func editDeleteClasses(classes []classData) (editedClasses []classData) {
	fmt.Println("登録授業の編集・削除を行います")
	if len(classes) == 0 {
		fmt.Println("登録授業なし")
		return classes
	}
	switch InputNum("0: 戻る, 1: 編集, 2: 削除") {
	case 1: editedClasses = editClasses(classes)
	case 2: editedClasses = deleteClasses(classes)
	default: return classes
	}
	return
}
/*メイン関数*/
func StartZoomMain() {
	filename := "classes.json"
	classes := loadClasses(filename)

	fmt.Println("\n" +
		"---------------------------------------------\n" +
		"----------------- StartZoom -----------------\n" +
		"----------- (made by RikuTsuzuki) -----------\n" +
		"---------------------------------------------")

	flg := 0
	for flg == 0 {
		switch InputNum("\n行いたい操作の番号を入力してください\n0: 終了, 1: 授業開始, 2: 授業登録, 3: 授業リスト, 4: 登録授業の編集・削除") {
		case 0:
			fmt.Println("終了します.")
			flg = 1
		case 1:
			startZoom(classes)
		case 2:
			fmt.Println("新しく授業を登録します。")
			classes = append(classes, makeClass())
			saveClasses(classes, filename)
		case 3:
			showClassList(classes)
		case 4:
			classes = editDeleteClasses(classes)
			saveClasses(classes, filename)
		default:
		}
	}
}
