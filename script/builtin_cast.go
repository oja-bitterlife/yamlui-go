package script

// **********************************************************************
// キャスト系のbuiltin
var castCmds = map[string]func(*VM, []Value) (Value, error){
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
		switch value.Type {
		case TypeBool:
			return value, nil
		case TypeNumber:
			return NewBool(value.Num != 0), nil
		case TypeString:
			return NewBool(value.Str != ""), nil
		default:
			return Value{}, castError("bool", value)
		}
	})
}

func castStr(vm *VM, args []Value) (Value, error) {
	return oneOp(vm, "str", args, func(vm *VM, value Value) (Value, error) {
		switch value.Type {
		case TypeBool:
			if value.Bool {
				return NewString("true"), nil
			}
			return NewString("false"), nil
		case TypeNumber:
			return NewString(Itoa(value.Num)), nil
		case TypeString:
			return value, nil
		default:
			return Value{}, castError("str", value)
		}
	})
}

func castNum(vm *VM, args []Value) (Value, error) {
	return oneOp(vm, "num", args, func(vm *VM, value Value) (Value, error) {
		switch value.Type {
		case TypeBool:
			if value.Bool {
				return NewNumber(1), nil
			}
			return NewNumber(0), nil
		case TypeNumber:
			return value, nil
		case TypeString:
			num, err := Atoi(value.Str)
			if err != nil {
				return Value{}, castError("num", value)
			}
			return NewNumber(num), nil
		default:
			return Value{}, castError("num", value)
		}
	})
}
