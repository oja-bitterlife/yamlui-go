package script

// **********************************************************************
// 評価器(VM)の実装
// ==================================================
// コマンド実装の関数型
type VMCmdFunc func(vm *VM, args []Value) (Value, error)

const (
	VM_VAR_CAPACITY  = 128 // varsの初期容量
	VM_RECURSION_MAX = 32  // 再帰の最大深さ
	VM_REPEAT_MAX    = 16  // repeatの最大回数
)

// 外部との連携用データ構造体
type VM struct {
	ast    []Value              // コンパイル済みコード
	cmds   map[string]VMCmdFunc // コマンド名と実装のマッピング
	vars   Value                // VMのメインメモリ的な
	Result Value                // 最後の評価結果を保存する場所
}

// VMの初期化
func NewVM(valueAST Value) (*VM, error) {
	// VMの初期化
	vm := &VM{
		ast:    valueAST.List,
		cmds:   make(map[string]VMCmdFunc),
		vars:   NewLitMap(make(map[string]Value, VM_VAR_CAPACITY)),
		Result: Value{},
	}
	SetBuiltinCmds(vm)
	vm.ClearVars()

	return vm, nil
}

func (vm *VM) Clone() *VM {
	// vars以外はそのままコピーしても問題ない
	clone := *vm

	// varsはCloneする
	clone.vars = clone.vars.Clone()

	return &clone
}

// **********************************************************************
// 実行
func (vm *VM) Run() error {
	results := make([]Value, len(vm.ast))

	// ASTを順番に評価していく
	for _, v := range vm.ast {
		// 深さのリセット
		vm.SetVar("vm_depth", NewNumber(0))
		vm.SetVar("vm_depth_max", NewNumber(0))

		// 評価
		if ret, err := vm.Eval(v); err != nil {
			return err
		} else {
			results = append(results, ret)
		}
	}

	// 結果
	// ----------------------------------------
	// 結果が1つだけならそのまま返す。複数あるならリストにして返す
	switch len(results) {
	case 0:
		vm.Result = NewNil() // 結果がない場合はnil
	case 1:
		vm.Result = results[0] // 結果が1つだけならそのまま返す
	default:
		vm.Result = NewList(results) // 結果が複数あるならリストにして返す
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

func (vm *VM) ClearVars() {
	clear(vm.vars.Map)
}
