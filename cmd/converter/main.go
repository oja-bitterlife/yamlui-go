package main

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
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

	// YAMLの型はJSONより多いので、一度JSONを経由させる
	// anyJson, err := json.Marshal(raw)
	// if err != nil {
	// 	script.LogFatal("failed to marshal to JSON: %v", err)
	// }
	// var uiData any
	// if err = json.Unmarshal(anyJson, &uiData); err != nil {
	// 	script.LogFatal("failed to unmarshal JSON: %v", err)
	// }

	// ソースとなるUIデータのMap からUIで使うValue型JSONに変換
	val, err := yamlui_json.UIMapToValue(raw)
	if err != nil {
		panic(err)
	}

	// ValueをValue型JSONに
	data, err := val.ToValueJSON()
	if err != nil {
		panic(err)
	}

	// 成功！
	fmt.Println(string(data))
}
