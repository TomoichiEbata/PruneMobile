/*
// server22.go ペアはclient9.go

// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// 使い方
// go run server13.go      (適当なシェルから)
// http://localhost:8080  (ブラウザ起動)
*/

package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"sync"

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

var mutex sync.Mutex

func echo2(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil) // cはサーバのコネクション
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	//mutex := new(sync.Mutex)

	for {
		//mt, message, err := c.ReadMessage() // クライアントからのメッセージの受信(mtはクライアント識別子)
		//_, _, err := c.ReadMessage() // クライアントからのメッセージの受信(mtはクライアント識別子)

		//mutex.Lock()  // ここに置くとデッドロックしてしまう

		gl := new(GetLoc)

		err := c.ReadJSON(&gl) // クライアントからのメッセージの受信

		mutex.Lock()

		// 原因不明の対処処理
		if gl.ID == 0 && gl.Lat < 0.01 && gl.Lng < 0.01 {
			mutex.Unlock()
			break
		} else if gl.ID < -1 { // 受理できないメッセージとして返信する
			//条件分岐 (変なIDが付与されているメッセージは潰す)
			//if (gl.ID > maxid) || (gl.ID < -1) { // 受理できないメッセージとして返信する

			gl.ID = -1
			gl.Lat = -999
			gl.Lng = -999
			err2 := c.WriteJSON(gl)
			if err2 != nil {
				log.Println("write1:", err2)
				mutex.Unlock()
				break
			}
		} else { // それ以外は転送する
			/*
				log.Printf("echo2 after c.WriteJSON(gl) ID:%d", gl.ID)
				log.Printf("echo2 after c.WriteJSON(gl) Lat:%f", gl.Lat)
				log.Printf("echo2 after c.WriteJSON(gl) Lng:%f", gl.Lng)
			*/
			if err != nil {
				log.Println("read:", err)
				mutex.Unlock()
				break
			}
			//fmt.Printf("echo2 before chan2_1 <- *gl\n")
			chan2_1 <- *gl
			//fmt.Printf("echo2 after chan2_1 <- *gl\n")

			//で、ここで受けとる
			//gl2 := new(GetLoc)
			//fmt.Printf("echo2 before gl2 := <-chan2_1\n")
			gl2 := <-chan2_1
			maxid = gl2.ID // ID最大値の更新
			/*
				log.Printf("echo2 after gl2 := <-chan2_1 ID:%d", gl2.ID)
				log.Printf("echo2 after gl2 := <-chan2_1 Lat:%f", gl2.Lat)
				log.Printf("echo2 after gl2 := <-chan2_1 Lng:%f", gl2.Lng)
			*/
			//fmt.Printf("echo2 before err2 := c.WriteJSON(gl2)\n")
			err2 := c.WriteJSON(gl2)
			//fmt.Printf("echo2 after err2 := c.WriteJSON(gl2)\n")
			if err2 != nil {
				log.Println("write2:", err2)
				mutex.Unlock()
				break
			}
			//fmt.Printf("end of echo2\n")

		}

		mutex.Unlock()
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

		//fmt.Printf("echo before gl := <-chan2_1\n")
		gl := <-chan2_1
		//fmt.Printf("echo after gl := <-chan2_1\n")

		//fmt.Printf("echo before err = c.WriteJSON(gl) gl2.id = %d\n", gl.ID)
		//fmt.Printf("echo before err = c.WriteJSON(gl) gl2.lat = %f\n", gl.Lat)
		//fmt.Printf("echo before err = c.WriteJSON(gl) gl2.lng= %f\n", gl.Lng)
		err = c.WriteJSON(gl)
		if err != nil {
			log.Println("WriteJSON1:", err)
		}
		//fmt.Printf("echo after err = c.WriteJSON(gl)\n")

		//fmt.Printf("echo before err = c.RreadJSON(gl)\n")
		gl2 := new(GetLoc)
		err2 := c.ReadJSON(&gl2)
		//fmt.Printf("echo after err = c.ReadJSON(&gl2) gl2.id = %d\n", gl2.ID)
		//fmt.Printf("echo after err = c.ReadJSON(&gl2) gl2.lat = %f\n", gl2.Lat)
		//fmt.Printf("echo after err = c.ReadJSON(&gl2) gl2.lng= %f\n", gl2.Lng)
		if err2 != nil {
			log.Println("ReadJSON:", err2)
		}
		// ここからチャネルで返す
		//fmt.Printf("echo before chan2_1 <- *gl2 gl2.id = %d\n", gl2.ID)
		//fmt.Printf("echo before chan2_1 <- *gl2 gl2.lat = %f\n", gl2.Lat)
		//fmt.Printf("echo before chan2_1 <- *gl2 gl2.lng = %f\n", gl2.Lng)
		chan2_1 <- *gl2
		//fmt.Printf("echo after chan2_1 <- *gl2\n")
		//fmt.Printf("end of echo\n")
	}

}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func smartphone(w http.ResponseWriter, r *http.Request) {
	smartphoneTemplate.Execute(w, "ws://"+r.Host+"/echo2")
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	http.HandleFunc("/echo2", echo2)           // echo関数を登録 (サーバとして必要)
	http.HandleFunc("/echo", echo)             // echo関数を登録 (サーバとして必要)
	http.HandleFunc("/", home)                 // home関数を登録
	http.HandleFunc("/smartphone", smartphone) // smartphone関数を登録
	log.Fatal(http.ListenAndServe(*addr, nil)) // localhost:8080で起動をセット
}

var smartphoneTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<head>
<meta charset="utf-8">
<script> 

function obj(id, lat, lng){
	this.id = id;
	this.lat = lat;
	this.lng = lng;
}

function random(min, max){
	return  Math.random()*(max-min) + min;
}

// var personal_id;

var lat = 35.654543;
var lng = 139.795534;  

window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
	};

	var personal_id = 0;

	
	///// 起動時のボタン
	// disabled属性を削除
	document.getElementById("open").removeAttribute("disabled");
	document.getElementById("open").style.color = "black";

	// disabled属性を設定 (closeボタンを非活性化)
	document.getElementById("close").setAttribute("disabled", true);
	document.getElementById("close").style.color = "White";	

	// disabled属性を設定 (sendボタンを非活性化)
	document.getElementById("send").setAttribute("disabled", true);
	document.getElementById("send").style.color = "White";	


		
	document.getElementById("open").onclick = function(evt) {
		console.log("document.getElementById open");

		// disabled属性を設定 (openボタンを非活性化)
		document.getElementById("open").setAttribute("disabled", true);
		document.getElementById("open").style.color = "White";		

		// disabled属性を削除
		document.getElementById("send").removeAttribute("disabled");
		document.getElementById("send").style.color = "black";	

		// disabled属性を削除
		document.getElementById("close").removeAttribute("disabled");
		document.getElementById("close").style.color = "black";	

		////////// 削除2
		// ws = new WebSocket("{{.}}");
		////////// 削除2終り


		////////// 削除1
		//var send_obj = new obj(0, 35.654543,139.795534);  // 最初は"0"でエントリ

		//console.log("open:send_obj");	
		//console.log(send_obj.id);	
		//console.log(send_obj.lat);
		//console.log(send_obj.lng);		

		//var json_obj = JSON.stringify(send_obj);
		//ws.send(json_obj);
		/////////// 削除1終り


        if (ws) {
            return false;
        }
		
		
		////////// 追加2
		ws = new WebSocket("{{.}}");
		////////// 追加2終り

		
        ws.onopen = function(evt) {
			print("OPEN");
		
			//ws = new WebSocket("{{.}}");
			
			////////// 追加1			
			var send_obj = new obj(0, 35.654543,139.795534);  // 最初は"0"でエントリ

			console.log("open:send_obj");	
			console.log(send_obj.id);	
			console.log(send_obj.lat);
			console.log(send_obj.lng);		

			var json_obj = JSON.stringify(send_obj);
			ws.send(json_obj);
			/////////// 追加1終り

		}
		
        ws.onclose = function(evt) {

			print("CLOSE");
            ws = null;
        }

		ws.onmessage = function(evt) {  // 受信したメッセージはここに飛んでくる
			print("RESPONSE: " + evt.data);  // jsonメッセージの内容を表示
			// データをJSON形式に変更
			var obj = JSON.parse(evt.data);

			personal_id = obj.id; // IDの取得(何回も取る必要はないが)
			console.log("personal_id");			
			console.log(personal_id);

			
			if ((Math.abs(obj.lat) > 90.0) || (Math.abs(obj.lng) > 180.0)){ // 異常な座標が入った場合は、マーカーを消去する
				console.log("before ws.close()");
				ws.close();
				console.log("after ws.close()");
			}
		}
		
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
	
	document.getElementById("send").onclick = function(evt) {

		console.log("document.getElementById send");

		// disabled属性を設定 (openボタンを非活性化)
		document.getElementById("open").setAttribute("disabled", true);
		document.getElementById("open").style.color = "White";	
	
		// disabled属性を削除
		document.getElementById("send").removeAttribute("disabled");
		document.getElementById("send").style.color = "black";	

		// disabled属性を削除
		document.getElementById("close").removeAttribute("disabled");
		document.getElementById("close").style.color = "black";	
	
		if (!ws) {
			console.log("return false send");
			return false;			
		}

		lat += random(0.5, -0.5) * 0.00001 * 10 * 5;
		lng += random(0.5, -0.5) * 0.00002 * 10 * 5

		
		//var send_obj = new obj(personal_id, 35.654543,139.795534);  // idでエントリ
		var send_obj = new obj(personal_id, lat, lng);  // idでエントリ

		console.log("send:send_obj");	
		console.log(send_obj.id);	
		console.log(send_obj.lat);
		console.log(send_obj.lng);		

		var json_obj = JSON.stringify(send_obj);
		ws.send(json_obj);		

		/*
        print("SEND: " + input.value);
        ws.send(input.value);
		return false;
		*/

		return false;	
    };

	document.getElementById("close").onclick = function(evt) {
		console.log(" document.getElementById close");

		// disabled属性を削除
		document.getElementById("open").removeAttribute("disabled");
		document.getElementById("open").style.color = "black";

		// disabled属性を設定 (closeボタンを非活性化)
		document.getElementById("close").setAttribute("disabled", true);
		document.getElementById("close").style.color = "White";	

		// disabled属性を設定 (sendボタンを非活性化)
		document.getElementById("send").setAttribute("disabled", true);
		document.getElementById("send").style.color = "White";			


        if (!ws) {
            return false;
		}
	
		var send_obj = new obj(personal_id, 999.9, 999.9);  // 最初は"0"でエントリ

		console.log("close:send_obj");
		console.log(send_obj.id);		
		console.log(send_obj.lat);
		console.log(send_obj.lng);		

		var json_obj = JSON.stringify(send_obj);
		ws.send(json_obj);

        //ws.close();  // これはws.onmessageの方で実施
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<!-- <p><input id="input" type="text" value="Hello world!"> -->
<button id="send">Send</button>
<button id="close">Close</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8" />
    <title>PruneMobile</title>

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
    // }).setView(new L.LatLng(35.598563, 139.475528), 18); 広袴

    L.tileLayer('http://{s}.tile.osm.org/{z}/{x}/{y}.png', {
        detectRetina: true,
        maxNativeZoom: 18
    }).addTo(map);

    var leafletView = new PruneClusterForLeaflet(1,1);  // (120,20)がデフォルト

	ws.onopen = function (event) {
	}

	var markers = [];

	//var helicopterIcon = L.icon({ iconUrl: 'http://sintef-9012.github.io/PruneCluster/examples/helicopter.png', iconSize: [48, 48] });
    //var airplaneIcon = L.icon({ iconUrl: 'http://sintef-9012.github.io/PruneCluster/examples/airplane.png', iconSize: [48, 48] });

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

			// 参考資料  http://embed.plnkr.co/WmtpkEqSDJFuFeuiYP54/
			//var marker = new PruneCluster.Marker(obj.lat, obj.lng, {
			//	//popup: "Bell 206 " + i,
			//	icon: helicopterIcon
			//});


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

					//leafletView.RemoveMarkers(markers[obj.id]);  // これでは消えてくれません
					// 1つのマーカーを消すのに、面倒でも以下の2行が必要
					var deleteList = markers.splice(i, 1);					
					leafletView.RemoveMarkers(deleteList);

					// 以下失敗例リスト
					//leafletView.RemoveMarkers(markers[i].hashCode);  //これはダメ
					//leafletView.RemoveMarkers(markers[obj.id],'item');
					//leafletView.ProcessView(); // 試しに入れてみる
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
	}, 1000);

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
