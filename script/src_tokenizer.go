package script

// 入力文字列をトークンに分割する
type LispTokenizer struct {
	src []byte
	pos int
}

func NewLispTokenizer(src string) *LispTokenizer {
	return &LispTokenizer{src: []byte(src), pos: 0}
}

// **********************************************************************
// トークン分割
// ==================================================
// 特別な文字
func isLispSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}
func isLispDelim(c byte) bool {
	return c == '(' || c == ')' || c == '{' || c == '}' || c == '"' || c == '\''
}

// ==================================================
// トークンを一つずつ返す。最後まで行ってたらnilを返す
func (tn *LispTokenizer) Next() ([]byte, error) {
	// 空白をスキップ
	for tn.pos < len(tn.src) && isLispSpace(tn.src[tn.pos]) {
		tn.pos++
	}
	if tn.pos >= len(tn.src) {
		return nil, nil
	}

	// トークンの開始位置を記録
	start := tn.pos
	c := tn.src[tn.pos]

	switch c {
	case '#': // コメント
		tn.pos++ // '#' をスキップ
		for tn.pos < len(tn.src) && tn.src[tn.pos] != '\n' {
			tn.pos++ // 行末までスキップ
		}
		return tn.Next() // コメントをスキップして次のトークンを取得

	case '(', ')', '{', '}': // 構造記号
		tn.pos++
		return tn.src[start:tn.pos], nil

	case '"', '\'': // 引用符
		quote := c // 開始したクォート文字を覚えておく
		tn.pos++   // 開始の"を飛ばす

		// 終了の"を探す
		for tn.pos < len(tn.src) {
			// エスケープ処理
			if tn.src[tn.pos] == '\\' {
				tn.pos++ // '\' をスキップ
				if tn.pos < len(tn.src) {
					tn.pos++ // エスケープされた次の文字（' や "）も無条件でスキップ
				}
				continue
			}

			// 閉じクォートを発見
			if tn.src[tn.pos] == quote {
				tn.pos++ // 閉じクォートを含める
				return tn.src[start:tn.pos], nil
			}

			tn.pos++
		}
		return tn.src[start:tn.pos], nil

	default: // 一塊の文字列
		for tn.pos < len(tn.src) && !isLispSpace(tn.src[tn.pos]) && !isLispDelim(tn.src[tn.pos]) {
			tn.pos++
		}
		return tn.src[start:tn.pos], nil
	}
}
