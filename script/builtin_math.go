package script

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"strings"
)

// **********************************************************************
// 数学系のbuiltin
var mathCmds = map[string]func(*VM, []Value) (Value, error){
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
			return NewNumber(float64(int(arg0.Num) % int(arg1.Num))), nil
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

func sqrt(vm *VM, args []Value) (Value, error) {
	return oneOp(vm, "sqrt", args, func(vm *VM, arg0 Value) (Value, error) {
		if arg0.Type == TypeNumber {
			// 虚数は無しで
			if arg0.Num < 0 {
				return Value{}, errors.New("cannot calculate sqrt of negative number: " + strconv.FormatFloat(arg0.Num, 'f', 4, 64))
			}
			return NewNumber(math.Sqrt(arg0.Num)), nil
		}
		return Value{}, oneOpTypeError("sqrt", arg0)
	})
}

func pow(vm *VM, args []Value) (Value, error) {
	return binOp(vm, "pow", args, func(vm *VM, arg0 Value, arg1 Value) (Value, error) {
		if arg0.Type == TypeNumber && arg1.Type == TypeNumber {
			return NewNumber(math.Pow(arg0.Num, arg1.Num)), nil
		}
		return Value{}, binOpTypeError("pow", arg0, arg1)
	})
}

func floor(vm *VM, args []Value) (Value, error) {
	return oneOp(vm, "floor", args, func(vm *VM, arg0 Value) (Value, error) {
		if arg0.Type == TypeNumber {
			return NewNumber(math.Floor(arg0.Num)), nil
		}
		return Value{}, oneOpTypeError("floor", arg0)
	})
}

func ceil(vm *VM, args []Value) (Value, error) {
	return oneOp(vm, "ceil", args, func(vm *VM, arg0 Value) (Value, error) {
		if arg0.Type == TypeNumber {
			return NewNumber(math.Ceil(arg0.Num)), nil
		}
		return Value{}, oneOpTypeError("ceil", arg0)
	})
}

func round(vm *VM, args []Value) (Value, error) {
	return oneOp(vm, "round", args, func(vm *VM, arg0 Value) (Value, error) {
		if arg0.Type == TypeNumber {
			return NewNumber(math.Round(arg0.Num)), nil
		}
		return Value{}, oneOpTypeError("round", arg0)
	})
}

func fit(vm *VM, args []Value) (Value, error) {
	// 引数の数をチェック
	if err := CheckCmdArgNum("fit", 5, args); err != nil {
		return Value{}, err
	}
	// 参照を展開
	values, err := vm.EvalAll(args)
	if err != nil {
		return Value{}, err
	}

	// 一つ目が入力値、二つ目が入力の最小値、三つ目が入力の最大値、四つ目が出力の最小値、五つ目が出力の最大値
	if values[0].Type != TypeNumber || values[1].Type != TypeNumber || values[2].Type != TypeNumber || values[3].Type != TypeNumber || values[4].Type != TypeNumber {
		jsonArgs, _ := json.Marshal(values)
		return Value{}, errors.New("fit: all arguments must be numbers: " + string(jsonArgs))
	}
	input := values[0].Num
	inMin := values[1].Num
	inMax := values[2].Num
	outMin := values[3].Num
	outMax := values[4].Num

	// 入力エラーの時はminを返す
	if inMax == inMin {
		return NewNumber(outMin), nil
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

	return NewNumber(fitValue), nil
}
