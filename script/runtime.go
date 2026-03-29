package script

import (
	"errors"
)

type Runtime struct {
	vm       *VM
	compiled []Value
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
		vm:       vm,
		compiled: v.List, // Listをそのままコンパイル済みコードとして保存
	}, nil
}

func (runtime *Runtime) Run() (*Value, error) {
	// リスト(root)に入ってやってくる
	var lastVal Value
	var err error
	for _, v := range runtime.compiled {
		lastVal, err = runtime.vm.Eval(v)
		if err != nil {
			return nil, err
		}
	}
	return &lastVal, nil
}

// ==================================================
// CompileしてRunまで一気にやる
func Run(src string) (*Value, error) {
	runtime, err := Compile(src)
	if err != nil {
		return nil, err
	}
	return runtime.Run()
}

// **********************************************************************
// Getter
func (runtime *Runtime) GetVM() *VM {
	return runtime.vm
}
