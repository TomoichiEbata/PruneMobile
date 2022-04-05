package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type GetLoc struct {
	ID    int     `json:"id"`
	Lat   float64 `json:"lat"`
	Lng   float64 `json:"lng"`
	TYPE  string  `json:"type"` // "PERSON","BUS","CONTROL
	POPUP int     `json:"popup"`
	//Address string  `json:"address"`
}

func echo(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("websocket connection err:", err)
		return
	}
	defer conn.Close()

	// redisサーバとの接続
	conn_redis, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		panic(err)
	}
	defer conn_redis.Close()

	// mapの作成

	// map処理を開始する
	type key struct {
		id  int
		att string
	}

	// 配列宣言
	m1 := make(map[key]int)

	psc := redis.PubSubConn{Conn: conn_redis}
	// 購読
	psc.Subscribe("channel_1", "channel_2", "channel_3")
	defer psc.Unsubscribe("channel_1", "channel_2", "channel_3")

	//var gl GetLoc

	var mutex sync.Mutex

	//var gl GetLoc
	//var gl2 GetLoc

	for {
		gl := new(GetLoc)
		gl2 := new(GetLoc)

		switch v := psc.Receive().(type) {
		case redis.Message:
			//fmt.Printf("%s: message: %s\n", v.Channel, v.Data)

			_ = json.Unmarshal(v.Data, &gl)
			//fmt.Println("start gl:", gl)
			fmt.Println()

			mutex.Lock()

			// 変数を使って、キーの存在を確認する
			value, isThere := m1[key{gl.ID, gl.TYPE}]

			fmt.Println("0:value:", value, "isThere:", isThere, "gl.ID:", gl.ID, "gl.TYPE", gl.TYPE)

			if math.Abs(gl.Lat) > 90.0 || math.Abs(gl.Lng) > 180.0 { // ありえない座標が投入されたら
				fmt.Println("enter 1")

				fmt.Println("1:gl:", gl)
				conn.WriteJSON(&gl) // 送って
				conn.ReadJSON(&gl2) // 戻して
				fmt.Println("1:gl2:", gl2)

				delete(m1, key{gl.ID, gl.TYPE}) // レコードを削除して終了する

			} else {
				if !isThere { // もしレコードが存在しなければ(新しいIDであれば)
					fmt.Println("enter 2")

					id := gl.ID
					gl.ID = -1 // 空番号 これでJavaScriptの方に

					fmt.Println("2:gl:", gl)

					conn.WriteJSON(&gl) // 送って
					conn.ReadJSON(&gl2) // 戻してもらって

					fmt.Println("2:gl2:", gl2)

					pm_id := gl2.ID // JavaScriptから与えられたIDで
					//fmt.Println("id:", id, ", pm_id:", pm_id)

					time.Sleep(time.Second * 1)

					m1[key{id, gl.TYPE}] = pm_id // レコードを追加する

				} else { //レコードが存在すれば、その値を使ってアイコンを動かす

					fmt.Println("enter 3")

					gl.ID = value // mapから見つけた値を使って、 (バグはここ、 pm_idが"0"の場合エントリーされてしまう)

					fmt.Println("3:gl:", gl)
					conn.WriteJSON(&gl) // アイコンを動かす
					conn.ReadJSON(&gl2)
					fmt.Println("3:gl2:", gl2)

				}
			}
			// 後で消すこと
			fmt.Println("m1:", m1)
			fmt.Println()

			mutex.Unlock()

		case redis.Subscription:
			fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			return
		}

	}

}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./js")))
	http.HandleFunc("/echo", echo)

	log.Println("server starting...", "http://localhost:5000")
	log.Fatal(http.ListenAndServe("localhost:5000", nil))
}
