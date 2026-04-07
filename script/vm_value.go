package script

// **********************************************************************
// ValueのType定義
// ==================================================
type ValueType int

const (
	// 評価済み
	TypeNil ValueType = iota
	TypeNumber
	TypeBool
	TypeString
	TypeLitList // 評価済みリスト
	TypeLitMap  // 評価済みマップ

	// 未評価
	TypeProperty // @i などの参照
	TypeList     // 未評価リスト
)

const (
	TypeNilStr      = "Nil"
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
	case TypeNil:
		return TypeNilStr
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
		return "Unknown:" + Itoa(int(t))
	}
}

// **********************************************************************
// Valueの定義
// Listに格納する型付きの値。リフレクションを避けるため全部入り
type Value struct {
	// 型情報
	Type ValueType

	// 基本データ
	Num int
	Str string

	// リスト構造 (S式や、分解済みの文字列フラグメント)
	List []Value
	Map  ValueMap
}

// クローンを作成
func (v Value) Clone() Value {
	// 参照以外はこれでOk
	res := v

	// 参照は再起でクローンを作成
	switch v.Type {
	case TypeList, TypeLitList:
		if v.List != nil {
			newList := make([]Value, len(v.List))
			for i, item := range v.List {
				newList[i] = item.Clone() // 再帰
			}
			res.List = newList
		}
	case TypeLitMap:
		if v.Map != nil {
			newMap := make(map[string]Value, len(v.Map))
			for k, val := range v.Map {
				newMap[k] = val.Clone() // 再帰
			}
			res.Map = newMap
		}
	}

	// 複製された値を返す
	return res
}

// BoolはNumに統合したので関数で取得する
func (v Value) Bool() bool {
	return v.Num != 0
}

// ==================================================
// 値の生成関数
func NewNil() Value         { return Value{Type: TypeNil} }
func NewNumber(i int) Value { return Value{Type: TypeNumber, Num: i} }
func NewBool(b bool) Value {
	if b {
		return Value{Type: TypeBool, Num: 1}
	} else {
		return Value{Type: TypeBool, Num: 0}
	}
}
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
		return NewBool(ToLower(v.Str) == "true")
	case TypeLitList, TypeList:
		return NewBool(len(v.List) > 0)
	case TypeLitMap:
		return NewBool(len(v.Map) > 0)
	default:
		LogWarn("cannot convert %s to bool", v.Type.String())
		return NewBool(false)
	}
}

func (v Value) ConvertNumber() Value {
	switch v.Type {
	case TypeBool:
		return NewNumber(v.Num)
	case TypeNumber:
		return v
	case TypeString:
		i, err := Atoi(v.Str)
		if err != nil {
			return NewNumber(0)
		}
		return NewNumber(i)
	case TypeLitList, TypeList:
		return NewNumber(len(v.List))
	case TypeLitMap:
		return NewNumber(len(v.Map))
	default:
		LogWarn("cannot convert %s to number", v.Type.String())
		return NewNumber(0)
	}
}

func (v Value) ConvertString() Value {
	switch v.Type {
	case TypeBool:
		if v.Num != 0 {
			return NewString("true")
		}
		return NewString("false")
	case TypeNumber:
		return NewString(Itoa(v.Num))
	case TypeString:
		return v
	default:
		return NewString("")
	}
}
