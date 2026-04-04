package script

import (
	"strings"
)

// **********************************************************************
// 評価関数
// ==================================================
// Valueを評価して、最終的な値を返す関数
func (vm *VM) Eval(v Value) (Value, error) {
	switch v.Type {
	case TypeNil:
		return v, nil // nilはそのまま返す
	case TypeNumber, TypeString, TypeBool, TypeLitList:
		return v, nil // リテラルはそのまま返す
	case TypeProperty:
		vv, ok := vm.vars.Map[v.Str] // 変数の値を返す
		if !ok {
			return Value{}, LogErr("undefined variable: " + v.Str)
		}
		return vv, nil
	case TypeList:
		// 再帰の深さをチェック
		// ----------------------------------------
		depth := vm.vars.Map["vm_depth"]
		depth.Num++ // 深さを増やす
		defer func() {
			depth.Num-- // 深さを戻す
			vm.vars.Map["vm_depth"] = depth
		}()
		if int(depth.Num) >= vm.maxRecursion {
			return Value{}, LogErr("maximum recursion depth exceeded: %d", int(depth.Num))
		}
		vm.vars.Map["vm_depth"] = depth

		// デバッグ用: 過去の最大深さを更新
		depthMax := vm.vars.Map["vm_depth_max"]
		if depth.Num > depthMax.Num {
			depthMax.Num = depth.Num
			vm.vars.Map["vm_depth_max"] = depthMax
		}

		// ここで再帰
		// ----------------------------------------
		list, err := vm.evalList(v.List)
		if err != nil {
			return Value{}, err
		}
		if !list.IsLiteral() {
			return Value{}, LogErr("list did not evaluate to a literal: %v", list)
		}

		return list, nil
	default:
		return Value{}, LogErr("unknown value type: %v", v.Type)
	}
}

// Valueのリストを評価する関数
func (vm *VM) EvalAll(values []Value) ([]Value, error) {
	results := make([]Value, len(values))
	for i, arg := range values {
		v, err := vm.Eval(arg)
		if err != nil {
			return nil, err
		}
		results[i] = v
	}
	return results, nil
}

// ==================================================
// リストを評価する関数
func (vm *VM) evalList(list []Value) (Value, error) {
	if len(list) == 0 {
		return Value{}, nil
	}

	// 1番目をコマンドとして解釈
	cmd := list[0].Str

	// !で始まるコマンドは先に引数を評価
	if strings.HasPrefix(cmd, "!") {
		// 引数は先に評価
		args, err := vm.EvalAll(list[1:])
		if err != nil {
			return Value{}, err
		}
		return vm.applyCmd(cmd, args)
	}

	// コマンドを適用
	return vm.applyCmd(cmd, list[1:])
}

// ==================================================
// コマンド(リストの先頭)を適用する関数
func (vm *VM) applyCmd(cmd string, args []Value) (Value, error) {
	// ($ val)の場合は展開だけ
	if cmd == "$" {
		// 引数が1つなら内容を返す、複数ならリストにして返す
		if len(args) == 1 {
			return args[0], nil
		}
		return NewLitList(args), nil
	}

	cleanCmd := strings.TrimPrefix(cmd, "$")

	// 組み込みコマンドはここで直接処理する
	prefixError := func() (Value, error) {
		return Value{}, LogErr("command '%s' cannot be used with '$' prefix", cleanCmd)
	}
	switch cleanCmd {
	case "set":
		if cmd != cleanCmd {
			return prefixError()
		}
		return vm.setVar(args)
	case "switch":
		if cmd != cleanCmd {
			return prefixError()
		}
		return vm.switch_(args)
	case "repeat":
		if cmd != cleanCmd {
			return prefixError()
		}
		return vm.repeat(args)
	case "do":
		if cmd != cleanCmd {
			return prefixError()
		}
		return vm.do(args)
	case "if":
		if cmd != cleanCmd {
			return prefixError()
		}
		return vm.if_(args)
	}

	// コマンドに応じた処理を実装
	fn, ok := vm.cmds[cleanCmd]
	if !ok {
		return Value{}, LogErr("unknown command: " + cleanCmd)
	}

	return fn(vm, args)
}

