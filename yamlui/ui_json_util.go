package yamlui

import (
	"fmt"

	"github.com/oja-bitterlife/yamlui-go/script"
)

func AnyToValue(data any) (script.Value, error) {
	switch data := data.(type) {
	case map[string]any, []any:
		return anyToValue(data)
	default:
		// UI用のデータは基本的にMapかListのはずなので、その他の型はエラーとする
		return script.Value{}, fmt.Errorf("Unsupported top-level data type: %T", data)
	}
}

func anyToValue(data any) (script.Value, error) {
	// まずはListかMapかを判定するために、キーの型を確認する
	switch data := data.(type) {
	case map[string]any:
		res := map[string]script.Value{}
		for k, val := range data {
			v, err := anyToValue(val)
			if err != nil {
				return script.Value{}, fmt.Errorf("Failed to convert value for key '%s': %v", k, err)
			}
			res[k] = v
		}
		return script.NewLitMap(res), nil
	case []any:
		list := make([]script.Value, len(data))
		for i, val := range data {
			v, err := anyToValue(val)
			if err != nil {
				return script.Value{}, fmt.Errorf("Failed to convert list element at index %d: %v", i, err)
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
		return script.Value{}, fmt.Errorf("Unsupported data type: %T", data)
	}
}
