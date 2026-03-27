package script

import "fmt"

// 仮想マシンの構造体
type VM struct {
	registry map[string]Value
	vars     map[string]Value
	cmds     map[string]func(args []Value) (Value, error)
}

// VMの初期化
func NewVM() *VM {
	return &VM{registry: make(map[string]Value)}
}

// **********************************************************************
// 評価関数
// ==================================================
// Valueを評価して、最終的な値を返す関数
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

// ==================================================
// リストを評価する関数
func (vm *VM) EvalList(list []Value) (Value, error) {
	if len(list) == 0 {
		return Value{}, nil
	}

	// 1番目をコマンドとして解釈
	cmd := list[0].Str
	args := make([]Value, len(list)-1)

	// 引数をすべて評価
	for i := 1; i < len(list); i++ {
		// リスト（ネストした式）だけは評価して結果を積む
		// それ以外（@y や 10）はそのまま積む
		if list[i].Type == TypeList {
			arg, err := vm.Eval(list[i])
			if err != nil {
				return Value{}, err
			}
			args[i-1] = arg
		} else {
			args[i-1] = list[i]
		}
	}
	return vm.Apply(cmd, args)
}

// ==================================================
// コマンド(リストの先頭)を適用する関数
func (vm *VM) Apply(cmd string, args []Value) (Value, error) {
	// コマンドに応じた処理を実装
	fn, ok := vm.cmds[cmd]
	if !ok {
		return Value{}, fmt.Errorf("unknown command: %s", cmd)
	}

	if cmd == "set" {
		return vm.SetVar(args)
	}

	return fn(args)
}

// **********************************************************************
// Builint-in コマンドの実装
// ==================================================
// set
func (vm *VM) SetVar(args []Value) (Value, error) {
	if len(args) < 2 {
		return Value{}, fmt.Errorf("set requires variable and value")
	}

	// 第一引数は保存先
	target := args[0].Prop

	// 第二引数以降は、このタイミングで Eval して「値」にする
	var final Value
	if len(args) == 2 {
		val, err := vm.Eval(args[1])
		if err != nil {
			return Value{}, err
		}
		final = val
	} else {
		// 可変長の場合
		results := make([]Value, 0, len(args)-1)
		for i := 1; i < len(args); i++ {
			val, err := vm.Eval(args[i])
			if err != nil {
				return Value{}, err
			}
			results = append(results, val)
		}
		final = NewList(results)
	}

	vm.vars[target] = final
	return final, nil
}
