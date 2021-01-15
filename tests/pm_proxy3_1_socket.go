// simple_csv_1_socket.go

// 2021/01/14
// simple_cvs.go の場合、数万のエージェント全部とwebsocket通信することになるため、
// サーバとの通信用のソケットを1つに限定して、この問題を回避した。
// ただエージェントは、従来通り、全部の数(数万から十万くらい？)作成したがそれでも。動いている。
// golang まじ凄い

// 変更点は、mainルーチンで、エージェント用のソケットを作って、goroutineでエージェント用のスレッド作る時に、それを渡している点
// あと通信は、送信→受信で1セットになるように、ミューテックスロックで競合回避を行った点(まあ通信が混乱するのを回避するため)
// 普通なら、これで相当の実行速度の低下が発生するはずなんだけど、体感的には遅くなかった。
// 現状の問題点は、chromoの方が先に落ちる、ということかな。まあ、数万のオブジェクトを1秒以内に動かされたら、chromoも文句の一つも言いたかろう。
// この問題は、メモリが潤沢に搭載されているPCでなら回避できるような気がするので、当面は放置することにする

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math"
	"net/url"
	"os"
	"strconv"
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
	uniName string // User Name: Example  6ca90e
	objType string // "Bus" or "User"
	simNum  int
	pmNum   int
	lon     float64
	lat     float64
}

var list = make([]unmTbl, 0)                                           // 構造体の動的リスト宣言
var addr = flag.String("addr", "0.0.0.0:8080", "http service address") // テスト

func main() {
	file, err := os.Open("1.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var wg sync.WaitGroup

	reader := csv.NewReader(file)
	var line []string

	// サーバとのコネクションを1つに統一

	//var upgrader = websocket.Upgrader{} // use default options
	_ = websocket.Upgrader{} // use default options

	// rand.Seed(time.Now().UnixNano())

	flag.Parse()
	log.SetFlags(0)
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo2"}
	//log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	for {

		time.Sleep(time.Millisecond * 1) // 0.001秒休む

		line, err = reader.Read()
		if err != nil {
			break
		}

		uniName := line[0]
		//fmt.Printf("%s\n", uniName)

		objType := line[9]
		//fmt.Printf("%s\n", objType)

		lon, _ := strconv.ParseFloat(line[8], 64)
		//fmt.Printf("%f\n", lon)

		lat, _ := strconv.ParseFloat(line[7], 64)
		//fmt.Printf("%f\n", lat)

		// 特定範囲に限定する

		if lon > 138.80830921713505 && lon < 140.88564728775447 && lat > 35.13735109891961 && lat < 36.31199054701465 {

			//if lon > 139.744330 && lon < 139.866586 && lat > 35.574777 && lat < 35.694479 {
			// if lon > 139.7583407156985 && lon < 139.81403350119444 && lat > 35.62835195825786 && lat < 35.66678018870369 {

			flag := 0

			for i := range list {
				if i != 0 && list[i].uniName == uniName { // 同一IDを発見したら
					list[i].lon = lon // 新しい経度情報の更新
					list[i].lat = lat // 新しい緯度情報の更新

					flag = 1
					break
				}
			}

			uniNum := len(list)

			if flag == 0 { // 新しいIDを発見した場合
				wg.Add(1) // goルーチンを実行する関数分だけAddする
				go movingObject(uniNum, uniName, objType, lon, lat, &wg, c)
			}
		}

	}

	// movingObjectに自己破壊メッセージを送信
	// 破壊情報の書き込み中は邪魔させない
	mutex.Lock()
	for i := range list {
		if i != 0 {
			list[i].lon = 999.9 // デタラメな経度情報の更新
			list[i].lat = 999.9 // デタラメな緯度情報の更新
		}
	}
	mutex.Unlock()

	// goルーチンで実行される関数が終了するまで待つ。
	// wg.Wait() // のを止める
	c.Close()
}

var mutex sync.Mutex

func movingObject(uniNum int, uniName string, objType string, lon float64, lat float64, wg *sync.WaitGroup, c *websocket.Conn) {

	fmt.Printf("start movingObject\n")

	defer wg.Done() // WaitGroupを最後に完了しないといけない。

	//defer c.Close()  // 単一通信だからこれが切れると困る

	// リストを作る前にテストをする
	//fmt.Printf("%s\n", objType)
	//fmt.Printf("%d\n", uniNum)
	//fmt.Printf("%f\n", lon)
	//fmt.Printf("%f\n", lat)

	ut := unmTbl{} // 構造体変数の初期化
	ut.uniName = uniName
	ut.objType = objType
	ut.simNum = uniNum
	ut.lat = lat
	ut.lon = lon

	gl := new(GetLoc)
	gl.ID = 0
	gl.Lat = ut.lat
	gl.Lng = ut.lon

	mutex.Lock()           // 送受信時にミューテックスロックしないと
	err := c.WriteJSON(gl) // PruneMobile登録用送信
	if err != nil {
		log.Println("write1:", err)
	}

	gl2 := new(GetLoc) // PruneMobile登録確認用受信
	err = c.ReadJSON(gl2)
	mutex.Unlock()

	ut.pmNum = gl2.ID // PrumeMobileから提供される番号

	//fmt.Printf("ut.objType=%v\n", ut.objType)
	list = append(list, ut) // 構造体をリストに動的追加

	// ここからは更新用のループ
	for {
		time.Sleep(time.Millisecond * 100) // 0.1秒休む

		// 前回との座標に差が認められれば、移動させる
		if math.Abs(list[uniNum].lat-gl.Lat) > 0.000000001 || math.Abs(list[uniNum].lon-gl.Lng) > 0.000000001 {

			fmt.Print("MOVING!\n")
			gl.Lat = list[uniNum].lat
			gl.Lng = list[uniNum].lon
			gl.ID = gl2.ID

			// 座標の送信

			mutex.Lock()
			err = c.WriteJSON(gl)
			if err != nil {
				log.Println("write2:", err)
			}

			// 応答受信
			gl3 := new(GetLoc)
			err = c.ReadJSON(gl3)
			mutex.Unlock()

			// 異常値によって、上記でブラウザのオブジェクトを消滅させ、さらに、ここでmovingObjectスレッドも消滅させる
			if list[uniNum].lat > 999.0 || list[uniNum].lon > 999.0 {

				println("stop movingObject！")
				return
			}

		}

	}

}
