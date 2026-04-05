package script

import (
	"strings"
)

// **********************************************************************
// 数学系のbuiltin
var mathCmds = map[string]func(*VM, []Value) (Value, error){
	"+":   add,
	"-":   sub,
	"*":   mul,
	"/":   div,
	"%":   mod,
	"abs": abs,
}

// ==================================================
// 四則演算
func add(vm *VM, args []Value) (Value, error) {
	return binOp(vm, "+", args, func(vm *VM, arg0 Value, arg1 Value) (Value, error) {
		if arg0.Type == TypeString && arg1.Type == TypeString {
			return NewString(arg0.Str + arg1.Str), nil
		}
		if arg0.Type == TypeNumber && arg1.Type == TypeNumber {
			return NewNumber(arg0.Num + arg1.Num), nil
		}
		return Value{}, binOpTypeError("+", arg0, arg1)
	})
}

func sub(vm *VM, args []Value) (Value, error) {
	return binOp(vm, "-", args, func(vm *VM, arg0 Value, arg1 Value) (Value, error) {
		if arg0.Type == TypeString && arg1.Type == TypeString {
			return NewString(arg0.Str + "-" + arg1.Str), nil
		}
		if arg0.Type == TypeNumber && arg1.Type == TypeNumber {
			return NewNumber(arg0.Num - arg1.Num), nil
		}
		return Value{}, binOpTypeError("-", arg0, arg1)
	})
}

func mul(vm *VM, args []Value) (Value, error) {
	return binOp(vm, "*", args, func(vm *VM, arg0 Value, arg1 Value) (Value, error) {
		repeatStr := func(s string, n int) string {
			var str strings.Builder
			for range n {
				str.WriteString(s)
			}
			return str.String()
		}
		if arg0.Type == TypeString && arg1.Type == TypeNumber {
			return NewString(repeatStr(arg0.Str, int(arg1.Num))), nil
		}
		if arg0.Type == TypeNumber && arg1.Type == TypeString {
			return NewString(repeatStr(arg1.Str, int(arg0.Num))), nil
		}
		if arg0.Type == TypeNumber && arg1.Type == TypeNumber {
			return NewNumber(arg0.Num * arg1.Num), nil
		}
		return Value{}, binOpTypeError("*", arg0, arg1)
	})
}

func div(vm *VM, args []Value) (Value, error) {
	return binOp(vm, "/", args, func(vm *VM, arg0 Value, arg1 Value) (Value, error) {
		if arg0.Type == TypeString && arg1.Type == TypeString {
			return NewString(arg0.Str + "/" + arg1.Str), nil
		}
		if arg0.Type == TypeNumber && arg1.Type == TypeNumber {
			return NewNumber(arg0.Num / arg1.Num), nil
		}
		return Value{}, binOpTypeError("/", arg0, arg1)
	})
}

func mod(vm *VM, args []Value) (Value, error) {
	return binOp(vm, "%", args, func(vm *VM, arg0 Value, arg1 Value) (Value, error) {
		if arg0.Type == TypeString && arg1.Type == TypeString {
			return NewString(arg0.Str + "%" + arg1.Str), nil
		}
		if arg0.Type == TypeNumber && arg1.Type == TypeNumber {
			return NewNumber(arg0.Num % arg1.Num), nil
		}
		return Value{}, binOpTypeError("%", arg0, arg1)
	})
}

// ==================================================
// その他の数学関数
func abs(vm *VM, args []Value) (Value, error) {
	return oneOp(vm, "abs", args, func(vm *VM, arg0 Value) (Value, error) {
		// マイナスをプラスに
		if arg0.Type == TypeNumber {
			if arg0.Num < 0 {
				return NewNumber(-arg0.Num), nil
			}
			return arg0, nil
		}

		// 数値以外はエラー
		return Value{}, oneOpTypeError("abs", arg0)
	})
}
