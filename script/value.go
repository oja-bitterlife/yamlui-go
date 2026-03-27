package script

// 値の型定義
type ValueType int

const (
	TypeNumber ValueType = iota
	TypeBool
	TypeString
	TypeProperty // @i などの参照（未評価）
	TypeList     // (op arg1 arg2)
	TypeQuote    // リテラル

	// 制御構造の型
	TypeSwitch // switch文
	TypeRepeat // repeat文
)

// Listに格納する型付きの値。リフレクションを避けるため全部入り
type Value struct {
	// 型情報
	Type ValueType

	// 基本データ
	Num  float64
	Bool bool
	Str  string

	// リスト構造 (S式や、分解済みの文字列フラグメント)
	List []Value
}

// 値の生成関数
func NewNumber(f float64) Value  { return Value{Type: TypeNumber, Num: f} }
func NewBool(b bool) Value       { return Value{Type: TypeBool, Bool: b} }
func NewString(s string) Value   { return Value{Type: TypeString, Str: s} }
func NewProperty(s string) Value { return Value{Type: TypeProperty, Str: s} }
func NewList(v []Value) Value    { return Value{Type: TypeList, List: v} }
func NewQuote(v []Value) Value   { return Value{Type: TypeQuote, List: v} }
