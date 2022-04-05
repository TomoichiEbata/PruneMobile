package main

import "fmt"

func main() {

	// 2次元マップ
	m := map[string]int{}
	//m := make(map[string]int)

	m["apple"] = 150
	m["banana"] = 200
	m["lemon"] = 300

	fmt.Println("0:", m["apple"])

	// 多次元mapは、2次元マップのように簡単には使えない
	// 比較的面倒のないやり方は、以下のように構造体を使う方法(らしい)
	// https://qiita.com/ruiu/items/476f65e7cec07fd3d4d7

	type key struct {
		id  int
		att string
	}

	// 配列宣言
	m1 := make(map[key]int)

	// レコード追加
	m1[key{0, "BUS"}] = 1
	m1[key{0, "PERSON"}] = 1

	m1[key{1, "BUS"}] = 2
	m1[key{1, "PERSON"}] = 11

	// レコードの一つを表示
	fmt.Println("1:", m1[key{1, "BUS"}])
	// レコードの全部を表示
	fmt.Println("2:", m1)

	// 変数を使う方法
	id0 := 1
	att0 := "PERSON"
	fmt.Println("3:", m1[key{id0, att0}])

	// レコードの削除
	delete(m1, key{1, "BUS"})
	fmt.Println("4:", m1)

	for i := range m1 {
		fmt.Println("5:", i)
	}

	// 変数を使って、キーの存在を確認する
	id0 = 1
	att0 = "PERSON"
	value, isThere := m1[key{id0, att0}]
	fmt.Println("6:", value, isThere)

	id0 = 2
	att0 = "PERSON"

	value, isThere = m1[key{id0, att0}]
	fmt.Println("7:", value, isThere)

	// レコード追加の例

	id0 = 3
	att0 = "BUS"
	new_id0 := 20

	fmt.Println("8:", m1)

	_, isThere = m1[key{id0, att0}]
	if !isThere {
		m1[key{id0, att0}] = new_id0
	}

	fmt.Println("9:", m1)

}
