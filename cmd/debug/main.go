package main

import (
	"fmt"
	"io"

	"github.com/oja-bitterlife/yamlui-go/script"
)

func main() {
	code := "(_.set @y (_.add @y 10.5))"

	tn := script.NewTokenizer(code)

	fmt.Println("Tokens:")
	for {
		token, err := tn.Next()
		if err != nil {
			if err == io.EOF {
				break // 正常終了
			}
			panic(err) // 予期しないエラー
		}

		if token == nil {
			break // 安全策
		}

		// ここで token に応じて Value 構造体を作る
		fmt.Printf("Token: [%s]\n", token)
	}

}
