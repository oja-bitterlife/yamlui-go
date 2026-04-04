package script

import (
	"strconv"
	"strings"
)

// **********************************************************************
// 構文解析してValueのASTを作る
func Compile(src string) (Value, error) {
	return parseLisp(src)
}

// **********************************************************************
// スクリプトの構文解析
// ==================================================
// parseLisp. スクリプト全体を解析して、Valueのツリー構造を作る
func parseLisp(src string) (Value, error) {
	tn := NewLispTokenizer(src)

	root := Value{Type: TypeList, List: []Value{}}

	// 明示的な ( がなくても、EOFまでトークンを読み続ける
	for {
		token, err := tn.Next()

		// エラー
		if err != nil {
			return Value{}, err
		}
		// 空トークンは終了
		if token == nil {
			break
		}

		// トークンを解析して root.List に append していく
		val, err := parseLispToken(tn, token)
		if err != nil {
			return Value{}, err
		}
		root.List = append(root.List, val)
	}

	// もし root.List の中身が複数あるなら、
	// VM側で「順番に実行するもの」として扱う
	return root, nil
}

// ==================================================
// parseLispToken. トークンを解析してValueに変換する
func parseLispToken(tn *LispTokenizer, token []byte) (Value, error) {
	switch token[0] {
	case '(':
		return parseLispList(tn, ')')
	case '{':
		return parseLispList(tn, '}')
	case ')', '}':
		return Value{}, LogErr("unexpected token: %s", string(token))
	case '"', '\'':
		// 前後の " を除去して文字列に
		if len(token) < 2 {
			return Value{}, LogErr("invalid string token: %s", string(token))
		}
		return NewString(string(token[1 : len(token)-1])), nil
	case '@', '_':
		return NewProperty(string(token)), nil
	default:
		// bool, number, or string
		s := string(token)

		// boolチェック
		if strings.ToLower(s) == "true" {
			return NewBool(true), nil
		}
		if strings.ToLower(s) == "false" {
			return NewBool(false), nil
		}

		// 数値にできる？
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return NewNumber(f), nil
		}

		// コマンドかPropety(""で囲まれていない文字列)
		return NewString(s), nil
	}
}

// ==================================================
// parseLispList. ()で囲まれたリストを解析する
func parseLispList(tn *LispTokenizer, terminate byte) (Value, error) {
	var list []Value
	for {
		token, err := tn.Next()

		// エラーチェック
		if err != nil {
			return Value{}, err
		}
		if token == nil {
			return Value{}, LogErr("unexpected end of input, expected '%c'", terminate)
		}

		// 終端
		if len(token) == 1 && token[0] == terminate {
			break
		}

		// 内容を再帰的に解析
		val, err := parseLispToken(tn, token)
		if err != nil {
			return Value{}, err
		}

		// リテラルブロックはPropertyやListを含められない
		if terminate == '}' {
			switch val.Type {
			case TypeNumber, TypeString, TypeBool:
				// OK
			default:
				return Value{}, LogErr("invalid value in literal block: only number, string, and bool are allowed, got %s", val.Type.String())
			}
		}

		list = append(list, val)
	}

	// リテラルブロックなら、リスト全体を文字列化してリテラルとして返す
	if terminate == '}' {
		return NewLitList(list), nil
	} else {
		return NewList(list), nil
	}
}
