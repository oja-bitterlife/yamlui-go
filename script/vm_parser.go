package script

import (
	"errors"
	"strconv"
	"strings"
)

// **********************************************************************
// スクリプトの構文解析
// ==================================================
// parse. スクリプト全体を解析して、Valueのツリー構造を作る
func parse(src string) (Value, error) {
	tn := NewTokenizer(src)

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
		val, err := parseToken(tn, token)
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
// parseToken. トークンを解析してValueに変換する
func parseToken(tn *Tokenizer, token []byte) (Value, error) {
	switch token[0] {
	case '(':
		return parseList(tn, ')')
	case '{':
		return parseList(tn, '}')
	case ')', '}':
		return Value{}, errors.New("unexpected '" + string(token[0]) + "'")
	case '"':
		// 前後の " を除去して文字列に
		if len(token) < 2 {
			return Value{}, errors.New("invalid string")
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
// parseList. ()で囲まれたリストを解析する
func parseList(tn *Tokenizer, terminate byte) (Value, error) {
	var list []Value
	for {
		token, err := tn.Next()

		// エラーチェック
		if err != nil {
			return Value{}, err
		}
		if token == nil {
			return Value{}, errors.New("unclosed parenthesis")
		}

		// 終端
		if len(token) == 1 && token[0] == terminate {
			break
		}

		// 内容を再帰的に解析
		val, err := parseToken(tn, token)
		if err != nil {
			return Value{}, err
		}

		// リテラルブロックはPropertyやListを含められない
		if terminate == '}' {
			switch val.Type {
			case TypeNumber, TypeString, TypeBool:
				// OK
			default:
				return Value{}, errors.New("literal block can only contain number, string, or bool: " + val.Type.ToStr())
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
