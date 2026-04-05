package yamlui_json

import (
	"encoding/json"

	"github.com/oja-bitterlife/yamlui-go/script"
)

// **********************************************************************
// Any型(通常の)JSONをValue型JSONに変換する
// JSON文字列で返す
func AnyJSONToValueJSON(fileData []byte) ([]byte, error) {
	val, err := AnyJSONToValue(fileData)
	if err != nil {
		return nil, err
	}

	return val.ToValueJSON()
}

// Value型で返す。すぐ使うならParseし直さなくていい
func AnyJSONToValue(fileData []byte) (script.Value, error) {
	// まずは普通にjsonのUnmarshal
	var data any
	if err := json.Unmarshal(fileData, &data); err != nil {
		return script.Value{}, script.LogErr("Failed to unmarshal JSON: %v", err)
	}

	// anyなのでValueに変換
	switch data := data.(type) {
	case map[string]any, []any:
		return anyJSONToValue(data)
	default:
		// UI用のデータは基本的にMapかListのはずなので、その他の型はエラーとする
		return script.Value{}, script.LogErr("Unsupported top-level JSON type: %T", data)
	}
}

// anyJSONToValue. any型のJSONデータを再帰的にscript.Valueに変換する
func anyJSONToValue(data any) (script.Value, error) {
	// まずはListかMapかを判定するために、キーの型を確認する
	switch data := data.(type) {
	case map[string]any:
		res := map[string]script.Value{}
		for k, val := range data {
			v, err := anyJSONToValue(val)
			if err != nil {
				return script.Value{}, err
			}
			res[k] = v
		}
		return script.NewLitMap(res), nil
	case []any:
		list := make([]script.Value, len(data))
		for i, val := range data {
			v, err := anyJSONToValue(val)
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
	case bool:
		return script.NewBool(data), nil
	default:
		return script.Value{}, script.LogErr("Unsupported JSON value type: %T", data)
	}
}
