// postgreSQLのDBにアクセスして取得した位置情報データを、PrumeMobile上に表示するプログラム
// 初版
// ペアは、server22-1.go

package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"net/url"
	"os"
	"sync"
	"time"

	"database/sql"

	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

const port int = 8910 // DBコンテナが公開しているポート番号

// GetLoc GetLoc
type GetLoc struct {
	ID  int     `json:"id"`
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
	//Address string  `json:"address"`
}

// 構造体の作り方
type unmTbl struct {
	uniName int    // User Name: Example  6ca90e
	objType string // "Bus" or "User"
	simNum  int
	pmNum   int
	lon     float64
	lat     float64
}

var list = make([]unmTbl, 0)                                           // 構造体の動的リスト宣言
var addr = flag.String("addr", "0.0.0.0:8080", "http service address") // テスト

func main() {

	var wg sync.WaitGroup

	_ = websocket.Upgrader{} // use default options
	flag.Parse()
	log.SetFlags(0)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo2"}

	c, _, err1 := websocket.DefaultDialer.Dial(u.String(), nil)
	if err1 != nil {
		log.Fatal("dial:", err1)
	}

	// db: データベースに接続するためのハンドラ
	var db *sql.DB
	// Dbの初期化
	dbParam := fmt.Sprintf("host=localhost port=%d user=postgres password=ca_sim dbname=ca_sim sslmode=disable", port)
	db, err := sql.Open("postgres", dbParam)
	if err != nil {
		fmt.Println("cannot open db")
		os.Exit(1)
	}

	sql := "SELECT id, to_char(time, 'HH24:MI:SS'), user_or_bus,x, y FROM position_log "
	rows, err := db.Query(sql)

	for rows.Next() {
		var id int
		var timestr string
		var user_or_bus string
		var x float64
		var y float64

		time.Sleep(time.Millisecond * 10) // 0.01秒休む

		rows.Scan(&id, &timestr, &user_or_bus, &x, &y)
		//fmt.Println(id, timestr, user_or_bus, x, y)

		uniName := id
		fmt.Printf("%d\n", uniName)

		objType := user_or_bus
		fmt.Printf("%s\n", objType)

		lon := x
		fmt.Printf("%f\n", lon)

		lat := y
		fmt.Printf("%f\n", lat)

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

func movingObject(uniNum int, uniName int, objType string, lon float64, lat float64, wg *sync.WaitGroup, c *websocket.Conn) {

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
