package script

import "fmt"

type Tokenizer struct {
	src []byte
	pos int
}

func NewTokenizer(src string) *Tokenizer {
	return &Tokenizer{src: []byte(src), pos: 0}
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}
func isDelim(c byte) bool {
	return c == '(' || c == ')' || c == '{' || c == '}' || c == '"'
}

func (tn *Tokenizer) Next() ([]byte, error) {
	// 空白をスキップ
	for tn.pos < len(tn.src) && isSpace(tn.src[tn.pos]) {
		tn.pos++
	}
	if tn.pos >= len(tn.src) {
		return nil, nil
	}

	start := tn.pos
	c := tn.src[tn.pos]

	switch c {
	case '(', ')', '{', '}': // 構造記号
		tn.pos++
		return tn.src[start:tn.pos], nil

	case '"': // 引用符
		tn.pos++ // 開始の"を飛ばす
		// 終了の"を探す
		for {
			if tn.pos >= len(tn.src) {
				// 終了の"が見つからなかった
				return nil, fmt.Errorf("unterminated string")
			}
			if tn.src[tn.pos] == '"' {
				break
			}
			tn.pos++
		}
		tn.pos++ // 終了の"を飛ばす
		return tn.src[start:tn.pos], nil

	default: // 一塊の文字列
		for tn.pos < len(tn.src) && !isSpace(tn.src[tn.pos]) && !isDelim(tn.src[tn.pos]) {
			tn.pos++
		}
		return tn.src[start:tn.pos], nil
	}
}
