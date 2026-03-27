//go:build debug

package script

import "fmt"

func (v Value) Dump() {
	v._dump(0)
}
func (v Value) _dump(indent int) {
	prefix := ""
	for i := 0; i < indent; i++ {
		prefix += "  " // 階層ごとにスペース2つ
	}

	switch v.Type {
	case TypeNumber:
		fmt.Printf("%s[Number] %v\n", prefix, v.Num)
	case TypeString:
		fmt.Printf("%s[String] %q\n", prefix, v.Str)
	case TypeProperty:
		fmt.Printf("%s[Property] %s\n", prefix, v.Str)
	case TypeList:
		fmt.Printf("%s[List] (\n", prefix)
		for _, child := range v.List {
			child._dump(indent + 1) // 再帰でインデントを増やす
		}
		fmt.Printf("%s)\n", prefix)
	case TypeQuote:
		fmt.Printf("%s[Quote] (\n", prefix)
		for _, child := range v.List {
			child._dump(indent + 1) // 再帰でインデントを増やす
		}
		fmt.Printf("%s)\n", prefix)
	default:
		fmt.Printf("%s[Unknown]\n", prefix)
	}
}

func (v Value) ToStr() string {
	return v._toStr()
}
func (v Value) _toStr() string {
	switch v.Type {
	case TypeNumber:
		return fmt.Sprintf("[Number] %v", v.Num)
	case TypeString:
		return fmt.Sprintf("[String] %q", v.Str)
	case TypeProperty:
		return fmt.Sprintf("[Property] %s", v.Str)
	case TypeList:
		str := "[List] ("
		for _, child := range v.List {
			str += child._toStr() // 再帰でインデントを増やす
		}
		str += ")"
		return str
	case TypeQuote:
		str := "[Quote] ("
		for _, child := range v.List {
			str += child._toStr()
		}
		str += ")"
		return str
	default:
		return "[Unknown]"
	}
}
