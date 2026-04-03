package script

import (
	"errors"
)

type Runtime struct {
	vm     *VM
	Result Value
	// コンパイル済みコード
	// トップは単純なコマンドを並べられる暗黙のBetinなので、Listになっている
	ast []Value
}

// クローンを作成
func (runtime *Runtime) Clone() *Runtime {
	newVM := *runtime.vm // vars以外はコピーでいい

	// varsはJSONを経由してクローンを作成
	vmVarsJSON, _ := runtime.vm.vars.MarshalJSON()
	newVM.vars.UnmarshalJSON(vmVarsJSON)

	// 新しいRuntimeを作成
	newRuntime := &Runtime{
		vm:  &newVM,
		ast: runtime.ast, // ASTは不変なのでそのまま参照を渡す
	}
	return newRuntime
}

// **********************************************************************
// 実行
// ==================================================
// ソースコードを実行する関数。とりあえずこれを呼べばOK
func Compile(src string) (*Runtime, error) {
	// ソースコードをパース
	v, err := parse(src)
	if err != nil {
		return nil, err
	}

	// Listが空の場合はエラー
	if len(v.List) == 0 {
		return nil, errors.New("no commands to execute")
	}

	vm := NewVM()
	SetBuiltinCmds(vm)

	return &Runtime{
		vm:  vm,
		ast: v.List, // Listをそのままコンパイル済みコードとして保存
	}, nil
}

func (runtime *Runtime) Run() error {
	results := make([]Value, len(runtime.ast))

	// ASTを順番に評価していく
	for _, v := range runtime.ast {
		// 深さのリセット
		runtime.vm.vars.Map["vm_depth"] = NewNumber(0)
		runtime.vm.vars.Map["vm_depth_max"] = NewNumber(0)

		// 評価
		if result, err := runtime.vm.Eval(v); err != nil {
			return err
		} else {
			results = append(results, result)
		}
	}

	// 結果
	// ----------------------------------------
	// 何も実行されなかった
	if len(results) == 0 {
		return errors.New("no commands were executed")
	}

	// 結果が1つだけならそのまま返す。複数あるならリストにして返す
	if len(results) == 1 {
		runtime.Result = results[0]
	} else {
		runtime.Result = NewList(results)
	}
	return nil
}

// **********************************************************************
// Getter/Setter
func (runtime *Runtime) GetVM() *VM {
	return runtime.vm
}

// コマンドを追加するショートカット
func (runtime *Runtime) RegisterCmd(name string, fn func(vm *VM, args []Value) (Value, error)) {
	runtime.vm.RegisterCmd(name, fn)
}

func (runtime *Runtime) GetAST() Value {
	return NewList(runtime.ast)
}

func (runtime *Runtime) GetVar(name string) Value {
	return runtime.vm.GetVar(name)
}

func (runtime *Runtime) SetVar(name string, value Value) {
	runtime.vm.SetVar(name, value)
}
