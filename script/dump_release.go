//go:build !debug

package script

func (v Value) Dump()         {}
func (v Value) ToStr() string { return "" }
