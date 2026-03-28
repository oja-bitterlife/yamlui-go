package builtin

import (
	"errors"
	"strings"

	"github.com/oja-bitterlife/yamlui-go/script"
)

// **********************************************************************
// 数学系のbuiltin
var MathCmds = map[string]func(*script.VM, []script.Value) (script.Value, error){
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
	return binOp(vm, "+", args, func(vm *script.VM, values []script.Value) (script.Value, error) {
		if values[0].Type == script.TypeString && values[1].Type == script.TypeString {
			return script.NewString(values[0].Str + values[1].Str), nil
		}
		if values[0].Type == script.TypeNumber && values[1].Type == script.TypeNumber {
			return script.NewNumber(values[0].Num + values[1].Num), nil
		}
		return script.Value{}, errors.New("invalid types for +: " + values[0].Type.ToStr() + " and " + values[1].Type.ToStr())
	})
}

func sub(vm *script.VM, args []script.Value) (script.Value, error) {
	return binOp(vm, "-", args, func(vm *script.VM, values []script.Value) (script.Value, error) {
		if values[0].Type == script.TypeString && values[1].Type == script.TypeString {
			return script.NewString(values[0].Str + "-" + values[1].Str), nil
		}
		if values[0].Type == script.TypeNumber && values[1].Type == script.TypeNumber {
			return script.NewNumber(values[0].Num - values[1].Num), nil
		}
		return script.Value{}, errors.New("invalid types for -: " + values[0].Type.ToStr() + " and " + values[1].Type.ToStr())
	})
}

func mul(vm *script.VM, args []script.Value) (script.Value, error) {
	return binOp(vm, "*", args, func(vm *script.VM, values []script.Value) (script.Value, error) {
		repeatStr := func(s string, n int) string {
			var str strings.Builder
			for range n {
				str.WriteString(s)
			}
			return str.String()
		}
		if values[0].Type == script.TypeString && values[1].Type == script.TypeNumber {
			return script.NewString(repeatStr(values[0].Str, int(values[1].Num))), nil
		}
		if values[0].Type == script.TypeNumber && values[1].Type == script.TypeString {
			return script.NewString(repeatStr(values[1].Str, int(values[0].Num))), nil
		}
		if values[0].Type == script.TypeNumber && values[1].Type == script.TypeNumber {
			return script.NewNumber(values[0].Num * values[1].Num), nil
		}
		return script.Value{}, errors.New("invalid types for *: " + values[0].Type.ToStr() + " and " + values[1].Type.ToStr())
	})
}

func div(vm *script.VM, args []script.Value) (script.Value, error) {
	return binOp(vm, "/", args, func(vm *script.VM, values []script.Value) (script.Value, error) {
		if values[0].Type == script.TypeString && values[1].Type == script.TypeString {
			return script.NewString(values[0].Str + "/" + values[1].Str), nil
		}
		if values[0].Type == script.TypeNumber && values[1].Type == script.TypeNumber {
			return script.NewNumber(values[0].Num / values[1].Num), nil
		}
		return script.Value{}, errors.New("invalid types for /: " + values[0].Type.ToStr() + " and " + values[1].Type.ToStr())
	})
}

func mod(vm *script.VM, args []script.Value) (script.Value, error) {
	return binOp(vm, "%", args, func(vm *script.VM, values []script.Value) (script.Value, error) {
		if values[0].Type == script.TypeString && values[1].Type == script.TypeString {
			return script.NewString(values[0].Str + "%" + values[1].Str), nil
		}
		if values[0].Type == script.TypeNumber && values[1].Type == script.TypeNumber {
			return script.NewNumber(float64(int(values[0].Num) % int(values[1].Num))), nil
		}
		return script.Value{}, errors.New("invalid types for %: " + values[0].Type.ToStr() + " and " + values[1].Type.ToStr())
	})
}

// ==================================================
// 単項演算子
func abs(vm *script.VM, args []script.Value) (script.Value, error) {
	// 引数の数をチェック
	if err := CheckCmdArgNum("abs", 1, args); err != nil {
		return script.Value{}, err
	}
	// 参照を展開
	value, err := vm.Eval(args[0])
	if err != nil {
		return script.Value{}, err
	}

	// マイナスをプラスに
	if value.Type == script.TypeNumber {
		if value.Num < 0 {
			return script.NewNumber(-value.Num), nil
		}
		return value, nil
	}

	// 数値以外はそのまま返す
	return value, nil
}

func not(vm *script.VM, args []script.Value) (script.Value, error) {
	// 引数の数をチェック
	if err := CheckCmdArgNum("not", 1, args); err != nil {
		return script.Value{}, err
	}
	// 参照を展開
	value, err := vm.Eval(args[0])
	if err != nil {
		return script.Value{}, err
	}

	// 0と0以外の数値を反転
	if value.Type == script.TypeNumber {
		if value.Num != 0 {
			return script.NewNumber(0), nil
		}
		return script.NewNumber(1), nil
	}
	// bolを反転
	if value.Type == script.TypeBool {
		return script.NewBool(!value.Bool), nil
	}

	// その他の値はエラー
	return script.Value{}, errors.New("invalid type for not: " + value.Type.ToStr())
}
