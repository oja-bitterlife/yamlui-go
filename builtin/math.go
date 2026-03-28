package builtin

import (
	"strings"

	"github.com/oja-bitterlife/yamlui-go/script"
)

// **********************************************************************
// 数学系のbuiltin
var mathCmds = map[string]func(*script.VM, []script.Value) (script.Value, error){
	"+":   add,
	"-":   sub,
	"*":   mul,
	"/":   div,
	"%":   mod,
	"abs": abs,
	"not": not,
}

// ==================================================
// 四則演算
func add(vm *script.VM, args []script.Value) (script.Value, error) {
	return binOp(vm, "+", args, func(vm *script.VM, arg0 script.Value, arg1 script.Value) (script.Value, error) {
		if arg0.Type == script.TypeString && arg1.Type == script.TypeString {
			return script.NewString(arg0.Str + arg1.Str), nil
		}
		if arg0.Type == script.TypeNumber && arg1.Type == script.TypeNumber {
			return script.NewNumber(arg0.Num + arg1.Num), nil
		}
		return script.Value{}, binOpTypeError("+", arg0, arg1)
	})
}

func sub(vm *script.VM, args []script.Value) (script.Value, error) {
	return binOp(vm, "-", args, func(vm *script.VM, arg0 script.Value, arg1 script.Value) (script.Value, error) {
		if arg0.Type == script.TypeString && arg1.Type == script.TypeString {
			return script.NewString(arg0.Str + "-" + arg1.Str), nil
		}
		if arg0.Type == script.TypeNumber && arg1.Type == script.TypeNumber {
			return script.NewNumber(arg0.Num - arg1.Num), nil
		}
		return script.Value{}, binOpTypeError("-", arg0, arg1)
	})
}

func mul(vm *script.VM, args []script.Value) (script.Value, error) {
	return binOp(vm, "*", args, func(vm *script.VM, arg0 script.Value, arg1 script.Value) (script.Value, error) {
		repeatStr := func(s string, n int) string {
			var str strings.Builder
			for range n {
				str.WriteString(s)
			}
			return str.String()
		}
		if arg0.Type == script.TypeString && arg1.Type == script.TypeNumber {
			return script.NewString(repeatStr(arg0.Str, int(arg1.Num))), nil
		}
		if arg0.Type == script.TypeNumber && arg1.Type == script.TypeString {
			return script.NewString(repeatStr(arg1.Str, int(arg0.Num))), nil
		}
		if arg0.Type == script.TypeNumber && arg1.Type == script.TypeNumber {
			return script.NewNumber(arg0.Num * arg1.Num), nil
		}
		return script.Value{}, binOpTypeError("*", arg0, arg1)
	})
}

func div(vm *script.VM, args []script.Value) (script.Value, error) {
	return binOp(vm, "/", args, func(vm *script.VM, arg0 script.Value, arg1 script.Value) (script.Value, error) {
		if arg0.Type == script.TypeString && arg1.Type == script.TypeString {
			return script.NewString(arg0.Str + "/" + arg1.Str), nil
		}
		if arg0.Type == script.TypeNumber && arg1.Type == script.TypeNumber {
			return script.NewNumber(arg0.Num / arg1.Num), nil
		}
		return script.Value{}, binOpTypeError("/", arg0, arg1)
	})
}

func mod(vm *script.VM, args []script.Value) (script.Value, error) {
	return binOp(vm, "%", args, func(vm *script.VM, arg0 script.Value, arg1 script.Value) (script.Value, error) {
		if arg0.Type == script.TypeString && arg1.Type == script.TypeString {
			return script.NewString(arg0.Str + "%" + arg1.Str), nil
		}
		if arg0.Type == script.TypeNumber && arg1.Type == script.TypeNumber {
			return script.NewNumber(float64(int(arg0.Num) % int(arg1.Num))), nil
		}
		return script.Value{}, binOpTypeError("%", arg0, arg1)
	})
}

// ==================================================
// 単項演算子
func abs(vm *script.VM, args []script.Value) (script.Value, error) {
	return oneOp(vm, "abs", args, func(vm *script.VM, arg0 script.Value) (script.Value, error) {
		// マイナスをプラスに
		if arg0.Type == script.TypeNumber {
			if arg0.Num < 0 {
				return script.NewNumber(-arg0.Num), nil
			}
			return arg0, nil
		}

		// 数値以外はエラー
		return script.Value{}, oneOpTypeError("abs", arg0)
	})
}

func not(vm *script.VM, args []script.Value) (script.Value, error) {
	return oneOp(vm, "not", args, func(vm *script.VM, arg0 script.Value) (script.Value, error) {
		// 0と0以外の数値を反転
		if arg0.Type == script.TypeNumber {
			if arg0.Num != 0 {
				return script.NewNumber(0), nil
			}
			return script.NewNumber(1), nil
		}
		// boolを反転
		if arg0.Type == script.TypeBool {
			return script.NewBool(!arg0.Bool), nil
		}

		// その他の値はエラー
		return script.Value{}, oneOpTypeError("not", arg0)
	})
}
