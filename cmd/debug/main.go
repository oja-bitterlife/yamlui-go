package main

import (
	"fmt"

	"github.com/oja-bitterlife/yamlui-go/script"
)

func main() {
	// テスト用のS式
	code := "(_.set @y (_.add @y 10.5))"
	fmt.Printf("Input: %s\n", code)

	// パース実行
	ast := script.Parse(code)

	// 結果の表示
	// %#v を使うと構造体の型情報付きで詳細に表示されます
	fmt.Printf("Parsed AST: %#v\n", ast)
}
