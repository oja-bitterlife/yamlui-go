package builtin

import (
	"errors"
	"strconv"
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
	">":   grator,
	"<":   less,
	"==":  eq,
	"!=":  neq,
	">=":  ge,
	"<=":  le,
}

func mathOp(vm *script.VM, cmdName string, args []script.Value, fn func(*script.VM, []script.Value) (script.Value, error)) (script.Value, error) {
	// 引数の数をチェック
	if err := script.ValidationArgNum(cmdName, 2, args); err != nil {
		return script.Value{}, err
	}
	// 参照を展開
	values, err := vm.EvalArgs(args)
	if err != nil {
		return script.Value{}, err
	}

	// リストなら個別に
	if values[0].Type == script.TypeList && values[1].Type == script.TypeList {
		// 同じ長でないとエラー
		if len(values[0].List) != len(values[1].List) {
			return script.Value{}, errors.New("different length of lists: " + strconv.Itoa(len(values[0].List)) + " and " + strconv.Itoa(len(values[1].List)))
		}
		// 個別に計算
		result := make([]script.Value, len(values[0].List))
		for i := 0; i < len(values[0].List); i++ {
			result[i], err = fn(vm, []script.Value{values[0].List[i], values[1].List[i]})
			if err != nil {
				return script.Value{}, err
			}
		}
		return script.NewList(result), nil
	}
	if values[0].Type == script.TypeList && values[1].Type != script.TypeList {
		// 第二引数をリストに適用
		result := make([]script.Value, len(values[0].List))
		for i := 0; i < len(values[0].List); i++ {
			result[i], err = fn(vm, []script.Value{values[0].List[i], values[1]})
			if err != nil {
				return script.Value{}, err
			}
		}
		return script.NewList(result), nil
	}
	if values[0].Type != script.TypeList && values[1].Type == script.TypeList {
		// 第一引数をリストに適用
		result := make([]script.Value, len(values[1].List))
		for i := 0; i < len(values[1].List); i++ {
			result[i], err = fn(vm, []script.Value{values[0], values[1].List[i]})
			if err != nil {
				return script.Value{}, err
			}
		}
		return script.NewList(result), nil
	}

	// 個別に計算
	return fn(vm, values)
}

// ==================================================
// 四則演算
func add(vm *script.VM, args []script.Value) (script.Value, error) {
	return mathOp(vm, "+", args, func(vm *script.VM, values []script.Value) (script.Value, error) {
		if values[0].Type == script.TypeString || values[1].Type == script.TypeString {
			return script.NewString(values[0].Str + values[1].Str), nil
		}
		if values[0].Type == script.TypeNumber && values[1].Type == script.TypeNumber {
			return script.NewNumber(values[0].Num + values[1].Num), nil
		}
		return script.Value{}, errors.New("invalid types for +: " + values[0].TypeStr() + " and " + values[1].TypeStr())
	})
}

func sub(vm *script.VM, args []script.Value) (script.Value, error) {
	return mathOp(vm, "-", args, func(vm *script.VM, values []script.Value) (script.Value, error) {
		if values[0].Type == script.TypeString || values[1].Type == script.TypeString {
			return script.NewString(values[0].Str + "-" + values[1].Str), nil
		}
		if values[0].Type == script.TypeNumber && values[1].Type == script.TypeNumber {
			return script.NewNumber(values[0].Num - values[1].Num), nil
		}
		return script.Value{}, errors.New("invalid types for -: " + values[0].TypeStr() + " and " + values[1].TypeStr())
	})
}

func mul(vm *script.VM, args []script.Value) (script.Value, error) {
	return mathOp(vm, "*", args, func(vm *script.VM, values []script.Value) (script.Value, error) {
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
		return script.Value{}, errors.New("invalid types for *: " + values[0].TypeStr() + " and " + values[1].TypeStr())
	})
}

func div(vm *script.VM, args []script.Value) (script.Value, error) {
	return mathOp(vm, "/", args, func(vm *script.VM, values []script.Value) (script.Value, error) {
		if values[0].Type == script.TypeString || values[1].Type == script.TypeString {
			return script.NewString(values[0].Str + "/" + values[1].Str), nil
		}
		if values[0].Type == script.TypeNumber && values[1].Type == script.TypeNumber {
			return script.NewNumber(values[0].Num / values[1].Num), nil
		}
		return script.Value{}, errors.New("invalid types for /: " + values[0].TypeStr() + " and " + values[1].TypeStr())
	})
}

func mod(vm *script.VM, args []script.Value) (script.Value, error) {
	return mathOp(vm, "%", args, func(vm *script.VM, values []script.Value) (script.Value, error) {
		if values[0].Type == script.TypeString || values[1].Type == script.TypeString {
			return script.NewString(values[0].Str + "%" + values[1].Str), nil
		}
		if values[0].Type == script.TypeNumber && values[1].Type == script.TypeNumber {
			return script.NewNumber(float64(int(values[0].Num) % int(values[1].Num))), nil
		}
		return script.Value{}, errors.New("invalid types for %: " + values[0].TypeStr() + " and " + values[1].TypeStr())
	})
}

