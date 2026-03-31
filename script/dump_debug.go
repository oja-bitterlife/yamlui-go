//go:build debug

package script

import (
	"encoding/json"
	"fmt"
)

// **********************************************************************
// デバッグ用のダンプ関数
// ==================================================
// JSON化して出力する関数。普段はこれ
func (v Value) Dump() {
	jsonData, err := json.Marshal(v)
	if err != nil {
		fmt.Printf("Failed to marshal Value: %v\n", err)
		return
	}
	fmt.Printf("%s\n", jsonData)
}
