package script

// ==================================================
// 組み込みコマンドの登録
// mapを渡して複数のコマンドを登録する
func (vm *VM) registerCmdList(cmds map[string]func(vm *VM, args []Value) (Value, error)) {
	for name, fn := range cmds {
		vm.RegisterCmd(name, fn)
	}
}

// 組み込みコマンドをまとめて登録
func (vm *VM) SetBuiltinCmds() {
	// 数学関数を追加
	vm.registerCmdList(mathCmds)

	// 比較系関数を追加
	vm.registerCmdList(compareCmds)

	// キャスト系関数を追加
	vm.registerCmdList(castCmds)
}

// ==================================================
// Validation
func CheckCmdArgNum(cmdName string, argNum int, args []Value) error {
	if len(args) != argNum {
		return LogErr("invalid number of arguments for " + cmdName + ": expected " + Itoa(argNum) + ", got " + Itoa(len(args)))
	}
	return nil
}

// ==================================================
// 二項演算の共通処理
func binOp(vm *VM, cmdName string, args []Value, fn func(*VM, Value, Value) (Value, error)) (Value, error) {
	// 引数の数をチェック
	if err := CheckCmdArgNum(cmdName, 2, args); err != nil {
		return Value{}, err
	}
	// 参照を展開
	values, err := vm.EvalAll(args)
	if err != nil {
		return Value{}, err
	}

	// リスト同士なら個別に
	if values[0].Type == TypeLitList && values[1].Type == TypeLitList {
		// 同じ長でないとエラー
		if len(values[0].List) != len(values[1].List) {
			return Value{}, LogErr("invalid arg lists for " + cmdName + ": lists length " + Itoa(len(values[0].List)) + " and " + Itoa(len(values[1].List)))
		}
		// 個別に計算
		result := make([]Value, len(values[0].List))
		for i := 0; i < len(values[0].List); i++ {
			result[i], err = fn(vm, values[0].List[i], values[1].List[i])
			if err != nil {
				return Value{}, err
			}
		}
		return NewLitList(result), nil
	}

	// リストと単一の値の計算
	if (values[0].Type == TypeLitList && values[1].Type != TypeLitList) ||
		(values[0].Type != TypeLitList && values[1].Type == TypeLitList) {

		var listArg, valueArg Value
		if values[0].Type == TypeLitList {
			listArg, valueArg = values[0], values[1]
		} else {
			listArg, valueArg = values[1], values[0]
		}

		// 第二引数をリストに適用
		result := make([]Value, len(listArg.List))
		for i := 0; i < len(listArg.List); i++ {
			result[i], err = fn(vm, listArg.List[i], valueArg)
			if err != nil {
				return Value{}, err
			}
		}
		return NewLitList(result), nil
	}

	// 単一同士で計算
	return fn(vm, values[0], values[1])
}

func binOpTypeError(cmd string, arg0 Value, arg1 Value) error {
	return LogErr("invalid types for %s: %s and %s", cmd, arg0.Type.String(), arg1.Type.String())
}

// ==================================================
// 単項演算子の共通処理
func oneOp(vm *VM, cmdName string, args []Value, fn func(*VM, Value) (Value, error)) (Value, error) {
	// 引数の数をチェック
	if err := CheckCmdArgNum(cmdName, 1, args); err != nil {
		return Value{}, err
	}
	// 参照を展開
	values, err := vm.EvalAll(args)
	if err != nil {
		return Value{}, err
	}

	// リストなら全部に適用
	if values[0].Type == TypeLitList {
		result := make([]Value, len(values[0].List))
		for i := 0; i < len(values[0].List); i++ {
			result[i], err = fn(vm, values[0].List[i])
			if err != nil {
				return Value{}, err
			}
		}
		return NewLitList(result), nil
	}

	// 個別に計算
	return fn(vm, values[0])
}

func oneOpTypeError(cmd string, value Value) error {
	return LogErr("invalid type for %s: %s", cmd, value.Type.String())
}
