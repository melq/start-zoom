package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type Config struct {
	ClassData ClassData
}
/*授業の情報を格納する構造体*/
type ClassData struct {
	Name    string `toml:"Name"`
	Weekday string `toml:"Weekday"`
	Start   string `toml:"Start"`
	End     string `toml:"End"`
	Url     string `toml:"Url"`
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
func loadClassesToml(filename string) (classes []ClassData) {
	tomlBytes, err := ioutil.ReadFile(filename) //json読み込み
	if err != nil {
		log.Fatal(err)
	}
	if len(tomlBytes) != 0 {
		if err := toml.Unmarshal(tomlBytes, &classes); err != nil {
			log.Fatal(err)
		}
	}
	return
}
/*jsonファイルを読み込んで構造体の配列を返す関数*/
func loadClasses(filename string) (classes []ClassData) {
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
func saveClassesToml(classes []ClassData, filename string) {
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(classes); err != nil {
		log.Fatal(err)
	}
	fp, err := os.OpenFile(filename, os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()
	if _, err = fp.WriteString(buf.String()); err != nil {
		log.Fatal(err)
	}
}
/*jsonファイルに書き込む関数*/
func saveClasses(classes []ClassData, filename string) {
	classJson, err := json.Marshal(classes)
	if err != nil {
		log.Fatal(err)
	}
	fp, err := os.OpenFile(filename, os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	_, err = fp.WriteString(string(classJson))
	if err != nil {
		panic(err)
	}
}
/*新規登録する授業の構造体を作成する関数*/
func registerClass() (cd ClassData) {
	fmt.Println("新しく授業を登録します。")
	fmt.Print("授業名を入力:")
	cd.Name = read()
	fmt.Println()
	tmp := InputNum("曜日を選択(開始時の曜日): 1: Sunday, 2: Monday, 3: Tuesday, 4: Wednesday, 5: Thursday, 6: Friday, 7: Saturday")
	switch tmp {
	case 1: cd.Weekday = "Sunday"
	case 2: cd.Weekday = "Monday"
	case 3: cd.Weekday = "Tuesday"
	case 4: cd.Weekday = "Wednesday"
	case 5: cd.Weekday = "Thursday"
	case 6: cd.Weekday = "Friday"
	case 7: cd.Weekday = "Saturday"
	}
	tmp = InputNum("開始時間を入力(例：14:30 => 1430 (半角数字))")
	cd.Start = strconv.Itoa(tmp / 100) + ":" + strconv.Itoa(tmp % 100)
	if tmp % 100 == 0 { cd.Start += "0" }
	tmp = InputNum("終了時間を入力:")
	cd.End = strconv.Itoa(tmp / 100) + ":" + strconv.Itoa(tmp % 100)
	if tmp % 100 == 0 { cd.End += "0" }
	fmt.Print("ZoomURLを入力:")
	cd.Url = read()

	return
}
/*ブラウザでZoomを開く関数*/
func startZoom(classes []ClassData) {
	fmt.Println("Zoomを開きます.")
	trueNow := time.Now()
	for _, class := range classes {
		if class.Weekday == trueNow.Weekday().String() {
			now, _ := time.Parse("15:04", strconv.Itoa(trueNow.Hour())+ ":" +strconv.Itoa(trueNow.Minute()))
			startTime, _ := time.Parse("15:04", class.Start)
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
/*登録授業のリストを表示する関数*/
func showClassList(classes []ClassData) {
	fmt.Println("\n登録されている授業を表示します.")
	fmt.Print("\n")
	if len(classes) == 0 {
		fmt.Println("登録授業なし")
	} else {
		for i, class := range classes {
			fmt.Println(i+1, ":", class.Name)
			fmt.Println("", class.Weekday, class.Start, "~", class.End)
			fmt.Println("", class.Url)
		}
	}
}
/*メイン関数*/
func StartZoomMain() {
	filename := "classes.toml"
	classes := loadClassesToml(filename)

	fmt.Println("\n" +
		"---------------------------------------------\n" +
		"----------------- StartZoom -----------------\n" +
		"----------- (made by RikuTsuzuki) -----------\n" +
		"---------------------------------------------")

	flg := 0
	for flg == 0 {
		switch InputNum("\n0: 終了, 1: 授業開始, 2: 授業登録, 3: 授業リスト") {
		case 0:
			fmt.Println("終了します.")
			flg = 1
		case 1:
			startZoom(classes)
		case 2:
			classes = append(classes, registerClass())
			saveClassesToml(classes, filename)
		case 3:
			showClassList(classes)
		}
	}
}
