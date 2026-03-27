package script

import "fmt"

type VM struct {
	registry map[string]Value
	vars     map[string]Value
	cmds     map[string]func(args []Value) (Value, error)
}

func NewVM() *VM {
	return &VM{registry: make(map[string]Value)}
}

func (vm *VM) Eval(v Value) (Value, error) {
	switch v.Type {
	case TypeNumber, TypeString:
		return v, nil // リテラルはそのまま返す
	case TypeProperty:
		return vm.vars[v.Prop], nil // 変数の値を返す
	case TypeList:
		return vm.EvalList(v.List) // ここで再帰
	default:
		return Value{}, nil
	}
}

func (vm *VM) EvalList(list []Value) (Value, error) {
	if len(list) == 0 {
		return Value{}, nil
	}

	// 1番目をコマンドとして解釈
	cmd := list[0].Str

	// 特殊形式: set
	if cmd == "set" {
		val, _ := vm.Eval(list[2])  // 右辺を再帰評価
		vm.vars[list[1].Prop] = val // 変数に値をセット
		return val, nil
	}

	// 通常の関数呼び出し: 引数をすべて評価
	args := make([]Value, len(list)-1)
	for i := 1; i < len(list); i++ {
		arg, _ := vm.Eval(list[i]) // 引数ごとに再帰
		args[i-1] = arg
	}

	return vm.Apply(cmd, args)
}

func (vm *VM) Apply(cmd string, args []Value) (Value, error) {
	// コマンドに応じた処理を実装
	fn, ok := vm.cmds[cmd]
	if !ok {
		return Value{}, fmt.Errorf("unknown command: %s", cmd)
	}
	return fn(args)
}
