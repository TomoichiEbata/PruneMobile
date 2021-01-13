package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// 構造体の作り方
type unmTbl struct {
	objType string // "Bus" or "User"
	simNum  int
	pmNum   int
}

var list = make([]unmTbl, 0) // 構造体の動的リスト宣言

var addr = flag.String("addr", "0.0.0.0:8080", "http service address") // テスト

func main() {
	// UDPのポート番号指定
	addr, _ := net.ResolveUDPAddr("udp", ":12345")
	sock, _ := net.ListenUDP("udp", addr)

	ut := unmTbl{} // テーブル先頭の作成
	list = append(list, ut)

	// list = make([]unmTbl, 0) // 構造体の動的リスト宣言

	for {
		buf := make([]byte, 1024)
		// rlen, _, err := sock.ReadFromUDP(buf)
		_, _, err := sock.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Println(string(buf[0:rlen]))
		strbuffer := string(buf) // convert read in file to a string
		// fmt.Println(strbuffer)

		slice := strings.Split(strbuffer, ",")
		/*
			for _, str := range slice {
				fmt.Printf("[%s]", str)
			}
		*/

		fmt.Printf("\n")
		objType := slice[0]
		fmt.Printf(objType)

		fmt.Printf("\n")
		num, _ := strconv.Atoi(slice[1])
		fmt.Printf("%d\n", num)

		lon, _ := strconv.ParseFloat(slice[2], 64)
		fmt.Printf("%f\n", lon)

		lat, _ := strconv.ParseFloat(slice[3], 64)
		fmt.Printf("%f\n", lat)

		// 新しいオブジェクトかどうかを確認する// リスト分、ループする
		//for i, _ := range list {
		flag := 0

		fmt.Printf("***************\n")
		fmt.Printf("list size = %d\n", len(list))
		fmt.Printf("***************\n")

		for i := range list {
			if i != 0 && list[i].objType == objType && list[i].simNum == num {
				flag = 1
				break
			}
		}

		uniNum := len(list) + 1
		if flag == 0 {
			movingObject(uniNum, objType, num, lon, lat)
		}

	}
}

func movingObject(uniNum int, objType string, num int, lon float64, lat float64) {

	// リストを作る前にテストをする
	fmt.Printf("%s\n", objType)
	fmt.Printf("%d\n", num)
	fmt.Printf("%f\n", lon)
	fmt.Printf("%f\n", lat)

	ut := unmTbl{} // 構造体変数の初期化
	ut.objType = objType
	ut.simNum = num
	ut.pmNum = 0

	fmt.Printf("ut.objType=%v\n", ut.objType)
	list = append(list, ut) // 構造体をリストに動的追加

	for i := range list {
		fmt.Printf("initEntry i:%v, list[i].objType:%v, list[i].simNum:%v, list[i].pmNum:%v\n",
			i, list[i].objType, list[i].simNum, list[i].pmNum)
	}

	//fmt.Println(list)

}
