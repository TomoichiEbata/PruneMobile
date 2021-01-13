// プログラムの動作チェックの為に、受信データにノイズを含める
// ベースはudp_server3.go
// サーバはserver22.go

package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// GetLoc GetLoc
type GetLoc struct {
	ID  int     `json:"id"`
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
	//Address string  `json:"address"`
}

// 構造体の作り方
type unmTbl struct {
	uniNum  int
	objType string // "Bus" or "User"
	simNum  int
	pmNum   int
	lon     float64
	lat     float64
}

var list = make([]unmTbl, 0) // 構造体の動的リスト宣言

var addr = flag.String("addr", "0.0.0.0:8080", "http service address") // テスト

// ノイズテスト用
func random(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func main() {
	// UDPのポート番号指定
	addr, _ := net.ResolveUDPAddr("udp", ":12345")
	sock, _ := net.ListenUDP("udp", addr)

	ut := unmTbl{} // テーブル先頭の作成 (何も入れない)
	list = append(list, ut)

	// list = make([]unmTbl, 0) // 構造体の動的リスト宣言

	var wg sync.WaitGroup

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

		// ノイズ混入
		lon += random(0.5, -0.5) * 0.00002 * 10 * 5

		lat, _ := strconv.ParseFloat(slice[3], 64)
		fmt.Printf("%f\n", lat)

		// ノイズ混入
		lat += random(0.5, -0.5) * 0.00001 * 10 * 5

		// 新しいオブジェクトかどうかを確認する// リスト分、ループする
		//for i, _ := range list {
		flag := 0

		for i := range list {
			if i != 0 && list[i].objType == objType && list[i].simNum == num {
				list[i].lon = lon // 新しい経度情報の更新
				list[i].lat = lat // 新しい緯度情報の更新

				flag = 1
				break
			}
		}

		uniNum := len(list)

		if flag == 0 {
			wg.Add(1) // goルーチンを実行する関数分だけAddする
			go movingObject(uniNum, objType, num, lon, lat, &wg)
		}

	}
	wg.Wait() // goルーチンで実行される関数が終了するまで待つ。

}

func movingObject(uniNum int, objType string, num int, lon float64, lat float64, wg *sync.WaitGroup) {

	fmt.Printf("start movingObject\n")

	defer wg.Done() // WaitGroupを最後に完了しないといけない。

	//var upgrader = websocket.Upgrader{} // use default options
	_ = websocket.Upgrader{} // use default options

	// rand.Seed(time.Now().UnixNano())

	flag.Parse()
	log.SetFlags(0)
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo2"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// リストを作る前にテストをする
	fmt.Printf("%s\n", objType)
	fmt.Printf("%d\n", num)
	fmt.Printf("%f\n", lon)
	fmt.Printf("%f\n", lat)

	ut := unmTbl{} // 構造体変数の初期化
	ut.uniNum = uniNum
	ut.objType = objType
	ut.simNum = num
	ut.lat = lat
	ut.lon = lon

	gl := new(GetLoc)
	gl.ID = 0
	gl.Lat = ut.lat
	gl.Lng = ut.lon

	err = c.WriteJSON(gl) // PruneMobile登録用送信
	if err != nil {
		log.Println("write:", err)
	}

	gl2 := new(GetLoc) // PruneMobile登録確認用受信
	err = c.ReadJSON(gl2)

	ut.pmNum = gl2.ID // PrumeMobileから提供される番号

	fmt.Printf("ut.objType=%v\n", ut.objType)
	list = append(list, ut) // 構造体をリストに動的追加

	for {
		gl.ID = gl2.ID
		gl.Lat = list[uniNum].lat
		gl.Lng = list[uniNum].lon

		// 座標の送信
		err = c.WriteJSON(gl)
		if err != nil {
			log.Println("write:", err)
		}

		// 応答受信
		gl3 := new(GetLoc)
		err = c.ReadJSON(gl3)

		time.Sleep(1 * time.Second) // 1秒休む

	}

}
