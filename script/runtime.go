package script

import (
	"errors"
)

type Runtime struct {
	vm  *VM
	ast []Value
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

func (runtime *Runtime) Run() (Value, error) {
	var lastVal Value
	var err error
	for _, v := range runtime.ast {
		// 深さのリセット
		runtime.vm.vars.Map["vm_depth"] = NewNumber(0)
		runtime.vm.vars.Map["vm_depth_max"] = NewNumber(0)

		// 評価
		lastVal, err = runtime.vm.Eval(v)
		if err != nil {
			return Value{}, err
		}
	}
	return lastVal, nil
}

// **********************************************************************
// Getter
func (runtime *Runtime) GetVM() *VM {
	return runtime.vm
}

func (runtime *Runtime) GetAST() Value {
	return NewList(runtime.ast)
}

func (runtime *Runtime) GetVar(name string) Value {
	return runtime.vm.GetVar(name)
}
