//go:build debug

package script

import (
	"encoding/json"
	"fmt"
)

// **********************************************************************
// デバッグ用のダンプ関数
// ==================================================
// 文字列化して出力する関数。普段はこれ
func (v Value) Dump() {
	fmt.Printf("%s\n", v.ToJSON())
}

// ==================================================
// JSON形式で出力するための構造体と関数
type debugJSONTree struct {
	Type  string
	Value any
}

// JSON形式で出力する関数。構造体に変換してからJSON化
// ----------------------------------------
func (v Value) ToJSON() string {
	root := v.toJSON()
	str, _ := json.Marshal(root)
	return string(str)
}

func (v Value) toJSON() debugJSONTree {
	if v.IsList() {
		list := make([]debugJSONTree, len(v.List))
		for i, item := range v.List {
			list[i] = item.toJSON()
		}
		return debugJSONTree{Type: v.Type.ToStr(), Value: list}
	} else {
		return debugJSONTree{Type: v.Type.ToStr(), Value: v.ToStr()}
	}
}
