package builtin

import (
	"errors"

	"github.com/oja-bitterlife/yamlui-go/script"
)

// **********************************************************************
// 数学系のbuiltin
var CompareCmds = map[string]func(*script.VM, []script.Value) (script.Value, error){
	">":  grator,
	"<":  less,
	"==": eq,
	"!=": neq,
	">=": ge,
	"<=": le,
}

// ==================================================
// 比較演算
func grator(vm *script.VM, args []script.Value) (script.Value, error) {
	return binOp(vm, ">", args, func(vm *script.VM, values []script.Value) (script.Value, error) {
		if values[0].Type == script.TypeString && values[1].Type == script.TypeString {
			return script.NewBool(values[0].Str > values[1].Str), nil
		}
		if values[0].Type == script.TypeNumber && values[1].Type == script.TypeNumber {
			return script.NewBool(values[0].Num > values[1].Num), nil
		}
		return script.Value{}, errors.New("invalid types for >: " + values[0].Type.ToStr() + " and " + values[1].Type.ToStr())
	})
}

func less(vm *script.VM, args []script.Value) (script.Value, error) {
	return binOp(vm, "<", args, func(vm *script.VM, values []script.Value) (script.Value, error) {
		if values[0].Type == script.TypeString && values[1].Type == script.TypeString {
			return script.NewBool(values[0].Str < values[1].Str), nil
		}
		if values[0].Type == script.TypeNumber && values[1].Type == script.TypeNumber {
			return script.NewBool(values[0].Num < values[1].Num), nil
		}
		return script.Value{}, errors.New("invalid types for <: " + values[0].Type.ToStr() + " and " + values[1].Type.ToStr())
	})
}

func eq(vm *script.VM, args []script.Value) (script.Value, error) {
	return binOp(vm, "==", args, func(vm *script.VM, values []script.Value) (script.Value, error) {
		if values[0].Type == script.TypeString && values[1].Type == script.TypeString {
			return script.NewBool(values[0].Str == values[1].Str), nil
		}
		if values[0].Type == script.TypeNumber && values[1].Type == script.TypeNumber {
			return script.NewBool(values[0].Num == values[1].Num), nil
		}
		if values[0].Type == script.TypeBool && values[1].Type == script.TypeBool {
			return script.NewBool(values[0].Bool == values[1].Bool), nil
		}
		return script.Value{}, errors.New("invalid types for ==: " + values[0].Type.ToStr() + " and " + values[1].Type.ToStr())
	})
}

func neq(vm *script.VM, args []script.Value) (script.Value, error) {
	return binOp(vm, "!=", args, func(vm *script.VM, values []script.Value) (script.Value, error) {
		if values[0].Type == script.TypeString && values[1].Type == script.TypeString {
			return script.NewBool(values[0].Str != values[1].Str), nil
		}
		if values[0].Type == script.TypeNumber && values[1].Type == script.TypeNumber {
			return script.NewBool(values[0].Num != values[1].Num), nil
		}
		if values[0].Type == script.TypeBool && values[1].Type == script.TypeBool {
			return script.NewBool(values[0].Bool != values[1].Bool), nil
		}
		return script.Value{}, errors.New("invalid types for !=: " + values[0].Type.ToStr() + " and " + values[1].Type.ToStr())
	})
}

func ge(vm *script.VM, args []script.Value) (script.Value, error) {
	return binOp(vm, ">=", args, func(vm *script.VM, values []script.Value) (script.Value, error) {
		if values[0].Type == script.TypeString && values[1].Type == script.TypeString {
			return script.NewBool(values[0].Str >= values[1].Str), nil
		}
		if values[0].Type == script.TypeNumber && values[1].Type == script.TypeNumber {
			return script.NewBool(values[0].Num >= values[1].Num), nil
		}
		return script.Value{}, errors.New("invalid types for >=: " + values[0].Type.ToStr() + " and " + values[1].Type.ToStr())
	})
}

func le(vm *script.VM, args []script.Value) (script.Value, error) {
	return binOp(vm, "<=", args, func(vm *script.VM, values []script.Value) (script.Value, error) {
		if values[0].Type == script.TypeString && values[1].Type == script.TypeString {
			return script.NewBool(values[0].Str <= values[1].Str), nil
		}
		if values[0].Type == script.TypeNumber && values[1].Type == script.TypeNumber {
			return script.NewBool(values[0].Num <= values[1].Num), nil
		}
		return script.Value{}, errors.New("invalid types for <=: " + values[0].Type.ToStr() + " and " + values[1].Type.ToStr())
	})
}
