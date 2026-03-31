package script

import (
	"bytes"
	"encoding/json"
	"errors"
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

	// 未評価
	TypeProperty // @i などの参照
	TypeList     // 未評価リスト
)

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
	default:
		return NewNumber(0)
	}
}

// **********************************************************************
// Marshal/Unmarshal
// json.Marshaler インターフェースの実装
func (v Value) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	switch v.Type {
	case TypeNumber:
		buf.WriteString(`{"Num":`)
		// 数値を直接書き込み
		b := buf.AvailableBuffer()
		b = strconv.AppendFloat(b, v.Num, 'f', 4, 64)
		buf.Write(b)
		buf.WriteByte('}')

	case TypeString, TypeProperty:
		switch v.Type {
		case TypeString:
			buf.WriteString(`{"Str":`)
		case TypeProperty:
			buf.WriteString(`{"Prop":`)
		}
		// strconv.AppendQuote を使えば、AvailableBuffer に直接
		// エスケープ済みの文字列 ("hello" など) を書き込めます！
		b := buf.AvailableBuffer()
		b = strconv.AppendQuote(b, v.Str)
		buf.Write(b)
		buf.WriteByte('}')

	case TypeBool:
		if v.Bool {
			buf.WriteString(`{"Bool":true}`)
		} else {
			buf.WriteString(`{"Bool":false}`)
		}
	case TypeList, TypeLitList:
		switch v.Type {
		case TypeList:
			buf.WriteString(`{"List":`)
		case TypeLitList:
			buf.WriteString(`{"Lit":`)
		}
		buf.WriteByte('[')
		for i, item := range v.List {
			if i > 0 {
				buf.WriteByte(',')
			}
			// 再帰的に呼び出される
			b, err := json.Marshal(item)
			if err != nil {
				return nil, err
			}
			buf.Write(b)
		}
		buf.WriteString(`]}`)
	default:
		return nil, errors.New("cannot marshal unknown type: " + strconv.Itoa(int(v.Type)))
	}

	return buf.Bytes(), nil
}

// json.Unmarshaler インターフェースの実装
// json.Unmarshalを使わず自前でパースする
func (v *Value) UnmarshalJSON(data []byte) error {
	// まずは {"Num":123} のような形式を想定して、キーと値を分割する
	s := strings.TrimSpace(string(data))
	if !strings.HasPrefix(s, "{") || !strings.HasSuffix(s, "}") {
		return errors.New("invalid JSON format for Value: " + s)
	}
	// {"Key":Value} の形式を想定して、キーと値を分割
	s = strings.TrimPrefix(s, "{")
	s = strings.TrimSuffix(s, "}")
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return errors.New("invalid JSON format for Value: " + s)
	}
	key := strings.TrimSpace(parts[0])
	valueStr := strings.TrimSpace(parts[1])

	switch key {
	case `"Num"`:
		f, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return err
		}
		v.Type = TypeNumber
		v.Num = f

	case `"Str"`:
		str, err := strconv.Unquote(valueStr)
		if err != nil {
			return err
		}
		v.Type = TypeString
		v.Str = str

	case `"Prop"`:
		str, err := strconv.Unquote(valueStr)
		if err != nil {
			return err
		}
		v.Type = TypeProperty
		v.Str = str

	case `"Bool"`:
		b, err := strconv.ParseBool(valueStr)
		if err != nil {
			return err
		}
		v.Type = TypeBool
		v.Bool = b

	case `"List"`, `"Lit"`:
		if key == `"List"` {
			v.Type = TypeList
		} else {
			v.Type = TypeLitList
		}

		if valueStr == "null" {
			v.List = []Value{}
		} else {
			var items []Value
			err := json.Unmarshal([]byte(valueStr), &items)
			if err != nil {
				return err
			}
			v.List = items
		}
	default:
		return errors.New("unknown key in JSON for Value: " + key)
	}

	return nil
}
