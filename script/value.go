package script

import (
	"strconv"
	"strings"
)

// ==================================================
// 値の型定義
type ValueType int

const (
	// 評価済み
	TypeNumber ValueType = iota
	TypeBool
	TypeString
	TypeLitList // 評価済みリスト

	// 未評価
	TypeProperty // @i などの参照
	TypeList     // 未評価リスト
)

// 値の型を文字列で返す. デバッグ用.
func (vt ValueType) ToStr() string {
	switch vt {
	case TypeNumber:
		return "Number"
	case TypeBool:
		return "Bool"
	case TypeString:
		return "String"
	case TypeLitList:
		return "LiteralList"
	case TypeProperty:
		return "Property"
	case TypeList:
		return "List"
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
}

// 値を文字列で返す. デバッグ用.
func (v Value) ToStr() string {
	switch v.Type {
	case TypeNumber:
		// 小数点以下が0なら整数として表示
		fmod := v.Num - float64(int64(v.Num))
		if fmod == 0 {
			return strconv.FormatInt(int64(v.Num), 10)
		} else {
			return strconv.FormatFloat(v.Num, 'f', 4, 64)
		}
	case TypeBool:
		return strconv.FormatBool(v.Bool)
	case TypeString, TypeProperty:
		return v.Str
	case TypeLitList, TypeList:
		strs := make([]string, len(v.List))
		for i, elem := range v.List {
			strs[i] = elem.ToStr()
		}
		return "[" + strings.Join(strs, ", ") + "]"
	default:
		return "Unknown"
	}
}

// ==================================================
// 値の生成関数
func NewNumber(f float64) Value  { return Value{Type: TypeNumber, Num: f} }
func NewBool(b bool) Value       { return Value{Type: TypeBool, Bool: b} }
func NewString(s string) Value   { return Value{Type: TypeString, Str: s} }
func NewLitList(v []Value) Value { return Value{Type: TypeLitList, List: v} }
func NewProperty(s string) Value { return Value{Type: TypeProperty, Str: s} }
func NewList(v []Value) Value    { return Value{Type: TypeList, List: v} }

// ==================================================
// リテラルチェック
func (v Value) IsLiteral() bool {
	switch v.Type {
	case TypeNumber, TypeBool, TypeString, TypeLitList:
		return true
	default:
		return false
	}
}

// リストチェック
func (v Value) IsList() bool {
	switch v.Type {
	case TypeList, TypeLitList:
		return true
	default:
		return false
	}
}
