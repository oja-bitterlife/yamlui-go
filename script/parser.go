package script

import (
	"fmt"
	"strconv"
)

type Parser struct {
	tn *Tokenizer
}

func NewParser(src string) *Parser {
	return &Parser{tn: NewTokenizer(src)}
}

// Parse はエントリーポイント
func (p *Parser) Parse() (Value, error) {
	token, err := p.tn.Next()
	if err != nil {
		return Value{}, err
	}
	if token == nil {
		return Value{}, nil // 終端
	}
	return p.parseToken(token)
}

// トークンを解析してValueに変換する
func (p *Parser) parseToken(token []byte) (Value, error) {
	switch token[0] {
	case '(':
		return p.parseList()
	case ')':
		return Value{}, fmt.Errorf("unexpected ')'")
	case '"':
		// 前後の " を除去して文字列に
		if len(token) < 2 {
			return Value{}, fmt.Errorf("invalid string")
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

// ()で囲まれたリストを解析する
func (p *Parser) parseList() (Value, error) {
	var list []Value
	for {
		token, err := p.tn.Next()

		// エラーチェック
		if err != nil {
			return Value{}, err
		}
		if token == nil {
			return Value{}, fmt.Errorf("unclosed parenthesis")
		}

		// 終端
		if len(token) == 1 && token[0] == ')' {
			break
		}

		// 内容を再帰的に解析
		val, err := p.parseToken(token)
		if err != nil {
			return Value{}, err
		}
		list = append(list, val)
	}
	return NewList(list), nil
}
