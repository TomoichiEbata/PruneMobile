// client9.go ペアは server22.go

package main

import (
	"flag"
	"log"
	"math/rand"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// GetLoc GetLoc
type GetLoc struct {
	ID   int     `json:"id"`
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
	TYPE string  `json:"type"` // "PERSON", "BUS", "CONTROL"
	//Address string  `json:"address"`
}

var addr = flag.String("addr", "0.0.0.0:8080", "http service address") // テスト

func random(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func main() {

	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ { // 3台分作る
		wg.Add(1) // goルーチンを実行する関数分だけAddする
		go bus(i, &wg)
		time.Sleep(1 * time.Second)
	}

	for i := 1; i < 100; i++ { // 100人分作る
		wg.Add(1) // goルーチンを実行する関数分だけAddする
		go passenger(i, &wg)
		time.Sleep(1 * time.Second)
	}

	// goルーチンで実行される関数が終了するまで待つ。
	wg.Wait()

}

func bus(count int, wg *sync.WaitGroup) {

	defer wg.Done() // WaitGroupを最後に完了しないといけない。

	//var upgrader = websocket.Upgrader{} // use default options
	_ = websocket.Upgrader{} // use default options

	rand.Seed(time.Now().UnixNano())

	flag.Parse()
	log.SetFlags(0)
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo2"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// 広袴
	//gl := GetLoc{
	//	ID:  0,
	//	Lat: 35.598563,
	//	Lng: 139.475528,
	//}

	// 豊洲
	gl := GetLoc{
		ID:   0,
		Lat:  35.653976,
		Lng:  139.796821,
		TYPE: "BUS",
	}

	err = c.WriteJSON(gl)
	if err != nil {
		log.Println("write:", err)
	}

	gl2 := new(GetLoc)
	err = c.ReadJSON(gl2)

	gl.ID = gl2.ID
	for i := 0; i < 20; i++ {
		gl.Lat += random(0.5, -0.5) * 0.00001 * 10 * 5
		gl.Lng += random(0.5, -0.5) * 0.00002 * 10 * 5

		err = c.WriteJSON(gl)
		if err != nil {
			log.Println("write:", err)
		}
		gl2 := new(GetLoc)
		err = c.ReadJSON(gl2)

		time.Sleep(1 * time.Second) // 1秒休む
	}

	gl.ID = gl2.ID
	gl.Lat = 999.9
	gl.Lng = 999.9

	err = c.WriteJSON(gl)
	err = c.ReadJSON(gl2)
}

func passenger(count int, wg *sync.WaitGroup) {

	//var addr = flag.String("addr", "0.0.0.0:8080", "http service address") // テスト
	//var addr = flag.String("addr", "localhost:8080", "http service address") // テスト

	defer wg.Done() // WaitGroupを最後に完了しないといけない。

	//var upgrader = websocket.Upgrader{} // use default options
	_ = websocket.Upgrader{} // use default options

	rand.Seed(time.Now().UnixNano())

	flag.Parse()
	log.SetFlags(0)
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo2"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// 広袴
	//gl := GetLoc{
	//	ID:  0,
	//	Lat: 35.598563,
	//	Lng: 139.475528,
	//}

	// 豊洲
	gl := GetLoc{
		ID:   0,
		Lat:  35.653976,
		Lng:  139.796821,
		TYPE: "PERSON",
	}

	log.Printf("count:%d before 1 ID:%d", count, gl.ID)
	log.Printf("count:%d before 1 Lat:%f", count, gl.Lat)
	log.Printf("count:%d before 1 Lng:%f", count, gl.Lng)
	log.Printf("count:%d before 1 TYPE:%s", count, gl.TYPE)
	//err = c.WriteJSON(mt, gl)
	err = c.WriteJSON(gl)
	if err != nil {
		log.Println("write:", err)
	}

	gl2 := new(GetLoc)
	err = c.ReadJSON(gl2)
	log.Printf("count:%d after1 ID:%d", count, gl2.ID)
	log.Printf("count:%d after1 Lat:%f", count, gl2.Lat)
	log.Printf("count:%d after1 Lng:%f", count, gl2.Lng)
	log.Printf("count:%d before 1 TYPE:%s", count, gl2.TYPE)

	gl.ID = gl2.ID
	for i := 0; i < 20; i++ {
		gl.Lat += random(0.5, -0.5) * 0.00001 * 10 * 5
		gl.Lng += random(0.5, -0.5) * 0.00002 * 10 * 5

		log.Printf("count:%d-%d before 2 ID:%d", count, i, gl.ID)
		log.Printf("count:%d-%d before 2 Lat:%f", count, i, gl.Lat)
		log.Printf("count:%d-%d before 2 Lng:%f", count, i, gl.Lng)
		log.Printf("count:%d-%d before 2 Lng:%s", count, i, gl.TYPE)

		err = c.WriteJSON(gl)
		if err != nil {
			log.Println("write:", err)
		}
		gl2 := new(GetLoc)
		err = c.ReadJSON(gl2)
		log.Printf("count:%d-%d after 2 ID:%d", count, i, gl2.ID)
		log.Printf("count:%d-%d after 2 Lat:%f", count, i, gl2.Lat)
		log.Printf("count:%d-%d after 2 Lng:%f", count, i, gl2.Lng)
		log.Printf("count:%d-%d after 2 Lng:%s", count, i, gl2.TYPE)

		time.Sleep(1 * time.Second) // 1秒休む
	}

	gl.ID = gl2.ID
	gl.Lat = 999.9
	gl.Lng = 999.9

	log.Printf("count:%d before 3 ID:%d", count, gl.ID)
	log.Printf("count:%d before 3 Lat:%f", count, gl.Lat)
	log.Printf("count:%d before 3 Lng:%f", count, gl.Lng)
	log.Printf("count:%d before 3 TYPE:%s", count, gl.TYPE)

	err = c.WriteJSON(gl)

	err = c.ReadJSON(gl2)
	log.Printf("count:%d after3 ID:%d", count, gl2.ID)
	log.Printf("count:%d after3 Lat:%f", count, gl2.Lat)
	log.Printf("count:%d after3 Lng:%f", count, gl2.Lng)
	log.Printf("count:%d after 3 TYPE:%s", count, gl.TYPE)
}
