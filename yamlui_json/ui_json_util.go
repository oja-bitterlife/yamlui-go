package yamlui_json

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/oja-bitterlife/yamlui-go/script"
)

func AnyJSONToValueJSON(fileData []byte) ([]byte, error) {
	// まずは普通にjsonのUnmarshal
	var data any
	if err := json.Unmarshal(fileData, &data); err != nil {
		return []byte{}, errors.New("Failed to parse JSON: " + err.Error())
	}

	// anyなのでValueに変換
	switch data := data.(type) {
	case map[string]any, []any:
		value, err := anyJSONToValueJSON(data)
		if err != nil {
			return []byte{}, errors.New("Failed to convert JSON to Value: " + err.Error())
		}
		return value.MarshalJSON()
	default:
		// UI用のデータは基本的にMapかListのはずなので、その他の型はエラーとする
		return []byte{}, errors.New("Unsupported top-level JSON type: expected object or array, got " + string(fileData))
	}
}

func anyJSONToValueJSON(data any) (script.Value, error) {
	// まずはListかMapかを判定するために、キーの型を確認する
	switch data := data.(type) {
	case map[string]any:
		res := map[string]script.Value{}
		for k, val := range data {
			v, err := anyJSONToValueJSON(val)
			if err != nil {
				return script.Value{}, errors.New("Failed to convert map value for key '" + k + "': " + err.Error())
			}
			res[k] = v
		}
		return script.NewLitMap(res), nil
	case []any:
		list := make([]script.Value, len(data))
		for i, val := range data {
			v, err := anyJSONToValueJSON(val)
			if err != nil {
				return script.Value{}, errors.New("Failed to convert list item at index " + fmt.Sprint(i) + ": " + err.Error())
			}
			list[i] = v
		}
		return script.NewLitList(list), nil
	case string:
		return script.NewString(data), nil
	case float64:
		return script.NewNumber(data), nil
	case bool:
		return script.NewBool(data), nil
	default:
		return script.Value{}, errors.New("Unsupported JSON value type: " + fmt.Sprintf("%T", data))
	}
}