// **********************************************************************
// 基本構文の実装
// ==================================================
// set
// 変数に値をセットする
func (vm *VM) setVar(args []Value) (Value, error) {
	if len(args) < 2 {
		return Value{}, LogErr("set requires a target and at least one value")
	}

	// 第一引数は保存先
	target := args[0].Str
	if !strings.HasPrefix(target, "@") && !strings.HasPrefix(target, "_") {
		return Value{}, LogErr("set target must start with '@' or '_': " + target)
	}

	// 第二引数以降は、このタイミングで Eval して「値」にする
	var result Value
	results, err := vm.EvalAll(args[1:])
	if err != nil {
		return Value{}, err
	}

	if len(args) == 2 {
		result = results[0] // 単一の値
	} else {
		// 可変長の場合
		result = NewLitList(results)
	}

	vm.vars.Map[target] = result
	return result, nil
}

// ==================================================
// switch
// 最初の引数を評価して、その値に応じたケースを実行する
func (vm *VM) switch_(args []Value) (Value, error) {
	if len(args) < 2 {
		return Value{}, LogErr("switch requires at least an expression and one case")
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
		return Value{}, LogErr("switch expression must evaluate to a number or boolean, got: %v", exprVal)
	}

	// caseNoの範囲をチェック
	if caseNo < 1 || caseNo >= len(args) {
		return Value{}, LogErr("switch case number out of range: %d (number of cases: %d)", caseNo, len(args)-1)
	}

	// switch先を評価する
	return vm.Eval(args[caseNo])
}

// ==================================================
// repeat
// 繰り返し処理。引数は、カウンタ変数、繰り返し回数、ブロック
func (vm *VM) repeat(args []Value) (Value, error) {
	if len(args) < 2 {
		return Value{}, LogErr("repeat requires at least a counter variable and repeat count")
	}

	// 第１引数は繰り返し用のカウンタ変数（例: @i）
	if args[0].Type != TypeProperty {
		return Value{}, LogErr("repeat counter variable must be a property, got: %v", args[0])
	}
	counterName := args[0].Str

	// 第２引数は繰り返し回数
	countVal, err := vm.Eval(args[1])
	if err != nil {
		return Value{}, err
	}
	if countVal.Type != TypeNumber {
		return Value{}, LogErr("repeat count must evaluate to a number, got: %v", countVal)
	}
	count := int(countVal.Num)

	// 繰り返し回数の上限をチェック
	if count > vm.maxRepeat {
		return Value{}, LogErr("repeat count exceeds maximum limit: %d (max: %d)", count, vm.maxRepeat)
	}

	// 繰り返し回数分ループ
	results := make([]Value, count)
	for i := range count {
		// カウンタ変数を現在の値で更新
		vm.vars.Map[counterName] = Value{Type: TypeNumber, Num: float64(i)}

		// 第3引数：ブロックを評価
		result, err := vm.Eval(args[2])
		if err != nil {
			return Value{}, err
		}
		results[i] = result
	}

	return NewLitList(results), nil
}

// ==================================================
// do
// 引数の式を順番に評価して、最後の値を返す
func (vm *VM) do(args []Value) (Value, error) {
	var lastVal Value
	var err error

	// args は Eval される前の「生の式（AST相当）」である必要があります
	for _, arg := range args {
		// 一行ずつ、その時の VM コンテキストで評価
		lastVal, err = vm.Eval(arg)
		if err != nil {
			return Value{}, err
		}
	}
	return lastVal, nil
}

// ==================================================
// 最初の引数を評価して、その値に応じたケースを実行する
func (vm *VM) if_(args []Value) (Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return Value{}, LogErr("if requires 2 or 3 arguments: condition, true case, [false case]")
	}

	// 最初の引数は評価する式
	condVal, err := vm.Eval(args[0])
	if err != nil {
		return Value{}, err
	}
	if condVal.Type != TypeBool {
		return Value{}, LogErr("if condition must evaluate to a boolean, got: %v", condVal)
	}

	// 引数が2つの場合は、条件が真のときのケースとみなす。偽のときは[]を返す
	if len(args) == 2 {
		// argsのお尻に偽のときのケースとして[] を追加する
		args = append(args, NewLitList([]Value{}))
	}
	return vm.switch_(args)
}
