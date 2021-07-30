// udp_sendto.go
// go run udp_sendto.go
// golangによるudpの送信"だけ"するプログラム

package main

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

func random(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

type unmTbl struct {
	uniNum  int
	objType string // "Bus" or "User"
	simNum  int
	pmNum   int
	lon     float64
	lat     float64
}

func main() {
	conn, err := net.Dial("udp4", "localhost:12345")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 初期化
	var ut [5]unmTbl
	for i, _ := range ut {
		ut[i].objType = "User"
		ut[i].uniNum = i
		ut[i].lat = 35.653976
		ut[i].lon = 139.796821
	}

	for {
		fmt.Println("Sending to server")
		for i, _ := range ut {
			ut[i].lat += random(0.5, -0.5) * 0.00001 * 10 * 5
			ut[i].lon += random(0.5, -0.5) * 0.00002 * 10 * 5

			str := ut[i].objType + "," + fmt.Sprint(ut[i].uniNum) + "," + fmt.Sprint(ut[i].lon) + "," + fmt.Sprint(ut[i].lat) + ","

			fmt.Println(str)

			_, err = conn.Write([]byte(str))
			if err != nil {
				panic(err)
			}
			time.Sleep(3 * time.Second) // 1秒休む

		}

	}

}


