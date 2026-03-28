package builtin

import (
	"errors"

	"github.com/oja-bitterlife/yamlui-go/script"
)

// **********************************************************************
// 数学系のbuiltin
var castCmds = map[string]func(*script.VM, []script.Value) (script.Value, error){
	"bool": bool,
}

func castError(typeName string, v script.Value) error {
	return errors.New("cannot cast to " + typeName + ": " + v.Type.ToStr())
}

// ==================================================
// キャスト
func bool(vm *script.VM, args []script.Value) (script.Value, error) {
	return oneOp(vm, "not", args, func(vm *script.VM, value script.Value) (script.Value, error) {
		switch value.Type {
		case script.TypeBool:
			return value, nil
		case script.TypeNumber:
			return script.NewBool(value.Num != 0), nil
		case script.TypeString:
			return script.NewBool(value.Str != ""), nil
		default:
			return script.Value{}, castError("bool", value)
		}
	})
}
