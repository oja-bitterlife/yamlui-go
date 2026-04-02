package script

import (
	"errors"
	"strconv"
	"strings"
)

// **********************************************************************
// 評価器(VM)の実装
// ==================================================
// 外部との連携用データ構造体
type VM struct {
	vars         Value
	cmds         map[string]func(vm *VM, args []Value) (Value, error)
	maxRecursion int // 再帰の最大深さ
	maxRepeat    int // repeatの最大回数
}

// VMの初期化
func NewVM() *VM {
	return &VM{
		vars:         NewLitMap(make(map[string]Value)),
		cmds:         make(map[string]func(vm *VM, args []Value) (Value, error)),
		maxRecursion: 64,
		maxRepeat:    256,
	}
}

// JSONでDumpする
// func (vm *VM) DumpVars() ([]byte, error) {
// 	buf := bytes.Buffer{}
// 	buf.WriteByte('{')
// 	count := 0
// 	for k, v := range vm.vars.Map {
// 		if count > 0 {
// 			buf.WriteByte(',')
// 		}
// 		buf.WriteString(strconv.Quote(k))
// 		buf.WriteByte(':')
// 		jsonData, err := v.MarshalJSON()
// 		if err != nil {
// 			return nil, err
// 		}
// 		buf.Write(jsonData)
// 	}
// 	buf.WriteByte('}')
// 	return buf.Bytes(), nil
// }

// コマンドを登録する関数
// ----------------------------------------
func (vm *VM) RegisterCmd(name string, fn func(vm *VM, args []Value) (Value, error)) {
	vm.cmds[name] = fn
}

func (vm *VM) RegisterCmdList(cmds map[string]func(vm *VM, args []Value) (Value, error)) {
	for name, fn := range cmds {
		vm.RegisterCmd(name, fn)
	}
}

// 変数を取得・設定する関数
func (vm *VM) GetVar(name string) Value {
	return vm.vars.Map[name]
}

func (vm *VM) GetVars() map[string]Value {
	return vm.vars.Map
}

func (vm *VM) SetVar(name string, value Value) {
	vm.vars.Map[name] = value
}

// デバッグ用cmds取得関数
func (vm *VM) GetCmds() map[string]func(vm *VM, args []Value) (Value, error) {
	return vm.cmds
}

// **********************************************************************
// 評価関数
// ==================================================
// Valueを評価して、最終的な値を返す関数
func (vm *VM) Eval(v Value) (Value, error) {
	switch v.Type {
	case TypeNumber, TypeString, TypeBool, TypeLitList:
		return v, nil // リテラルはそのまま返す
	case TypeProperty:
		v = vm.vars.Map[v.Str] // 変数の値を返す
		if !v.IsLiteral() {
			return Value{}, errors.New("expected literal value")
		}
		return v, nil
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
			return Value{}, errors.New("maximum recursion depth exceeded")
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
			return Value{}, errors.New("expected literal value from list evaluation")
		}

		return list, nil
	default:
		return Value{}, errors.New("unknown value type")
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
	// (! val)の場合は展開だけ
	if cmd == "!" {
		// 引数が1つなら内容を返す、複数ならリストにして返す
		if len(args) == 1 {
			return args[0], nil
		}
		return NewLitList(args), nil
	}

	cleanCmd := strings.TrimPrefix(cmd, "!")

	// 組み込みコマンドはここで直接処理する
	switch cleanCmd {
	case "set":
		if cmd != cleanCmd {
			return Value{}, errors.New("set cannot be used with '!' prefix")
		}
		return vm.setVar(args)
	case "switch":
		if cmd != cleanCmd {
			return Value{}, errors.New("switch cannot be used with '!' prefix")
		}
		return vm.switch_(args)
	case "repeat":
		if cmd != cleanCmd {
			return Value{}, errors.New("repeat cannot be used with '!' prefix")
		}
		return vm.repeat(args)
	case "do":
		if cmd != cleanCmd {
			return Value{}, errors.New("do cannot be used with '!' prefix")
		}
		return vm.do(args)
	}

	// コマンドに応じた処理を実装
	fn, ok := vm.cmds[cleanCmd]
	if !ok {
		return Value{}, errors.New("unknown command: " + cmd)
	}

	return fn(vm, args)
}

// **********************************************************************
// Builint-in コマンドの実装
// ==================================================
// set
// 変数に値をセットする
func (vm *VM) setVar(args []Value) (Value, error) {
	if len(args) < 2 {
		return Value{}, errors.New("set requires variable and value")
	}

	// 第一引数は保存先
	target := args[0].Str
	if !strings.HasPrefix(target, "@") && !strings.HasPrefix(target, "_") {
		return Value{}, errors.New("set target must start with '@' or '_': '" + target + "'")
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
		return Value{}, errors.New("case number out of range: " + strconv.Itoa(caseNo))
	}

	// switch先を評価する
	return vm.Eval(args[caseNo])
}

// ==================================================
// repeat
// 繰り返し処理。引数は、カウンタ変数、繰り返し回数、ブロック
func (vm *VM) repeat(args []Value) (Value, error) {
	if len(args) < 2 {
		return Value{}, errors.New("repeat requires count and block")
	}

	// 第１引数は繰り返し用のカウンタ変数（例: @i）
	if args[0].Type != TypeProperty {
		return Value{}, errors.New("first argument must be a property (e.g., @i)")
	}
	counterName := args[0].Str

	// 第２引数は繰り返し回数
	countVal, err := vm.Eval(args[1])
	if err != nil {
		return Value{}, err
	}
	if countVal.Type != TypeNumber {
		return Value{}, errors.New("repeat count must be a number")
	}
	count := int(countVal.Num)

	// 繰り返し回数の上限をチェック
	if count > vm.maxRepeat {
		return Value{}, errors.New("repeat count exceeds maximum: " + strconv.Itoa(count))
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
