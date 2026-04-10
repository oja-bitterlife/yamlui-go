package yamlui_json

import (
	"fmt"

	"github.com/oja-bitterlife/yamlui-go/script"
)

// **********************************************************************
// JSONがMarshalされたMapをscript.Valueに変換する.呼び出し口
func UIMapToValue(data any) (script.Value, error) {
	switch data := data.(type) {
	case map[string]any, []any:
		// 変換本体
		value, err := uiMapToValue(data)
		if err != nil {
			return script.Value{}, script.LogErr("Failed to convert JSON to Value: %v", err)
		}

		// Scriptのコンパイルもここでやっちゃう
		return CompileScripts(value), nil

	default:
		// UI用のデータは基本的にMapかListのはずなので、その他の型はエラーとする
		return script.Value{}, script.LogErr("Unsupported top-level JSON type: %T", data)
	}
}

// AnyMapを再帰的にscript.Valueに変換する
func uiMapToValue(data any) (script.Value, error) {
	// まずはListかMapかを判定するために、キーの型を確認する
	switch data := data.(type) {
	case map[string]any:
		res := map[string]script.Value{}
		for k, val := range data {
			v, err := uiMapToValue(val)
			if err != nil {
				return script.Value{}, err
			}
			res[k] = v
		}
		return script.NewLitMap(res), nil
	case []any:
		list := make([]script.Value, len(data))
		for i, val := range data {
			v, err := uiMapToValue(val)
			if err != nil {
				return script.Value{}, err
			}
			list[i] = v
		}
		return script.NewLitList(list), nil
	case string:
		return script.NewString(data), nil
	case float64:
		return script.NewNumber(int(data)), nil
	case uint64: // YAML対応
		return script.NewNumber(int(data)), nil
	case int64: // YAML/TOML対応
		return script.NewNumber(int(data)), nil
	case bool:
		return script.NewBool(data), nil
	default:
		return script.Value{}, script.LogErr("Unsupported JSON value type: %T", data)
	}
}

// **********************************************************************
// CompileScripts は Value 型のツリーを走査し、
// "script" キーを持つ Map 内の文字列をコンパイル済みの型に差し替えます
func CompileScripts(v script.Value) script.Value {
	switch v.Type {
	case script.TypeLitMap:
		for k, child := range v.Map {
			if k == "script" && child.Type == script.TypeString {
				// ここで文字列をコンパイル！
				valueAST, err := script.Compile(child.Str)
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
