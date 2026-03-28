package builtin

import (
	"errors"
	"math"
	"strings"

	"github.com/oja-bitterlife/yamlui-go/script"
)

// **********************************************************************
// 数学系のbuiltin
var mathCmds = map[string]func(*script.VM, []script.Value) (script.Value, error){
	"+":     add,
	"-":     sub,
	"*":     mul,
	"/":     div,
	"%":     mod,
	"abs":   abs,
	"sqrt":  sqrt,
	"pow":   pow,
	"floor": floor,
	"ceil":  ceil,
	"round": round,
	"fit":   fit,
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
// その他の数学関数
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

func sqrt(vm *script.VM, args []script.Value) (script.Value, error) {
	return oneOp(vm, "sqrt", args, func(vm *script.VM, arg0 script.Value) (script.Value, error) {
		if arg0.Type == script.TypeNumber {
			// 虚数は無しで
			if arg0.Num < 0 {
				return script.Value{}, errors.New("cannot calculate sqrt of negative number: " + arg0.ToStr())
			}
			return script.NewNumber(math.Sqrt(arg0.Num)), nil
		}
		return script.Value{}, oneOpTypeError("sqrt", arg0)
	})
}

func pow(vm *script.VM, args []script.Value) (script.Value, error) {
	return binOp(vm, "pow", args, func(vm *script.VM, arg0 script.Value, arg1 script.Value) (script.Value, error) {
		if arg0.Type == script.TypeNumber && arg1.Type == script.TypeNumber {
			return script.NewNumber(math.Pow(arg0.Num, arg1.Num)), nil
		}
		return script.Value{}, binOpTypeError("pow", arg0, arg1)
	})
}

func floor(vm *script.VM, args []script.Value) (script.Value, error) {
	return oneOp(vm, "floor", args, func(vm *script.VM, arg0 script.Value) (script.Value, error) {
		if arg0.Type == script.TypeNumber {
			return script.NewNumber(math.Floor(arg0.Num)), nil
		}
		return script.Value{}, oneOpTypeError("floor", arg0)
	})
}

func ceil(vm *script.VM, args []script.Value) (script.Value, error) {
	return oneOp(vm, "ceil", args, func(vm *script.VM, arg0 script.Value) (script.Value, error) {
		if arg0.Type == script.TypeNumber {
			return script.NewNumber(math.Ceil(arg0.Num)), nil
		}
		return script.Value{}, oneOpTypeError("ceil", arg0)
	})
}

func round(vm *script.VM, args []script.Value) (script.Value, error) {
	return oneOp(vm, "round", args, func(vm *script.VM, arg0 script.Value) (script.Value, error) {
		if arg0.Type == script.TypeNumber {
			return script.NewNumber(math.Round(arg0.Num)), nil
		}
		return script.Value{}, oneOpTypeError("round", arg0)
	})
}

func fit(vm *script.VM, args []script.Value) (script.Value, error) {
	// 引数の数をチェック
	if err := CheckCmdArgNum("fit", 5, args); err != nil {
		return script.Value{}, err
	}
	// 参照を展開
	values, err := vm.EvalAll(args)
	if err != nil {
		return script.Value{}, err
	}

	// 一つ目が入力値、二つ目が入力の最小値、三つ目が入力の最大値、四つ目が出力の最小値、五つ目が出力の最大値
	if values[0].Type != script.TypeNumber || values[1].Type != script.TypeNumber || values[2].Type != script.TypeNumber || values[3].Type != script.TypeNumber || values[4].Type != script.TypeNumber {
		return script.Value{}, errors.New("invalid types for fit: " + values[0].Type.ToStr() + ", " + values[1].Type.ToStr() + ", " + values[2].Type.ToStr() + ", " + values[3].Type.ToStr() + ", " + values[4].Type.ToStr())
	}
	input := values[0].Num
	inMin := values[1].Num
	inMax := values[2].Num
	outMin := values[3].Num
	outMax := values[4].Num

	// 入力エラーの時はminを返す
	if inMax == inMin {
		return script.NewNumber(outMin), nil
	}

	// 入力値を0から1の範囲に正規化
	norm := (input - inMin) / (inMax - inMin)
	// 正規化された値を出力の範囲にスケーリング
	fitValue := outMin + norm*(outMax-outMin)

	// min/maxでクランプ
	realMin := math.Min(outMin, outMax)
	realMax := math.Max(outMin, outMax)
	if fitValue < realMin {
		fitValue = realMin
	} else if fitValue > realMax {
		fitValue = realMax
	}

	return script.NewNumber(fitValue), nil
}
