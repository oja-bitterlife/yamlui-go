package builtin

import (
	"github.com/oja-bitterlife/yamlui-go/script"
)

// **********************************************************************
// 数学系のbuiltin
var MathCmds = map[string]func(*script.VM, []script.Value) (script.Value, error){
	"+": add,
	"-": sub,
	"*": mul,
	"/": div,
	"%": mod,
}

func mathOp(vm *script.VM, cmdName string, argNum int, args []script.Value, fn func(*script.VM, []script.Value) (script.Value, error)) (script.Value, error) {
	if err := script.ValidationArgNum(cmdName, argNum, args); err != nil {
		return script.Value{}, err
	}
	values, err := vm.EvalArgs(args)
	if err != nil {
		return script.Value{}, err
	}
	return fn(vm, values)
}

// ==================================================
// 四則演算
func add(vm *script.VM, args []script.Value) (script.Value, error) {
	return mathOp(vm, "+", 2, args, func(vm *script.VM, values []script.Value) (script.Value, error) {
		return script.NewNumber(values[0].Num + values[1].Num), nil
	})
}
func sub(vm *script.VM, args []script.Value) (script.Value, error) {
	return mathOp(vm, "-", 2, args, func(vm *script.VM, values []script.Value) (script.Value, error) {
		return script.NewNumber(values[0].Num - values[1].Num), nil
	})
}
func mul(vm *script.VM, args []script.Value) (script.Value, error) {
	return mathOp(vm, "*", 2, args, func(vm *script.VM, values []script.Value) (script.Value, error) {
		return script.NewNumber(values[0].Num * values[1].Num), nil
	})
}
func div(vm *script.VM, args []script.Value) (script.Value, error) {
	return mathOp(vm, "/", 2, args, func(vm *script.VM, values []script.Value) (script.Value, error) {
		return script.NewNumber(values[0].Num / values[1].Num), nil
	})
}
func mod(vm *script.VM, args []script.Value) (script.Value, error) {
	return mathOp(vm, "%", 2, args, func(vm *script.VM, values []script.Value) (script.Value, error) {
		return script.NewNumber(float64(int(values[0].Num) % int(values[1].Num))), nil
	})
}
