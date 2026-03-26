package script

import (
	"strconv"
	"strings"
)

// Parse は Lisp 形式の文字列を AST (Value) に変換します
func Parse(src string) Value {
	// 1. トークナイズ: カッコを分離して空白で分割
	src = strings.ReplaceAll(src, "(", " ( ")
	src = strings.ReplaceAll(src, ")", " ) ")
	tokens := strings.Fields(src)

	// 2. パース用のルートとスタック
	// 便宜上、全体を包む隠れたリストを root とする
	root := &Value{Type: TypeList}
	stack := []*Value{root}

	for _, t := range tokens {
		current := stack[len(stack)-1]

		switch t {
		case "(":
			// 新しいリストを作成し、現在のリストに追加
			newList := Value{Type: TypeList}
			current.List = append(current.List, newList)
			// 追加した要素のポインタをスタックに積む
			stack = append(stack, &current.List[len(current.List)-1])
		case ")":
			// リストが終了したのでスタックを戻す
			if len(stack) > 1 {
				stack = stack[:len(stack)-1]
			}
		default:
			// 文字列、数値、プロパティなどを判定して追加
			current.List = append(current.List, parseToken(t))
		}
	}

	// 最初にパースされた要素を返す
	if len(root.List) > 0 {
		return root.List[0]
	}
	return Value{Type: TypeNil}
}

// parseToken は単一のトークンを適切な型に変換します
func parseToken(t string) Value {
	// Property: @x, @y など
	if strings.HasPrefix(t, "@") {
		return Value{Type: TypeProperty, Str: t[1:]}
	}
	// Number: 10, -5.5 など
	if f, err := strconv.ParseFloat(t, 64); err == nil {
		return Value{Type: TypeNumber, Number: f}
	}
	// Bool: true, false
	if b, err := strconv.ParseBool(t); err == nil {
		return Value{Type: TypeBool, Bool: b}
	}
	// Symbol: _.add, set など
	return Value{Type: TypeSymbol, Str: t}
}
