package convert

import (
	"strconv"
	"strings"

	"github.com/oja-bitterlife/yamlui-go/script"
)

// **********************************************************************
// スクリプトの構文解析
// ==================================================
// parse. スクリプト全体を解析して、script.Valueのツリー構造を作る
func parse(src string) (script.Value, error) {
	tn := NewTokenizer(src)

	root := script.Value{Type: script.TypeList, List: []script.Value{}}

	// 明示的な ( がなくても、EOFまでトークンを読み続ける
	for {
		token, err := tn.Next()

		// エラー
		if err != nil {
			return script.Value{}, err
		}
		// 空トークンは終了
		if token == nil {
			break
		}

		// トークンを解析して root.List に append していく
		val, err := parseToken(tn, token)
		if err != nil {
			return script.Value{}, err
		}
		root.List = append(root.List, val)
	}

	// もし root.List の中身が複数あるなら、
	// VM側で「順番に実行するもの」として扱う
	return root, nil
}

// ==================================================
// parseToken. トークンを解析してscript.Valueに変換する
func parseToken(tn *Tokenizer, token []byte) (script.Value, error) {
	switch token[0] {
	case '(':
		return parseList(tn, ')')
	case '{':
		return parseList(tn, '}')
	case ')', '}':
		return script.Value{}, script.LogErr("unexpected token: %s", string(token))
	case '"', '\'':
		// 前後の " を除去して文字列に
		if len(token) < 2 {
			return script.Value{}, script.LogErr("invalid string token: %s", string(token))
		}
		return script.NewString(string(token[1 : len(token)-1])), nil
	case '@', '_':
		return script.NewProperty(string(token)), nil
	default:
		// bool, number, or string
		s := string(token)

		// boolチェック
		if strings.ToLower(s) == "true" {
			return script.NewBool(true), nil
		}
		if strings.ToLower(s) == "false" {
			return script.NewBool(false), nil
		}

		// 数値にできる？
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return script.NewNumber(f), nil
		}

		// コマンドかPropety(""で囲まれていない文字列)
		return script.NewString(s), nil
	}
}

// ==================================================
// parseList. ()で囲まれたリストを解析する
func parseList(tn *Tokenizer, terminate byte) (script.Value, error) {
	var list []script.Value
	for {
		token, err := tn.Next()

		// エラーチェック
		if err != nil {
			return script.Value{}, err
		}
		if token == nil {
			return script.Value{}, script.LogErr("unexpected end of input, expected '%c'", terminate)
		}

		// 終端
		if len(token) == 1 && token[0] == terminate {
			break
		}

		// 内容を再帰的に解析
		val, err := parseToken(tn, token)
		if err != nil {
			return script.Value{}, err
		}

		// リテラルブロックはPropertyやListを含められない
		if terminate == '}' {
			switch val.Type {
			case script.TypeNumber, script.TypeString, script.TypeBool:
				// OK
			default:
				return script.Value{}, script.LogErr("invalid value in literal block: only number, string, and bool are allowed, got %s", val.Type.String())
			}
		}

		list = append(list, val)
	}

	// リテラルブロックなら、リスト全体を文字列化してリテラルとして返す
	if terminate == '}' {
		return script.NewLitList(list), nil
	} else {
		return script.NewList(list), nil
	}
}
