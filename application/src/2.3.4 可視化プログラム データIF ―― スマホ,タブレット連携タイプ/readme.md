# Amazon Lightsail を使った、スマホの現在位置の表示方法

PruneMobileは、シミュレータ等で計算した位置情報を、ブラウザで表示することを目的としたものですが、これを、現実のスマホの位置の検知にも使えるようにしました(要するに「ココセコム」としても使える、ということです)

これを実現する為には、インターネット上に(クラウド)サーバを置かなければなりません。AWSのVPS(仮想専用サーバー)が思いつきますが、AWSのEC2は運用が面倒な上に使用料が高価です。そこで「月額 500 円で使えるAWSクラウドのVPS」を使う方法について記載しておきます。

- Amazon Lightsail の立ち上げ方法については、<a href="https://wp.kobore.net/江端さんの技術メモ/post-1513/" target="_blank">こちら</a> を参考にして下さい。

- ここでは、"sea-anemone.tech"という架空のドメインを例として使っていますが、外部(例えば「お名前.com」)でドメインを得た場合は、その名前に置き換えて読んで下さい。


- 公開鍵の取得方法については、<a href="https://wp.kobore.net/江端さんの技術メモ/post-1550/" target="_blank">こちら</a>を参考にして下さい(ここに記載されている、"go_template/server_test"は、"PruneMobile\vps_server"と読み換えて下さい)


## Step 1 サーバの起動

Amazon Lightsailのシェルから適当なシェルを立ち上げて
```
$ cd PruneMobile\vps_server    (江端の環境では、~/go_template/server_test/ )
$ go run serverXX.go (Xは数字)
```

と起動して下さい。

## Step 2 地図画面(マーカ表示画面)の起動
Chromoブラウザ(他のブラウザのことは知らん)から、
```
https://sea-anemone.tech:8080/
```
と入力して下さい。現在は、東京のある地域が表示されますが、serverXX.go の中に記載れている、位置情報、35.60000, 139.60000 を片っぱしから、任意の位置(自宅の位置等)に変更することで、自宅付近での実証実験ができます。
自宅の情報は、GoogleMAPから取得できます。

## Step 3 移動オブジェクト(マーカの対象)の起動
スマホのブラウザから、
```
https://sea-anemone.tech:8080/smartphone
```
として、[open]ボタンを押して下さい。スマホで位置測位が開始されます(この際、位置情報を提供して良いか、と聞きあれることがありますので、"OK"として下さい)。
[close]ボタンを押下すると地図画面からマーカが消えます。

## 現時点で確認している問題点で、いずれ直すもの

- ローカルのjs(javascript)のローディングに失敗した為、江端のプライベートサーバ(kobore.net)からローディングしている。PruneMobile\vps_server\serverXX.goの以下を参照
```
	<script src="http://kobore.net/PruneCluster.js"></script>
	<link rel="stylesheet" href="http://kobore.net/examples.css"/>
```
- 動作中にwebsocketが切断してしまった時(スマホの閉じる、別のブラウザ画面を立ち上げた時)、オブジェクトが放置されて、システム全体が動かなくなる


## ローカルでの通信のペア
server22.go と client9.go 

## デモシミュレーションのペア
server22-1.go と pm_proxy3_1_socket.go

