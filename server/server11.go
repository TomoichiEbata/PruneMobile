/*
// server11.go ペアはclient5.go

// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// 使い方
// go run server9.go      (適当なシェルから)
// http://localhost:8080  (ブラウザ起動)
*/

package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// GetLoc GetLoc
type GetLoc struct {
	ID  int     `json:"id"`
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
	//Address string  `json:"address"`
}

//var addr = flag.String("addr", "localhost:8080", "http service address")
var addr = flag.String("addr", "0.0.0.0:8080", "http service address") // テスト

var upgrader = websocket.Upgrader{} // use default options

var chan2_1 = make(chan GetLoc)

var maxid = 0

func echo2(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil) // cはサーバのコネクション
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for {
		//mt, message, err := c.ReadMessage() // クライアントからのメッセージの受信(mtはクライアント識別子)
		//_, _, err := c.ReadMessage() // クライアントからのメッセージの受信(mtはクライアント識別子)

		gl := new(GetLoc)

		err := c.ReadJSON(&gl) // クライアントからのメッセージの受信(mtはクライアント識別子)

		//条件分岐 (変なIDが付与されているメッセージは潰す)
		//if (gl.ID > maxid) || (gl.ID < -1) { // 受理できないメッセージとして返信する
		if gl.ID < -1 { // 受理できないメッセージとして返信する
			gl.ID = -1
			gl.Lat = -999
			gl.Lng = -999
			err2 := c.WriteJSON(gl)
			if err2 != nil {
				log.Println("write:", err2)
				break
			}
		} else { // それ以外は転送する
			fmt.Printf("71\n")
			log.Printf("first ID:%d", gl.ID)
			log.Printf("first Lat:%f", gl.Lat)
			log.Printf("first Lng:%f", gl.Lng)

			if err != nil {
				log.Println("read:", err)
				break
			}
			fmt.Printf("80\n")
			chan2_1 <- *gl
			fmt.Printf("82\n")
			//で、ここで受けとる
			//gl2 := new(GetLoc)
			gl2 := <-chan2_1
			fmt.Printf("86\n")

			maxid = gl2.ID // ID最大値の更新
			log.Printf("final ID:%d", gl2.ID)
			log.Printf("final Lat:%f", gl2.Lat)
			log.Printf("final Lng:%f", gl2.Lng)

			err2 := c.WriteJSON(gl2)
			fmt.Printf("94\n")
			if err2 != nil {
				log.Println("write:", err2)
				break
			}
			fmt.Printf("99\n")
		}
	}
}

func echo(w http.ResponseWriter, r *http.Request) {

	c, err := upgrader.Upgrade(w, r, nil) // cはサーバのコネクション
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	/*	ここでロックして待つ */

	for {
		fmt.Printf("113 before gl = <-chan2_1\n")

		gl := <-chan2_1

		fmt.Printf("117 after gl = <-chan2_1\n")

		//log.Printf("recv_serv: %s", gl)
		//err = c.WriteJSON(mt, gl)
		err = c.WriteJSON(gl)
		if err != nil {
			log.Println("WriteJSON:", err)
		}

		fmt.Printf("126\n")
		/*
			_, message, err2 := c.ReadMessage() // クライアントからのメッセージの受信(mtはクライアント識別子)
			if err2 != nil {
				log.Println("read:", err2)
				break
			}
			log.Printf("recv_serv: %s", message) // 受信したメッセージの表示
		*/

		gl2 := new(GetLoc)
		err2 := c.ReadJSON(&gl2)
		fmt.Printf("138\n")
		fmt.Printf("after gl2.id = %d\n", gl2.ID)
		fmt.Printf("after gl2.lat = %f\n", gl2.Lat)
		fmt.Printf("after gl2.lng= %f\n", gl2.Lng)
		if err2 != nil {
			log.Println("ReadJSON:", err2)
		}
		fmt.Printf("145\n")
		// ここからチャネルで返す
		chan2_1 <- *gl2
		fmt.Printf("148\n")
	}

}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	http.HandleFunc("/echo2", echo2)           // echo関数を登録 (サーバとして必要)
	http.HandleFunc("/echo", echo)             // echo関数を登録 (サーバとして必要)
	http.HandleFunc("/", home)                 // home関数を登録
	log.Fatal(http.ListenAndServe(*addr, nil)) // localhost:8080で起動をセット
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8" />
    <title>PruneCluster - Realworld 50k</title>

	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.0.0-beta.2.rc.2/leaflet.css"/>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.0.0-beta.2.rc.2/leaflet.js"></script>

	<script src="http://kobore.net/PruneCluster.js"></script>           <!-- これ、いずれローカルホストから取れるように換える -->
	<link rel="stylesheet" href="http://kobore.net/examples.css"/>      <!-- これも、いずれローカルホストから取れるように換える -->

	<!-- goのテンプレートのローカルって、どこになるんだろう？ -->

