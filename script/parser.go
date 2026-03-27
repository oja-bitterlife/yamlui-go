package script

import (
	"errors"
	"strconv"
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
		return parseList(tn)
	case ')':
		return Value{}, errors.New("unexpected ')'")
	case '"':
		// 前後の " を除去して文字列に
		if len(token) < 2 {
			return Value{}, errors.New("invalid string")
		}
		return NewString(string(token[1 : len(token)-1])), nil
	case '@':
		return NewProperty(string(token)), nil
	default:
		// 数値か、それ以外の識別子（コマンド名など）
		s := string(token)
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return NewNumber(f), nil
		}
		return NewString(s), nil
	}
}

// ==================================================
// parseList. ()で囲まれたリストを解析する
func parseList(tn *Tokenizer) (Value, error) {
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
		if len(token) == 1 && token[0] == ')' {
			break
		}

		// 内容を再帰的に解析
		val, err := parseToken(tn, token)
		if err != nil {
			return Value{}, err
		}
		list = append(list, val)
	}
	return NewList(list), nil
}
