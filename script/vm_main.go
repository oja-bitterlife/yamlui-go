package script

// **********************************************************************
// 評価器(VM)の実装
// ==================================================
// コマンド実装の関数型
type VMCmdFunc func(vm *VM, args []Value) (Value, error)
type VMCmdDispatcher func(cmdName string) (VMCmdFunc, bool)

const (
	VM_VAR_CAPACITY  = 32 // varsの初期容量
	VM_RECURSION_MAX = 16 // 再帰の最大深さ
	VM_REPEAT_MAX    = 16 // repeatの最大回数
)

// 外部との連携用データ構造体
type VM struct {
	AST    Value             // コンパイル済みコード
	cmds   []VMCmdDispatcher // コマンド名と実装のマッピング
	Vars   ValueMap          // VMのメインメモリ的な
	Result Value             // 最後の評価結果を保存する場所
}

// VMの初期化
func NewVM(valueAST Value) VM {
	// VMの初期化
	vm := VM{
		AST:    valueAST,
		cmds:   []VMCmdDispatcher{BuiltinCmds},
		Vars:   make(map[string]Value, VM_VAR_CAPACITY),
		Result: Value{},
	}
	return vm
}

func (vm *VM) Clone() VM {
	// vars以外はそのままコピーしても問題ない
	clone := *vm

	// varsはCloneする
	for k, v := range vm.Vars {
		clone.Vars[k] = v.Clone()
	}

	return clone
}

// **********************************************************************
// 実行
func (vm *VM) Run() error {
	// resultを受け取る準備
	if len(vm.AST.List) <= 1 {
		vm.Result = NewNil() // ASTが空ならnil, ASTが1つでも上書きを期待してnil
	} else {
		vm.Result = NewList(make([]Value, len(vm.AST.List))) // ASTが複数ならリスト
	}

	// ASTを順番に評価していく
	for i, v := range vm.AST.List {
		// 深さのリセット
		vm.Vars["vm_depth"] = NewNumber(0)
		vm.Vars["vm_depth_max"] = NewNumber(0)

		// 評価
		if ret, err := vm.Eval(v); err != nil {
			return err
		} else {
			// 結果を保存
			if vm.Result.IsList() {
				vm.Result.List[i] = ret
			} else {
				vm.Result = ret
			}
		}
	}

	return nil
}

// **********************************************************************
// getter/setter
func (vm *VM) AddCmdIF(cmdDispatcher VMCmdDispatcher) {
	vm.cmds = append(vm.cmds, cmdDispatcher)
	for _, cmdsFunc := range vm.cmds {
		Log("cmds: %v", cmdsFunc)
	}
}

func (vm *VM) HasScript() bool {
	return len(vm.AST.List) > 0
}
