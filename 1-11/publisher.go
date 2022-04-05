package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/gomodule/redigo/redis"
)

type GetLoc struct {
	ID    int     `json:"id"`
	Lat   float64 `json:"lat"`
	Lng   float64 `json:"lng"`
	TYPE  string  `json:"type"` // "PERSON","BUS","CONTROL
	POPUP int     `json:"popup"`
	//Address string  `json:"address"`
}

func main() {
	// 接続
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	count := 0

	var gl GetLoc

	// 100回ループする
	for i := 0; i < 100; i++ {

		if i < 90 {
			gl.TYPE = "BUS"
			gl.Lat = 35.654543 + (rand.Float64()-0.5)*0.00001*100
			gl.Lng = 139.795534 + (rand.Float64()-0.5)*0.00002*100
			gl.ID = rand.Int() % 5
		} else {
			gl.TYPE = "BUS"
			gl.Lat = 181.0
			gl.Lng = 181.0
			gl.ID = (i - 90)
			fmt.Println("gl.ID:", gl.ID)
		}

		// パブリッシュ
		//hello_string := "hello" + strconv.Itoa(count)
		// r, err := redis.Int(conn.Do("PUBLISH", "channel_1", "hello"))
		//r, err := redis.Int(conn.Do("PUBLISH", "channel_1", hello_string))
		// r, err := redis.Int(conn.Do("PUBLISH", "channel_1", gl))

		json_gl, _ := json.Marshal(gl)

		//		r, err := redis.Int(conn.Do("PUBLISH", "channel_1", gl))
		r, err := redis.Int(conn.Do("PUBLISH", "channel_1", json_gl))
		if err != nil {
			panic(err)
		}
		fmt.Println(r)

		time.Sleep(time.Millisecond * 1000)

		count++
	}

}
