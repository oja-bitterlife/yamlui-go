package script

// **********************************************************************
// 評価器(VM)の実装
// ==================================================
// コマンド実装の関数型
type VMCmdFunc func(vm *VM, args []Value) (Value, error)

const (
	VM_VAR_CAPACITY  = 64 // varsの初期容量
	VM_RECURSION_MAX = 16 // 再帰の最大深さ
	VM_REPEAT_MAX    = 16 // repeatの最大回数
)

// 外部との連携用データ構造体
type VM struct {
	ast    Value                // コンパイル済みコード
	cmds   map[string]VMCmdFunc // コマンド名と実装のマッピング
	vars   Value                // VMのメインメモリ的な
	Result Value                // 最後の評価結果を保存する場所
}

// VMの初期化
func NewVM(valueAST Value) VM {
	// VMの初期化
	vm := VM{
		ast:    valueAST,
		cmds:   make(map[string]VMCmdFunc),
		vars:   NewLitMap(make(map[string]Value, VM_VAR_CAPACITY)),
		Result: Value{},
	}
	vm.SetBuiltinCmds()
	vm.ClearVars()

	return vm
}

func (vm *VM) Clone() VM {
	// vars以外はそのままコピーしても問題ない
	clone := *vm

	// varsはCloneする
	clone.vars = clone.vars.Clone()

	return clone
}

func (vm *VM) ToValue() Value {
	return vm.ast
}

func (vm *VM) HasScript() bool {
	return len(vm.ast.List) > 0
}

// **********************************************************************
// 実行
func (vm *VM) Run() error {
	// resultを受け取る準備
	if len(vm.ast.List) <= 1 {
		vm.Result = NewNil() // ASTが空ならnil, ASTが1つでも上書きを期待してnil
	} else {
		vm.Result = NewList(make([]Value, len(vm.ast.List))) // ASTが複数ならリスト
	}

	// ASTを順番に評価していく
	for i, v := range vm.ast.List {
		// 深さのリセット
		vm.SetVar("vm_depth", NewNumber(0))
		vm.SetVar("vm_depth_max", NewNumber(0))

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
// コマンド用
func (vm *VM) RegisterCmd(name string, fn VMCmdFunc) {
	vm.cmds[name] = fn
}

// デバッグ用cmds取得関数
func (vm *VM) GetCmds() map[string]VMCmdFunc {
	return vm.cmds
}

// **********************************************************************
// vars用
func (vm *VM) GetVar(name string) Value {
	return vm.vars.Map[name]
}

func (vm *VM) GetVars() map[string]Value {
	return vm.vars.Map
}

func (vm *VM) SetVar(name string, value Value) {
	vm.vars.Map[name] = value
}

func (vm *VM) DeleteVar(name string) {
	delete(vm.vars.Map, name)
}

func (vm *VM) ClearVars() {
	clear(vm.vars.Map)
}
