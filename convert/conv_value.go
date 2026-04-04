package convert

import (
	"bytes"
	"errors"
	"strconv"

	"github.com/oja-bitterlife/yamlui-go/script"
)

// **********************************************************************
// ValueのJSON形式からscript.Valueへの変換
// ==================================================
// parseValue. JSON形式(つまり文字列)のValueをscript.Valueに変換する
func parseValue(data []byte) (script.Value, error) {
	data = bytes.TrimSpace(data)

	// まずはブラケットを取り除く
	data = bytes.TrimPrefix(data, []byte("{"))
	data = bytes.TrimSuffix(data, []byte("}"))

	// キーと値を分割
	parts := bytes.SplitN(data, []byte(":"), 2)
	if len(parts) != 2 {
		return script.Value{}, errors.New("invalid JSON format: expected key-value pair")
	}

	key, err := strconv.Unquote(string(bytes.TrimSpace(parts[0])))
	if err != nil {
		return script.Value{}, err
	}
	value := bytes.TrimSpace(parts[1])

	switch string(key) {
	case script.TypeNumberStr:
		f, err := strconv.ParseFloat(string(value), 64)
		if err != nil {
			return script.Value{}, err
		}
		return script.NewNumber(f), nil
	case script.TypeStringStr:
		s, err := strconv.Unquote(string(value))
		if err != nil {
			return script.Value{}, err
		}
		return script.NewString(s), nil
	case script.TypePropertyStr:
		s, err := strconv.Unquote(string(value))
		if err != nil {
			return script.Value{}, err
		}
		return script.NewProperty(s), nil
	case script.TypeBoolStr:
		b, err := strconv.ParseBool(string(value))
		if err != nil {
			return script.Value{}, err
		}
		return script.NewBool(b), nil
	case script.TypeListStr:
		list, err := parseValueList(value)
		if err != nil {
			return script.Value{}, err
		}
		return script.NewList(list), nil
	case script.TypeLitListStr:
		list, err := parseValueList(value)
		if err != nil {
			return script.Value{}, err
		}
		return script.NewLitList(list), nil
	case script.TypeLitMapStr:
		m, err := parseValueMap(value)
		if err != nil {
			return script.Value{}, err
		}
		return script.NewLitMap(m), nil
	default:
		return script.Value{}, errors.New("invalid JSON format: unknown key \"" + string(key) + "\"")
	}
}

// **********************************************************************
// ヘルパー
// ==================================================
// findStartEnd. dataのstart位置から、openCharで始まりcloseCharで終わる部分の開始位置と終了位置を返す
func findStartEnd(data []byte, start int, openChar, closeChar byte) (int, int) {
	depth := 0
	inString := false
	escaped := false

	// 文字列は入れ子にならない。オブジェクトのときだけ入れ子対応
	isObject := openChar != '"'
	if isObject {
		depth = 0
	}

	// 最初の開き文字を見つける
	offset := bytes.IndexByte(data[start:], openChar)
	if offset == -1 {
		return -1, -1
	}
	start += offset

	// 開き文字を見つけた位置から閉じる文字を探す
	for i := start; i < len(data); i++ {
		b := data[i]

		// 文字列内のエスケープ処理
		if inString {
			if escaped {
				escaped = false
				continue
			}
			if b == '\\' {
				escaped = true
				continue
			}
			if b == '"' {
				inString = false
				// 文字列単体を探していた場合（キーのパースなど）
				if !isObject {
					return start, i + 1
				}
			}
			continue
		}

		// 文字列外の処理
		switch b {
		case '"':
			inString = true
		case openChar:
			if isObject {
				depth++
			}
		case closeChar:
			if isObject {
				depth--
				if depth == 0 {
					return start, i + 1
				}
			}
		}
	}
	return -1, -1
}

// ==================================================
// List部分のパース
func parseValueList(data []byte) ([]script.Value, error) {
	// 中身はValu型なので、{"Num":1} のような形式で入っているはず
	data = bytes.TrimSpace(data)

	// まずはブラケットを取り除く
	data = bytes.TrimPrefix(data, []byte("["))
	data = bytes.TrimSuffix(data, []byte("]"))

	var list []script.Value

	// Listの各要素は{}のペアなので、{}ごとに処理をしていく
	current := 0
	for current < len(data) {
		// 最初の{}のペアを見つける
		start, end := findStartEnd(data, current, '{', '}')
		if start == -1 || end == -1 {
			return nil, errors.New("invalid JSON format: unmatched brackets in list")
		}

		// {}の中身を解析
		val, err := parseValue(data[start:end])
		if err != nil {
			return nil, err
		}
		list = append(list, val)

		// 次の{}を探すために、現在の位置を更新
		current = end
	}

	return list, nil
}

// ==================================================
// Map部分のパース
func parseValueMap(data []byte) (map[string]script.Value, error) {
	// 中身はValu型なので、{"key":{"Num":1}} のような形式で入っているはず
	data = bytes.TrimSpace(data)

	// まずはブラケットを取り除く
	data = bytes.TrimPrefix(data, []byte("{"))
	data = bytes.TrimSuffix(data, []byte("}"))

	m := make(map[string]script.Value)

	// キーと値のペアを処理
	current := 0
	for current < len(data) {
		// キーを見つける
		// シングルクオートはJSONでは使えないらしいのでダブルクオートのみ
		// ----------------------------------------
		keyStart, keyEnd := findStartEnd(data, current, '"', '"')
		if keyStart == -1 || keyEnd == -1 {
			break
		}

		key, err := strconv.Unquote(string(data[keyStart:keyEnd]))
		if err != nil {
			return nil, err
		}

		// コロンを見つける
		// ----------------------------------------
		colonIndex := bytes.IndexByte(data[keyEnd:], ':')
		if colonIndex == -1 {
			return nil, errors.New("invalid JSON format: missing colon after key")
		}

		// 値を見つける
		// ----------------------------------------
		valueStart, valueEnd := findStartEnd(data, keyEnd+colonIndex, '{', '}')
		if valueStart == -1 || valueEnd == -1 {
			return nil, errors.New("invalid JSON format: unmatched brackets in value")
		}

		val, err := parseValue(data[valueStart:valueEnd])
		if err != nil {
			return nil, err
		}

		// マップに追加して次へ
		// ----------------------------------------
		m[key] = val
		current = valueEnd
	}

	return m, nil
}
