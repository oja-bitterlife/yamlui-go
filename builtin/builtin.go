package builtin

import (
	"errors"
	"strconv"

	"github.com/oja-bitterlife/yamlui-go/script"
)

// ==================================================
// 組み込みコマンドの登録
func SetBuiltinCmds(vm *script.VM) {
	// 数学関数を追加
	vm.RegisterCmdList(mathCmds)

	// 円関数を追加
	vm.RegisterCmdList(circularCmds)

	// 比較系関数を追加
	vm.RegisterCmdList(compareCmds)

	// キャスト系関数を追加
	vm.RegisterCmdList(castCmds)
}

// ==================================================
// Validation
func CheckCmdArgNum(cmdName string, argNum int, args []script.Value) error {
	if len(args) != argNum {
		return errors.New(cmdName + " expects " + strconv.Itoa(argNum) + " arguments, but got " + strconv.Itoa(len(args)))
	}
	return nil
}

// ==================================================
// 二項演算の共通処理
func binOp(vm *script.VM, cmdName string, args []script.Value, fn func(*script.VM, script.Value, script.Value) (script.Value, error)) (script.Value, error) {
	// 引数の数をチェック
	if err := CheckCmdArgNum(cmdName, 2, args); err != nil {
		return script.Value{}, err
	}
	// 参照を展開
	values, err := vm.EvalAll(args)
	if err != nil {
		return script.Value{}, err
	}

	// リストなら個別に
	if values[0].Type == script.TypeLitList && values[1].Type == script.TypeLitList {
		// どちらもリストなら同じ長でないとエラー
		// 同じ長でないとエラー
		if len(values[0].List) != len(values[1].List) {
			return script.Value{}, errors.New("different length of lists: " + strconv.Itoa(len(values[0].List)) + " and " + strconv.Itoa(len(values[1].List)))
		}
		// 個別に計算
		result := make([]script.Value, len(values[0].List))
		for i := 0; i < len(values[0].List); i++ {
			result[i], err = fn(vm, values[0].List[i], values[1].List[i])
			if err != nil {
				return script.Value{}, err
			}
		}
		return script.NewLitList(result), nil
	}
	if values[0].Type == script.TypeLitList && values[1].Type != script.TypeLitList {
		// 第二引数をリストに適用
		result := make([]script.Value, len(values[0].List))
		for i := 0; i < len(values[0].List); i++ {
			result[i], err = fn(vm, values[0].List[i], values[1])
			if err != nil {
				return script.Value{}, err
			}
		}
		return script.NewLitList(result), nil
	}
	if values[0].Type != script.TypeLitList && values[1].Type == script.TypeLitList {
		// 第一引数をリストに適用
		result := make([]script.Value, len(values[1].List))
		for i := 0; i < len(values[1].List); i++ {
			result[i], err = fn(vm, values[0], values[1].List[i])
			if err != nil {
				return script.Value{}, err
			}
		}
		return script.NewLitList(result), nil
	}

	// 個別に計算
	return fn(vm, values[0], values[1])
}

func binOpTypeError(cmd string, arg0 script.Value, arg1 script.Value) error {
	return errors.New("invalid types for " + cmd + ": " + arg0.Type.ToStr() + " and " + arg1.Type.ToStr())
}

// ==================================================
// 単項演算子の共通処理
func oneOp(vm *script.VM, cmdName string, args []script.Value, fn func(*script.VM, script.Value) (script.Value, error)) (script.Value, error) {
	// 引数の数をチェック
	if err := CheckCmdArgNum(cmdName, 1, args); err != nil {
		return script.Value{}, err
	}
	// 参照を展開
	values, err := vm.EvalAll(args)
	if err != nil {
		return script.Value{}, err
	}

	// リストなら全部に適用
	if values[0].Type == script.TypeLitList {
		result := make([]script.Value, len(values[0].List))
		for i := 0; i < len(values[0].List); i++ {
			result[i], err = fn(vm, values[0].List[i])
			if err != nil {
				return script.Value{}, err
			}
		}
		return script.NewLitList(result), nil
	}

	// 個別に計算
	return fn(vm, values[0])
}

func oneOpTypeError(cmd string, value script.Value) error {
	return errors.New("invalid type for " + cmd + ": " + value.Type.ToStr())
}
