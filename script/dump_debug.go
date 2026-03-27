//go:build debug

package script

import "fmt"

func (v Value) Dump() {
	fmt.Printf("%s\n", v.toStr())
}

func (v Value) ToStr() string {
	return v.toStr()
}
func (v Value) toStr() string {
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
			str += child.toStr()
		}
		str += ")"
		return str
	default:
		return "[Unknown]"
	}
}
