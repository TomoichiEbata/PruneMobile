// position_log_sample.go

// position_logから座標情報を読み出すサンプル (要pq)

// go get github.com/lib/pq

package main

import (
	"fmt"
	"os"

	"database/sql"

	_ "github.com/lib/pq"
)

const port int = 8910 // DBコンテナが公開しているポート番号

func trackUser(db *sql.DB, userID int, dateIndex int) {
	// SQLステートメント
	sql := "SELECT id, to_char(time, 'HH24:MI:SS'), x, y, satisfaction FROM position_log "
	sql += " WHERE date_index = $1 AND user_or_bus = 'USER' AND id = $2;"
	prepared, err := db.Prepare(sql)
	rows, err := prepared.Query(dateIndex, userID)

	if err != nil {
		fmt.Printf("error : %v", err)
		os.Exit(1)
	}

	fmt.Printf("time, x, y, satisfaction (userID: %v)\n", userID)
	for rows.Next() {
		var id int
		var timestr string
		var x float64
		var y float64
		var satisfaction float64 // 満足度・現実装では、リクエスト発行直後に変化する
		rows.Scan(&id, &timestr, &x, &y, &satisfaction)
		fmt.Println("-", timestr, x, y, satisfaction)
	}
}

func trackBus(db *sql.DB, busID int, dateIndex int) {
	// SQLステートメント
	sql := "SELECT id, to_char(time, 'HH24:MI:SS'), x, y FROM position_log "
	sql += " WHERE date_index = $1 AND user_or_bus = 'BUS' AND id = $2;"
	prepared, err := db.Prepare(sql)
	rows, err := prepared.Query(dateIndex, busID)

	if err != nil {
		fmt.Printf("error : %v", err)
		os.Exit(1)
	}

	fmt.Printf("time, x, y (busID: %v)\n", busID)
	for rows.Next() {
		var id int
		var timestr string
		var x float64
		var y float64
		// var satisfaction float64  バスの満足度はNULL
		rows.Scan(&id, &timestr, &x, &y)
		fmt.Println("-", timestr, x, y)
	}
}

func main() {
	// db: データベースに接続するためのハンドラ
	var db *sql.DB
	// Dbの初期化
	dbParam := fmt.Sprintf("host=localhost port=%d user=postgres password=ca_sim dbname=ca_sim sslmode=disable", port)
	db, err := sql.Open("postgres", dbParam)
	if err != nil {
		fmt.Println("cannot open db")
		os.Exit(1)
	}
	dateIndex := 1 // 1日目
	busID := 8     // バス番号
	userID := 11   // ユーザ番号
	trackBus(db, busID, dateIndex)
	trackUser(db, userID, dateIndex)
}
