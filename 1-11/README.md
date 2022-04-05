# 1. RedisサーバのPubSub(ブロードキャスト)と、Golangのmapを使って、複数のWebにPrumeMobileの表示をさせるプログラム

<!-- TOC -->

- [1. RedisサーバのPubSubブロードキャストと、Golangのmapを使って、複数のWebにPrumeMobileの表示をさせるプログラム](#1-redis%E3%82%B5%E3%83%BC%E3%83%90%E3%81%AEpubsub%E3%83%96%E3%83%AD%E3%83%BC%E3%83%89%E3%82%AD%E3%83%A3%E3%82%B9%E3%83%88%E3%81%A8golang%E3%81%AEmap%E3%82%92%E4%BD%BF%E3%81%A3%E3%81%A6%E8%A4%87%E6%95%B0%E3%81%AEweb%E3%81%ABprumemobile%E3%81%AE%E8%A1%A8%E7%A4%BA%E3%82%92%E3%81%95%E3%81%9B%E3%82%8B%E3%83%97%E3%83%AD%E3%82%B0%E3%83%A9%E3%83%A0)
- [2. RedisプログラムのインストールとRedigoのインストール](#2-redis%E3%83%97%E3%83%AD%E3%82%B0%E3%83%A9%E3%83%A0%E3%81%AE%E3%82%A4%E3%83%B3%E3%82%B9%E3%83%88%E3%83%BC%E3%83%AB%E3%81%A8redigo%E3%81%AE%E3%82%A4%E3%83%B3%E3%82%B9%E3%83%88%E3%83%BC%E3%83%AB)
- [3. プログラムの稼動方法](#3-%E3%83%97%E3%83%AD%E3%82%B0%E3%83%A9%E3%83%A0%E3%81%AE%E7%A8%BC%E5%8B%95%E6%96%B9%E6%B3%95)

<!-- /TOC -->

# 2. RedisプログラムのインストールとRedigoのインストール

https://wp.kobore.net/%e6%b1%9f%e7%ab%af%e3%81%95%e3%82%93%e3%81%ae%e6%8a%80%e8%a1%93%e3%83%a1%e3%83%a2/post-5817/ を参照のこと

# 3. プログラムの稼動方法
- "go run publisher.go"を起動
- "go run main.go"を起動
- http://localhost:5000 で起動

