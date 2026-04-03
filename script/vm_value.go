package script

import (
	"strconv"
	"strings"
)

// **********************************************************************
// 値の型定義
// ==================================================
// 値の種類
type ValueType int

const (
	// 評価済み
	TypeNumber ValueType = iota
	TypeBool
	TypeString
	TypeLitList // 評価済みリスト
	TypeLitMap  // 評価済みマップ

	// 未評価
	TypeProperty // @i などの参照
	TypeList     // 未評価リスト
)

const (
	TypeNumberStr   = "Num"
	TypeBoolStr     = "Bool"
	TypeStringStr   = "Str"
	TypeLitListStr  = "LitList"
	TypeLitMapStr   = "LitMap"
	TypePropertyStr = "Prop"
	TypeListStr     = "List"
)

func (t ValueType) String() string {
	switch t {
	case TypeNumber:
		return TypeNumberStr
	case TypeBool:
		return TypeBoolStr
	case TypeString:
		return TypeStringStr
	case TypeLitList:
		return TypeLitListStr
	case TypeLitMap:
		return TypeLitMapStr
	case TypeProperty:
		return TypePropertyStr
	case TypeList:
		return TypeListStr
	default:
		return "Unknown"
	}
}

// ==================================================
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
	Map  ValueMap
}

func (v Value) String() string {
	jsonBytes, err := v.MarshalJSON()
	if err != nil {
		return "Error: " + err.Error()
	}
	return string(jsonBytes)
}

// ==================================================
// 値の生成関数
func NewNumber(f float64) Value  { return Value{Type: TypeNumber, Num: f} }
func NewBool(b bool) Value       { return Value{Type: TypeBool, Bool: b} }
func NewString(s string) Value   { return Value{Type: TypeString, Str: s} }
func NewLitList(v []Value) Value { return Value{Type: TypeLitList, List: v} }
func NewLitMap(v ValueMap) Value { return Value{Type: TypeLitMap, Map: v} }
func NewProperty(s string) Value { return Value{Type: TypeProperty, Str: s} }
func NewList(v []Value) Value    { return Value{Type: TypeList, List: v} }

// ==================================================
// リテラルチェック
func (v Value) IsLiteral() bool {
	switch v.Type {
	case TypeNumber, TypeBool, TypeString, TypeLitList, TypeLitMap:
		return true
	default:
		return false
	}
}

// リストチェック
func (v Value) IsList() bool {
	switch v.Type {
	case TypeList, TypeLitList, TypeLitMap:
		return true
	default:
		return false
	}
}

// ==================================================
// コンバート
func (v Value) ConvertBool() Value {
	switch v.Type {
	case TypeBool:
		return v
	case TypeNumber:
		return NewBool(v.Num != 0)
	case TypeString:
		return NewBool(strings.ToLower(v.Str) == "true")
	case TypeLitList, TypeList:
		return NewBool(len(v.List) > 0)
	case TypeLitMap:
		return NewBool(len(v.Map) > 0)
	default:
		return NewBool(false)
	}
}

func (v Value) ConvertNumber() Value {
	switch v.Type {
	case TypeNumber:
		return v
	case TypeBool:
		if v.Bool {
			return NewNumber(1)
		} else {
			return NewNumber(0)
		}
	case TypeString:
		f, err := strconv.ParseFloat(v.Str, 64)
		if err != nil {
			return NewNumber(0)
		}
		return NewNumber(f)
	case TypeLitList, TypeList:
		return NewNumber(float64(len(v.List)))
	case TypeLitMap:
		return NewNumber(float64(len(v.Map)))
	default:
		return NewNumber(0)
	}
}
