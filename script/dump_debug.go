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
	fmt.Printf("%s\n", v.ToJson())
}

// ==================================================
// JSON形式で出力するための構造体と関数
type debugJsonTree struct {
	Type  string
	Value any
}

// JSON形式で出力する関数。構造体に変換してからJSON化
// ----------------------------------------
func (v Value) ToJson() string {
	root := v.toJson()
	str, _ := json.Marshal(root)
	return string(str)
}

func (v Value) toJson() debugJsonTree {
	switch v.Type {
	case TypeNumber:
		return debugJsonTree{Type: "Number", Value: v.Num}
	case TypeString:
		return debugJsonTree{Type: "String", Value: v.Str}
	case TypeProperty:
		return debugJsonTree{Type: "Property", Value: v.Str}
	case TypeList:
		list := make([]debugJsonTree, len(v.List))
		for i, item := range v.List {
			list[i] = item.toJson()
		}
		return debugJsonTree{Type: "List", Value: list}
	default:
		return debugJsonTree{Type: "Unknown", Value: v.Type}
	}
}
