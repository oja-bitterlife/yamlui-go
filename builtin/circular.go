package builtin

import (
	"errors"
	"math"

	"github.com/oja-bitterlife/yamlui-go/script"
)

// **********************************************************************
// 円関数(三角関数)系のbuiltin
var circularCmds = map[string]func(*script.VM, []script.Value) (script.Value, error){
	"sin":   sin,
	"cos":   cos,
	"tan":   tan,
	"atan2": atan2,
}

// ==================================================
// 円関数(三角関数)
func sin(vm *script.VM, args []script.Value) (script.Value, error) {
	return oneOp(vm, "sin", args, func(vm *script.VM, arg0 script.Value) (script.Value, error) {
		if arg0.Type == script.TypeNumber {
			return script.NewNumber(math.Sin(arg0.Num)), nil
		}
		return script.Value{}, oneOpTypeError("sin", arg0)
	})
}

func cos(vm *script.VM, args []script.Value) (script.Value, error) {
	return oneOp(vm, "cos", args, func(vm *script.VM, arg0 script.Value) (script.Value, error) {
		if arg0.Type == script.TypeNumber {
			return script.NewNumber(math.Cos(arg0.Num)), nil
		}
		return script.Value{}, oneOpTypeError("cos", arg0)
	})
}

func tan(vm *script.VM, args []script.Value) (script.Value, error) {
	return oneOp(vm, "tan", args, func(vm *script.VM, arg0 script.Value) (script.Value, error) {
		if arg0.Type == script.TypeNumber {
			return script.NewNumber(math.Tan(arg0.Num)), nil
		}
		return script.Value{}, oneOpTypeError("tan", arg0)
	})
}

func atan2(vm *script.VM, args []script.Value) (script.Value, error) {
	return binOp(vm, "atan2", args, func(vm *script.VM, arg0 script.Value, arg1 script.Value) (script.Value, error) {
		if arg0.Type == script.TypeNumber && arg1.Type == script.TypeNumber {
			return script.NewNumber(math.Atan2(arg0.Num, arg1.Num)), nil
		}
		return script.Value{}, binOpTypeError("atan2", arg0, arg1)
	})
}

// ==================================================
// min/max
func min(vm *script.VM, args []script.Value) (script.Value, error) {
	// 引数が1つでリストならリストの最小値を返す
	if len(args) == 1 && args[0].Type == script.TypeLitList {
		return oneOp(vm, "min", args, func(vm *script.VM, arg0 script.Value) (script.Value, error) {
			if len(arg0.List) == 0 {
				return script.Value{}, errors.New("min expects a non-empty list")
			}
			minVal := arg0.List[0]
			for _, v := range arg0.List[1:] {
				if v.ConvertNumber().Num < minVal.ConvertNumber().Num {
					minVal = v
				}
			}
			return minVal, nil
		})
	}

	// 引数が2つなら小さい方を返す
	return binOp(vm, "min", args, func(vm *script.VM, arg0 script.Value, arg1 script.Value) (script.Value, error) {
		if arg0.ConvertNumber().Num < arg1.ConvertNumber().Num {
			return arg0, nil
		} else {
			return arg1, nil
		}
	})
}

func max(vm *script.VM, args []script.Value) (script.Value, error) {
	// 引数が1つでリストならリストの最大値を返す
	if len(args) == 1 && args[0].Type == script.TypeLitList {
		return oneOp(vm, "max", args, func(vm *script.VM, arg0 script.Value) (script.Value, error) {
			if len(arg0.List) == 0 {
				return script.Value{}, errors.New("max expects a non-empty list")
			}
			maxVal := arg0.List[0]
			for _, v := range arg0.List[1:] {
				if v.ConvertNumber().Num > maxVal.ConvertNumber().Num {
					maxVal = v
				}
			}
			return maxVal, nil
		})
	}

	// 引数が2つなら大きい方を返す
	return binOp(vm, "max", args, func(vm *script.VM, arg0 script.Value, arg1 script.Value) (script.Value, error) {
		if arg0.ConvertNumber().Num > arg1.ConvertNumber().Num {
			return arg0, nil
		} else {
			return arg1, nil
		}
	})
}