</head>
<body>
<div id="map"></div>

<script>

	ws = new WebSocket("{{.}}"); // websocketの確立

	/*
	var print = function(message) {
		var d = document.createElement("div");
		d.textContent = message;
		output.appendChild(d);
	};
	*/

	// 引数にはミリ秒を指定。（例：5秒の場合は5000）
	function sleep(a){
  		var dt1 = new Date().getTime();
  		var dt2 = new Date().getTime();
  		while (dt2 < dt1 + a){
			dt2 = new Date().getTime();
		}
  		return;
	}

    var map = L.map("map", {
        attributionControl: false,
        zoomControl: false
    }).setView(new L.LatLng(35.654543, 139.795534), 18);

    L.tileLayer('http://{s}.tile.osm.org/{z}/{x}/{y}.png', {
        detectRetina: true,
        maxNativeZoom: 18
    }).addTo(map);

    var leafletView = new PruneClusterForLeaflet(1,1);  // (120,20)がデフォルト

	ws.onopen = function (event) {
	}

	var markers = [];

	// 受信すると、勝手にここに飛んでくる
	ws.onmessage = function (event) {
		// データをJSON形式に変更
		var obj = JSON.parse(event.data);

		console.log("233");	
		console.log(obj.id);
		console.log(obj.lat);						
		console.log(obj.lng);	


		if (obj.id == 0){  // idが未登録の場合
			console.log("obj.id == 0")
			// データをマーカーとして登録
			var marker = new PruneCluster.Marker(obj.lat, obj.lng);
			console.log(marker.hashCode);		
			markers.push(marker);
	
			leafletView.RegisterMarker(marker);
	
			console.log(markers);
			console.log(markers.length)

			obj.id = marker.hashCode;
			//ws.send(marker.hashCode); // テキスト送信
			var json_obj = JSON.stringify(obj);
			ws.send(json_obj);			
		} else if ((Math.abs(obj.lat) > 90.0) || (Math.abs(obj.lng) > 180.0)){ // 異常な座標が入った場合は、マーカーを消去する
			console.log("Math.abs(obj.lat) > 180.0)")
			for (let i = 0; i < markers.length; ++i) {
				if (obj.id == markers[i].hashCode){
					console.log(i)
					console.log(obj.id)										
					console.log("obj.id == markers[i].hashCode")
					leafletView.RemoveMarkers(markers[obj.id]);
					//leafletView.RemoveMarkers(markers[i-1]);
					//leafletView.RemoveMarkers(markers);					
					break;
				}
			}
			obj.lat = 91.0;
			obj.lng = 181.0;
			var json_obj = JSON.stringify(obj);
			ws.send(json_obj);				
		} else {
			// 位置情報更新
			console.log("else")
			for (let i = 0; i < markers.length; ++i) {
				if (obj.id == markers[i].hashCode){
					var ll = markers[i].position;
					ll.lat = obj.lat;
					ll.lng = obj.lng;
					break;
				}
			}
			var json_obj = JSON.stringify(obj);
			ws.send(json_obj);	
		}
	}

	// 位置情報の更新
    window.setInterval(function () {
        leafletView.ProcessView();  // 変更が行われたときに呼び出されれなければならない
	}, 500);

	// サーバを止めると、ここに飛んでくる
	ws.onclose = function(event) {
		//print("CLOSE");
		ws = null;
	}


    map.addLayer(leafletView);
</script>



</body>
</html>
`))
