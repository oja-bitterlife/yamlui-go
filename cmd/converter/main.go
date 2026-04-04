package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/oja-bitterlife/yamlui-go/convert"
	"github.com/oja-bitterlife/yamlui-go/script"
	"github.com/oja-bitterlife/yamlui-go/yamlui_json"
)

func main() {
	// YAML ファイルを読み込む
	yamlData, _ := os.ReadFile("assets/ui.yaml")

	// YAML->map[string]any -> JSON
	var raw any
	if err := yaml.Unmarshal(yamlData, &raw); err != nil {
		script.LogFatal("failed to unmarshal YAML: %v", err)
	}
	anyJson, err := json.Marshal(raw)
	if err != nil {
		script.LogFatal("failed to marshal to JSON: %v", err)
	}

	// JSON からValue型JSONに変換
	val, err := yamlui_json.AnyJSONToValue(anyJson)
	if err != nil {
		panic(err)
	}

	// script部分をコンパイル
	val = CompileScripts(val)

	// ValueをValue型JSONに
	data, err := val.ToValueJSON()
	if err != nil {
		panic(err)
	}

	// 成功！
	fmt.Println(string(data))
}

// CompileScripts は Value 型のツリーを走査し、
// "script" キーを持つ Map 内の文字列をコンパイル済みの型に差し替えます
func CompileScripts(v script.Value) script.Value {
	switch v.Type {
	case script.TypeLitMap:
		for k, child := range v.Map {
			if k == "script" && child.Type == script.TypeString {
				// ここで文字列をコンパイル！
				valueAST, err := convert.Compile(child.Str)
				if err != nil {
					panic(fmt.Sprintf("failed to compile script: %v", err))
				}
				v.Map[k] = valueAST
			} else {
				// それ以外はさらに深く潜る
				v.Map[k] = CompileScripts(child)
			}
		}
	case script.TypeList, script.TypeLitList:
		for i := range v.List {
			v.List[i] = CompileScripts(v.List[i])
		}
	}
	return v
}
