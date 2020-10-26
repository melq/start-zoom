package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

var sc = bufio.NewScanner(os.Stdin)

func read() string {
	sc.Scan()
	return sc.Text()
}

func InputNum (sc *bufio.Scanner, msg string) int {
	for {
		fmt.Println(msg)
		sc.Scan()
		i, e := strconv.Atoi(sc.Text())
		if e != nil {
			continue
		}
		return i
	}
}

type classData struct {
	Name string `json:name`
	Day string `json:day`
	Start string `json:start`
	End string `json:end`
	Url string `json:url`
}

/*jsonファイルに書き込む関数*/
func SaveClass (classes []classData, filename string) {
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
func RegisterClass() (cd classData) {
	fmt.Println("新しく授業を登録します。")
	fmt.Print("授業名を入力:")
	cd.Name = read()
	fmt.Print("曜日を入力:")
	cd.Day = read()
	fmt.Print("開始時間を入力(例：14:30)")
	cd.Start = read()
	fmt.Print("終了時間を入力:")
	cd.End = read()
	fmt.Print("ZoomURLを入力:")
	cd.Url = read()

	return
}

func ExStartZoom(classes []classData) {
	now := time.Now()
	fmt.Println(now.Weekday())

	fmt.Println("class start.")
}

func StartZoom() {
	filename := "classes.json"
	bytes, err := ioutil.ReadFile(filename)	//json読み込み
	if err != nil {
		log.Fatal(err)
	}
	var classes []classData
	if len(bytes) != 0 {
		if err := json.Unmarshal(bytes, &classes); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("\n--- StartZoom (made by Riku Tsuzuki) --- ")

	flg := 0
	for flg == 0 {
		switch InputNum(sc, "\n0: 終了, 1: 授業登録, 2: 授業開始, 3: 授業リスト") {
		case 0:
			fmt.Println("終了します.")
			flg = 1
		case 1:
			fmt.Println("授業を登録します.")
			classes = append(classes, RegisterClass())
			SaveClass(classes, filename)
		case 2:
			fmt.Println("Zoomを開きます.")
			ExStartZoom(classes)
		case 3:
			fmt.Println("\n登録されている授業を表示します.")
			fmt.Print("\n")
			if len(classes) == 0 {
				fmt.Println("no class.")
			} else {
				for i, v := range classes {
					fmt.Println(i+1, ":", v.Name)
					fmt.Println("", v.Day, v.Start, "~", v.End)
					fmt.Println("", v.Url)
				}
			}
		}
	}
}
