<!-- TOC -->

- [1. simple_udpについて](#1-simple_udp%E3%81%AB%E3%81%A4%E3%81%84%E3%81%A6)
- [2. ビッグデータを使ったシミュレーションについて](#2-%E3%83%93%E3%83%83%E3%82%B0%E3%83%87%E3%83%BC%E3%82%BF%E3%82%92%E4%BD%BF%E3%81%A3%E3%81%9F%E3%82%B7%E3%83%9F%E3%83%A5%E3%83%AC%E3%83%BC%E3%82%B7%E3%83%A7%E3%83%B3%E3%81%AB%E3%81%A4%E3%81%84%E3%81%A6)
- [3. ローカルでの通信のペア](#3-%E3%83%AD%E3%83%BC%E3%82%AB%E3%83%AB%E3%81%A7%E3%81%AE%E9%80%9A%E4%BF%A1%E3%81%AE%E3%83%9A%E3%82%A2)
    - [3.1. 基本形](#31-%E5%9F%BA%E6%9C%AC%E5%BD%A2)
    - [3.2. オリジナルのアイコンの稼動例](#32-%E3%82%AA%E3%83%AA%E3%82%B8%E3%83%8A%E3%83%AB%E3%81%AE%E3%82%A2%E3%82%A4%E3%82%B3%E3%83%B3%E3%81%AE%E7%A8%BC%E5%8B%95%E4%BE%8B)

<!-- /TOC -->

# 1. simple_udpについて
simple_udpについては、こちら↓に説明がある(Windowsバージョン等の説明もある)
https://wp.kobore.net/江端さんの技術メモ/post-1959/

ちなみに、golangを使ったUDPの通信については、以下に記述がある(参考)
https://wp.kobore.net/江端さんの技術メモ/post-2051/


# 2. ビッグデータを使ったシミュレーションについて

(2)ビッグデータを使ったシミュレーションについては、以下のようにする(思い出せ>私)
```
>go run pm_proxy3_1_socket.go
>go run server22-1.go
>http://localhost:8080
```

# 3. ローカルでの通信のペア

## 3.1. 基本形
server22.go と client9.go // 江端のローカル環境では、 ~/go_template/tests にある

## 3.2. オリジナルのアイコンの稼動例

server23.go と client10.go // 江端のローカル環境では、 ~/go_template/tests 

**アイコンは、~/go_template/tests/staticの中に入っている。**

ちなみに、golangで作ったサーバから、js, css,をローディングできるように、
tests/static, tests/static の中に必要なファイルを入れてあるが、 https://wp.kobore.net/江端さんの技術メモ/post-5280/ で作成した領域からローディングしようとしてが失敗している(サーバが立ち上がる時間に間にあわない様子)
 




