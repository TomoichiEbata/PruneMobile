// client5.go ペアは server10.go

package main

import (
	"flag"
	"log"
	"net/url"

	"github.com/gorilla/websocket"

	"math/rand"
	"time"
)

// GetLoc GetLoc
type GetLoc struct {
	ID  int     `json:"id"`
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
	//Address string  `json:"address"`
}

var addr = flag.String("addr", "0.0.0.0:8080", "http service address") // テスト
//var addr = flag.String("addr", "localhost:8080", "http service address") // テスト
var upgrader = websocket.Upgrader{} // use default options

func random(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func main() {
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

	gl := GetLoc{
		ID:  0,
		Lat: 35.653976,
		Lng: 139.796821,
	}

	log.Printf("ID:%d", gl.ID)
	log.Printf("Lat:%f", gl.Lat)
	log.Printf("Lng:%f", gl.Lng)
	//err = c.WriteJSON(mt, gl)
	err = c.WriteJSON(gl)
	if err != nil {
		log.Println("write:", err)
	}

	gl2 := new(GetLoc)
	err = c.ReadJSON(gl2)
	log.Printf("after ID:%d", gl2.ID)
	log.Printf("after Lat:%f", gl2.Lat)
	log.Printf("after Lng:%f", gl2.Lng)

	gl.ID = gl2.ID
	for i := 0; i < 20; i++ {
		gl.Lat += random(0.5, -0.5) * 0.00001 * 100
		gl.Lng += random(0.5, -0.5) * 0.00002 * 100
		err = c.WriteJSON(gl)
		if err != nil {
			log.Println("write:", err)
		}
		gl2 := new(GetLoc)
		err = c.ReadJSON(gl2)
		log.Printf("after ID:%d", gl2.ID)
		log.Printf("after Lat:%f", gl2.Lat)
		log.Printf("after Lng:%f", gl2.Lng)

		time.Sleep(1 * time.Second) // 1秒休む
	}

	gl.ID = gl2.ID
	gl.Lat = 999.9
	gl.Lng = 999.9

	err = c.WriteJSON(gl)

	err = c.ReadJSON(gl2)
	log.Printf("after ID:%d", gl2.ID)
	log.Printf("after Lat:%f", gl2.Lat)
	log.Printf("after Lng:%f", gl2.Lng)
}
