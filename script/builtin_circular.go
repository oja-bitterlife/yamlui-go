package script

import (
	"math"
)

// **********************************************************************
// 円関数(三角関数)系のbuiltin
var circularCmds = map[string]func(*VM, []Value) (Value, error){
	"sin":   sin,
	"cos":   cos,
	"tan":   tan,
	"atan2": atan2,
}

// ==================================================
// 円関数(三角関数)
func sin(vm *VM, args []Value) (Value, error) {
	return oneOp(vm, "sin", args, func(vm *VM, arg0 Value) (Value, error) {
		if arg0.Type == TypeNumber {
			return NewNumber(math.Sin(arg0.Num)), nil
		}
		return Value{}, oneOpTypeError("sin", arg0)
	})
}

func cos(vm *VM, args []Value) (Value, error) {
	return oneOp(vm, "cos", args, func(vm *VM, arg0 Value) (Value, error) {
		if arg0.Type == TypeNumber {
			return NewNumber(math.Cos(arg0.Num)), nil
		}
		return Value{}, oneOpTypeError("cos", arg0)
	})
}

func tan(vm *VM, args []Value) (Value, error) {
	return oneOp(vm, "tan", args, func(vm *VM, arg0 Value) (Value, error) {
		if arg0.Type == TypeNumber {
			return NewNumber(math.Tan(arg0.Num)), nil
		}
		return Value{}, oneOpTypeError("tan", arg0)
	})
}

func atan2(vm *VM, args []Value) (Value, error) {
	return binOp(vm, "atan2", args, func(vm *VM, arg0 Value, arg1 Value) (Value, error) {
		if arg0.Type == TypeNumber && arg1.Type == TypeNumber {
			return NewNumber(math.Atan2(arg0.Num, arg1.Num)), nil
		}
		return Value{}, binOpTypeError("atan2", arg0, arg1)
	})
}

// ==================================================
// min/max
func min(vm *VM, args []Value) (Value, error) {
	// 引数が1つでリストならリストの最小値を返す
	if len(args) == 1 && args[0].Type == TypeLitList {
		return oneOp(vm, "min", args, func(vm *VM, arg0 Value) (Value, error) {
			if len(arg0.List) == 0 {
				return Value{}, LogErr("min expects a non-empty list")
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
	return binOp(vm, "min", args, func(vm *VM, arg0 Value, arg1 Value) (Value, error) {
		if arg0.ConvertNumber().Num < arg1.ConvertNumber().Num {
			return arg0, nil
		} else {
			return arg1, nil
		}
	})
}

func max(vm *VM, args []Value) (Value, error) {
	// 引数が1つでリストならリストの最大値を返す
	if len(args) == 1 && args[0].Type == TypeLitList {
		return oneOp(vm, "max", args, func(vm *VM, arg0 Value) (Value, error) {
			if len(arg0.List) == 0 {
				return Value{}, LogErr("max expects a non-empty list")
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
	return binOp(vm, "max", args, func(vm *VM, arg0 Value, arg1 Value) (Value, error) {
		if arg0.ConvertNumber().Num > arg1.ConvertNumber().Num {
			return arg0, nil
		} else {
			return arg1, nil
		}
	})
}
