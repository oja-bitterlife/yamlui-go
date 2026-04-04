package convert

import (
	"bytes"
	"errors"
	"strconv"

	"github.com/oja-bitterlife/yamlui-go/script"
)

// **********************************************************************
// Marshalの実装
// ==================================================
// json.MarshalJSONを使ってしまわないよう、関数名を違うものにしておく
func ToJSON(v script.Value) ([]byte, error) {
	var buf bytes.Buffer

	switch v.Type {
	case script.TypeNumber:
		buf.WriteString(`{"`)
		buf.WriteString(script.TypeNumberStr)
		buf.WriteString(`":`)
		// 数値を直接書き込み
		b := buf.AvailableBuffer()
		b = strconv.AppendFloat(b, v.Num, 'f', 4, 64)
		buf.Write(b)
		buf.WriteByte('}')

	case script.TypeString, script.TypeProperty:
		switch v.Type {
		case script.TypeString:
			buf.WriteString(`{"`)
			buf.WriteString(script.TypeStringStr)
			buf.WriteString(`":`)
		case script.TypeProperty:
			buf.WriteString(`{"`)
			buf.WriteString(script.TypePropertyStr)
			buf.WriteString(`":`)
		}
		// strconv.AppendQuote を使えば、AvailableBuffer に直接
		// エスケープ済みの文字列 ("hello" など) を書き込めます！
		b := buf.AvailableBuffer()
		b = strconv.AppendQuote(b, v.Str)
		buf.Write(b)
		buf.WriteByte('}')

	case script.TypeBool:
		if v.Bool {
			buf.WriteString(`{"`)
			buf.WriteString(script.TypeBoolStr)
			buf.WriteString(`":true}`)
		} else {
			buf.WriteString(`{"`)
			buf.WriteString(script.TypeBoolStr)
			buf.WriteString(`":false}`)
		}
	case script.TypeList, script.TypeLitList:
		switch v.Type {
		case script.TypeList:
			buf.WriteString(`{"`)
			buf.WriteString(script.TypeListStr)
			buf.WriteString(`":`)
		case script.TypeLitList:
			buf.WriteString(`{"`)
			buf.WriteString(script.TypeLitListStr)
			buf.WriteString(`":`)
		}
		buf.WriteByte('[')
		for i, item := range v.List {
			if i > 0 {
				buf.WriteByte(',')
			}
			// 再帰的に呼び出される
			b, err := ToJSON(item)
			if err != nil {
				return nil, err
			}
			buf.Write(b)
		}
		buf.WriteString(`]}`)
	case script.TypeLitMap:
		buf.WriteString(`{"`)
		buf.WriteString(script.TypeLitMapStr)
		buf.WriteString(`":{`)
		i := 0
		for k, v := range v.Map {
			if i > 0 {
				buf.WriteByte(',')
			}
			// キーは文字列としてエスケープして書き込む
			b := buf.AvailableBuffer()
			b = strconv.AppendQuote(b, k)
			buf.Write(b)
			buf.WriteByte(':')
			// 値は再帰的に呼び出される
			b, err := ToJSON(v)
			if err != nil {
				return nil, err
			}
			buf.Write(b)
			i++
		}
		buf.WriteString(`}}`)

	default:
		return nil, errors.New("cannot marshal unknown type: " + v.Type.String())
	}

	return buf.Bytes(), nil
}

// **********************************************************************
// Unmarshalの実装
// ==================================================
// 内部関すの実装から
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

// ==================================================
// Unmarshalの呼び出し口
// json.UnmarshalJSONを使ってしまわないよう、関数名を違うものにしておく
func ValueFromJSON(data []byte) (script.Value, error) {
	data = bytes.TrimSpace(data)

	if len(data) == 0 {
		return script.Value{}, errors.New("invalid JSON format: empty input")
	}

	// まずはブラケットの整合性をチェック
	if bytes.HasPrefix(data, []byte("{")) && !bytes.HasSuffix(data, []byte("}")) ||
		bytes.HasPrefix(data, []byte("[")) && !bytes.HasSuffix(data, []byte("]")) {
		return script.Value{}, errors.New("invalid JSON format: mismatched brackets")
	}

	// JSONの形式に従ってscript.Valueを構築
	val, err := parseValue(data)
	if err != nil {
		return script.Value{}, err
	}

	// 構築されたscript.Valueをvにコピー
	return val, nil
}
