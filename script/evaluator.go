package script

import "errors"

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
		return Value{}, errors.New("unknown command: " + cmd)
	}

	// 組み込みコマンドはここで直接処理する
	switch cmd {
	case "set":
		return vm.SetVar(args)
	case "switch":
		return vm.Switch(args)
	case "repeat":
		return vm.Repeat(args)
	}

	return fn(args)
}

// **********************************************************************
// Builint-in コマンドの実装
// ==================================================
// set
func (vm *VM) SetVar(args []Value) (Value, error) {
	if len(args) < 2 {
		return Value{}, errors.New("set requires variable and value")
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

// ==================================================
// switch
func (vm *VM) Switch(args []Value) (Value, error) {
	if len(args) < 2 {
		return Value{}, errors.New("switch requires an expression and cases")
	}

	// 最初の引数は評価する式
	exprVal, err := vm.Eval(args[0])
	if err != nil {
		return Value{}, err
	}

	// caseNoとして数値にする
	var caseNo int
	switch exprVal.Type {
	case TypeNumber:
		caseNo = int(exprVal.Num) + 1
	case TypeBool:
		if exprVal.Bool {
			caseNo = 1
		} else {
			caseNo = 2
		}
	default:
		return Value{}, errors.New("switch expression must be a number or bool")
	}

	// caseNoの範囲をチェック
	if caseNo < 1 || caseNo >= len(args) {
		return Value{}, errors.New("case number out of range")
	}

	// switch先を評価する
	return vm.Eval(args[caseNo])
}

// ==================================================
// repeat
func (vm *VM) Repeat(args []Value) (Value, error) {
	if len(args) < 2 {
		return Value{}, errors.New("repeat requires count and block")
	}

	// 第１引数は繰り返し用のカウンタ変数（例: @i）
	if args[0].Type != TypeProperty {
		return Value{}, errors.New("first argument must be a property (e.g., @i)")
	}
	counterName := args[0].Prop

	// 第２引数は繰り返し回数
	countVal, err := vm.Eval(args[0])
	if err != nil {
		return Value{}, err
	}
	if countVal.Type != TypeNumber {
		return Value{}, errors.New("repeat count must be a number")
	}
	count := int(countVal.Num)

	// 繰り返し回数分ループ
	var lastVal Value
	for i := range count {
		// カウンタ変数を現在の値で更新
		vm.vars[counterName] = Value{Type: TypeNumber, Num: float64(i)}

		// 第3引数：ブロックを評価
		val, err := vm.Eval(args[2])
		if err != nil {
			return Value{}, err
		}
		lastVal = val
	}

	return lastVal, nil
}
