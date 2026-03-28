package builtin

import (
	"errors"
	"strconv"

	"github.com/oja-bitterlife/yamlui-go/script"
)

// **********************************************************************
// キャスト系のbuiltin
var castCmds = map[string]func(*script.VM, []script.Value) (script.Value, error){
	"bool": bool,
	"str":  str,
	"num":  num,
}

func castError(typeName string, v script.Value) error {
	return errors.New("cannot cast to " + typeName + ": " + v.Type.ToStr())
}

// ==================================================
// 基本型キャスト
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

func str(vm *script.VM, args []script.Value) (script.Value, error) {
	return oneOp(vm, "str", args, func(vm *script.VM, value script.Value) (script.Value, error) {
		switch value.Type {
		case script.TypeBool:
			if value.Bool {
				return script.NewString("true"), nil
			}
			return script.NewString("false"), nil
		case script.TypeNumber:
			return script.NewString(value.ToStr()), nil
		case script.TypeString:
			return value, nil
		default:
			return script.Value{}, castError("str", value)
		}
	})
}

func num(vm *script.VM, args []script.Value) (script.Value, error) {
	return oneOp(vm, "num", args, func(vm *script.VM, value script.Value) (script.Value, error) {
		switch value.Type {
		case script.TypeBool:
			if value.Bool {
				return script.NewNumber(1), nil
			}
			return script.NewNumber(0), nil
		case script.TypeNumber:
			return value, nil
		case script.TypeString:
			num, err := strconv.ParseFloat(value.Str, 64)
			if err != nil {
				return script.Value{}, castError("num", value)
			}
			return script.NewNumber(num), nil
		default:
			return script.Value{}, castError("num", value)
		}
	})
}
