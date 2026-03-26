package script

type ValueType int

const (
	TypeNil ValueType = iota
	TypeNumber
	TypeBool
	TypeStr
	TypeSymbol
	TypeProperty
	TypeList // AST（ツリー構造）の時のみ使用
)

type Value struct {
	Type   ValueType
	Number float64
	Bool   bool
	Str    string
	List   []Value // 子要素（S式の入れ子）
}