// ==================================================
// 単項演算子
func abs(vm *script.VM, args []script.Value) (script.Value, error) {
	// 引数の数をチェック
	if err := script.ValidationArgNum("abs", 1, args); err != nil {
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
	if err := script.ValidationArgNum("not", 1, args); err != nil {
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
	return script.Value{}, errors.New("invalid type for not: " + value.TypeStr())
}

// ==================================================
// 比較演算
func grator(vm *script.VM, args []script.Value) (script.Value, error) {
	return mathOp(vm, ">", args, func(vm *script.VM, values []script.Value) (script.Value, error) {
		if values[0].Type == script.TypeString && values[1].Type == script.TypeString {
			return script.NewBool(values[0].Str > values[1].Str), nil
		}
		if values[0].Type == script.TypeNumber && values[1].Type == script.TypeNumber {
			return script.NewBool(values[0].Num > values[1].Num), nil
		}
		return script.Value{}, errors.New("invalid types for >: " + values[0].TypeStr() + " and " + values[1].TypeStr())
	})
}

func less(vm *script.VM, args []script.Value) (script.Value, error) {
	return mathOp(vm, "<", args, func(vm *script.VM, values []script.Value) (script.Value, error) {
		if values[0].Type == script.TypeString && values[1].Type == script.TypeString {
			return script.NewBool(values[0].Str < values[1].Str), nil
		}
		if values[0].Type == script.TypeNumber && values[1].Type == script.TypeNumber {
			return script.NewBool(values[0].Num < values[1].Num), nil
		}
		return script.Value{}, errors.New("invalid types for <: " + values[0].TypeStr() + " and " + values[1].TypeStr())
	})
}

func eq(vm *script.VM, args []script.Value) (script.Value, error) {
	return mathOp(vm, "==", args, func(vm *script.VM, values []script.Value) (script.Value, error) {
		if values[0].Type == script.TypeString && values[1].Type == script.TypeString {
			return script.NewBool(values[0].Str == values[1].Str), nil
		}
		if values[0].Type == script.TypeNumber && values[1].Type == script.TypeNumber {
			return script.NewBool(values[0].Num == values[1].Num), nil
		}
		if values[0].Type == script.TypeBool && values[1].Type == script.TypeBool {
			return script.NewBool(values[0].Bool == values[1].Bool), nil
		}
		return script.Value{}, errors.New("invalid types for ==: " + values[0].TypeStr() + " and " + values[1].TypeStr())
	})
}

func neq(vm *script.VM, args []script.Value) (script.Value, error) {
	return mathOp(vm, "!=", args, func(vm *script.VM, values []script.Value) (script.Value, error) {
		if values[0].Type == script.TypeString && values[1].Type == script.TypeString {
			return script.NewBool(values[0].Str != values[1].Str), nil
		}
		if values[0].Type == script.TypeNumber && values[1].Type == script.TypeNumber {
			return script.NewBool(values[0].Num != values[1].Num), nil
		}
		if values[0].Type == script.TypeBool && values[1].Type == script.TypeBool {
			return script.NewBool(values[0].Bool != values[1].Bool), nil
		}
		return script.Value{}, errors.New("invalid types for !=: " + values[0].TypeStr() + " and " + values[1].TypeStr())
	})
}

func ge(vm *script.VM, args []script.Value) (script.Value, error) {
	return mathOp(vm, ">=", args, func(vm *script.VM, values []script.Value) (script.Value, error) {
		if values[0].Type == script.TypeString && values[1].Type == script.TypeString {
			return script.NewBool(values[0].Str >= values[1].Str), nil
		}
		if values[0].Type == script.TypeNumber && values[1].Type == script.TypeNumber {
			return script.NewBool(values[0].Num >= values[1].Num), nil
		}
		return script.Value{}, errors.New("invalid types for >=: " + values[0].TypeStr() + " and " + values[1].TypeStr())
	})
}

func le(vm *script.VM, args []script.Value) (script.Value, error) {
	return mathOp(vm, "<=", args, func(vm *script.VM, values []script.Value) (script.Value, error) {
		if values[0].Type == script.TypeString && values[1].Type == script.TypeString {
			return script.NewBool(values[0].Str <= values[1].Str), nil
		}
		if values[0].Type == script.TypeNumber && values[1].Type == script.TypeNumber {
			return script.NewBool(values[0].Num <= values[1].Num), nil
		}
		return script.Value{}, errors.New("invalid types for <=: " + values[0].TypeStr() + " and " + values[1].TypeStr())
	})
}
