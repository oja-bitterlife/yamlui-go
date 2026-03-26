package script

// OpCode は VM が理解する命令の種類
type OpCode int

const (
	OpPushNum  OpCode = iota // 数値をスタックに積む
	OpPushBool               // 真偽値をスタックに積む
	OpPushStr                // 文字列をスタックに積む
	OpPushSym                // 関数名などのシンボルを積む
	OpPushProp               // @y などのプロパティを積む
	OpCall                   // 関数を呼び出す
)

// Instruction は 1 つの実行単位
type Instruction struct {
	Op  OpCode
	Val Value
}

// Compile は AST (Value) を平坦な命令列に変換します
func Compile(v Value) []Instruction {
	var insts []Instruction

	if v.Type == TypeList {
		// (fn arg1 arg2) の形式
		// 1. 引数を先にコンパイル（スタックに積む順序）
		// Lisp では [0] が関数名なので、[1:] を先に処理
		for i := 1; i < len(v.List); i++ {
			insts = append(insts, Compile(v.List[i])...)
		}
		// 2. 最後にその関数 ([0]) を呼び出す命令を追加
		if len(v.List) > 0 {
			insts = append(insts, Instruction{
				Op:  OpCall,
				Val: v.List[0],
			})
		}
	} else {
		// 単体（数値やプロパティ）は PUSH 命令
		insts = append(insts, Instruction{
			Op:  getPushOp(v.Type),
			Val: v,
		})
	}
	return insts
}

func getPushOp(t ValueType) OpCode {
	switch t {
	case TypeNumber:
		return OpPushNum
	case TypeBool:
		return OpPushBool
	case TypeProperty:
		return OpPushProp
	case TypeSymbol:
		return OpPushSym
	default:
		return OpPushStr
	}
}
