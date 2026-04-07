package script

// **********************************************************************
// キャスト系のbuiltin
var castCmds = map[string]VMCmdFunc{
	"bool": castBool,
	"str":  castStr,
	"num":  castNum,
}

func castError(typeName string, v Value) error {
	return LogErr("cannot cast %s to %s", v.Type.String(), typeName)
}

// ==================================================
// 基本型キャスト
func castBool(vm *VM, args []Value) (Value, error) {
	return oneOp(vm, "not", args, func(vm *VM, value Value) (Value, error) {
		return value.ConvertBool(), nil
	})
}

func castStr(vm *VM, args []Value) (Value, error) {
	return oneOp(vm, "str", args, func(vm *VM, value Value) (Value, error) {
		switch value.Type {
		case TypeBool, TypeNumber, TypeString:
			return value.ConvertString(), nil
		default:
			return Value{}, castError("str", value)
		}
	})
}

func castNum(vm *VM, args []Value) (Value, error) {
	return oneOp(vm, "num", args, func(vm *VM, value Value) (Value, error) {
		return value.ConvertNumber(), nil
	})
}
