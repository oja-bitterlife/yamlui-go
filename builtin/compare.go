package builtin

import (
	"github.com/oja-bitterlife/yamlui-go/script"
)

// **********************************************************************
// 数学系のbuiltin
var compareCmds = map[string]func(*script.VM, []script.Value) (script.Value, error){
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
	return binOp(vm, ">", args, func(vm *script.VM, arg0 script.Value, arg1 script.Value) (script.Value, error) {
		if arg0.Type == script.TypeString && arg1.Type == script.TypeString {
			return script.NewBool(arg0.Str > arg1.Str), nil
		}
		if arg0.Type == script.TypeNumber && arg1.Type == script.TypeNumber {
			return script.NewBool(arg0.Num > arg1.Num), nil
		}
		return script.Value{}, binOpTypeError(">", arg0, arg1)
	})
}

func less(vm *script.VM, args []script.Value) (script.Value, error) {
	return binOp(vm, "<", args, func(vm *script.VM, arg0 script.Value, arg1 script.Value) (script.Value, error) {
		if arg0.Type == script.TypeString && arg1.Type == script.TypeString {
			return script.NewBool(arg0.Str < arg1.Str), nil
		}
		if arg0.Type == script.TypeNumber && arg1.Type == script.TypeNumber {
			return script.NewBool(arg0.Num < arg1.Num), nil
		}
		return script.Value{}, binOpTypeError("<", arg0, arg1)
	})
}

func eq(vm *script.VM, args []script.Value) (script.Value, error) {
	return binOp(vm, "==", args, func(vm *script.VM, arg0 script.Value, arg1 script.Value) (script.Value, error) {
		if arg0.Type == script.TypeString && arg1.Type == script.TypeString {
			return script.NewBool(arg0.Str == arg1.Str), nil
		}
		if arg0.Type == script.TypeNumber && arg1.Type == script.TypeNumber {
			return script.NewBool(arg0.Num == arg1.Num), nil
		}
		if arg0.Type == script.TypeBool && arg1.Type == script.TypeBool {
			return script.NewBool(arg0.Bool == arg1.Bool), nil
		}
		return script.Value{}, binOpTypeError("==", arg0, arg1)
	})
}

func neq(vm *script.VM, args []script.Value) (script.Value, error) {
	return binOp(vm, "!=", args, func(vm *script.VM, arg0 script.Value, arg1 script.Value) (script.Value, error) {
		if arg0.Type == script.TypeString && arg1.Type == script.TypeString {
			return script.NewBool(arg0.Str != arg1.Str), nil
		}
		if arg0.Type == script.TypeNumber && arg1.Type == script.TypeNumber {
			return script.NewBool(arg0.Num != arg1.Num), nil
		}
		if arg0.Type == script.TypeBool && arg1.Type == script.TypeBool {
			return script.NewBool(arg0.Bool != arg1.Bool), nil
		}
		return script.Value{}, binOpTypeError("!=: ", arg0, arg1)
	})
}

func ge(vm *script.VM, args []script.Value) (script.Value, error) {
	return binOp(vm, ">=", args, func(vm *script.VM, arg0 script.Value, arg1 script.Value) (script.Value, error) {
		if arg0.Type == script.TypeString && arg1.Type == script.TypeString {
			return script.NewBool(arg0.Str >= arg1.Str), nil
		}
		if arg0.Type == script.TypeNumber && arg1.Type == script.TypeNumber {
			return script.NewBool(arg0.Num >= arg1.Num), nil
		}
		return script.Value{}, binOpTypeError(">=: ", arg0, arg1)
	})
}

func le(vm *script.VM, args []script.Value) (script.Value, error) {
	return binOp(vm, "<=", args, func(vm *script.VM, arg0 script.Value, arg1 script.Value) (script.Value, error) {
		if arg0.Type == script.TypeString && arg1.Type == script.TypeString {
			return script.NewBool(arg0.Str <= arg1.Str), nil
		}
		if arg0.Type == script.TypeNumber && arg1.Type == script.TypeNumber {
			return script.NewBool(arg0.Num <= arg1.Num), nil
		}
		return script.Value{}, binOpTypeError("<=: ", arg0, arg1)
	})
}
