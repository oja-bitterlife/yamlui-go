package script

import (
	"bytes"
)

// **********************************************************************
// Unmarshalの実装
// json.UnmarshalJSONを使ってしまわないよう、関数名を違うものにしておく
func NewFromValueJSON(data []byte) (Value, error) {
	data = bytes.TrimSpace(data)

	if len(data) == 0 {
		return Value{}, LogErr("empty JSON data")
	}

	// まずはブラケットの整合性をチェック
	if bytes.HasPrefix(data, []byte("{")) && !bytes.HasSuffix(data, []byte("}")) ||
		bytes.HasPrefix(data, []byte("[")) && !bytes.HasSuffix(data, []byte("]")) {
		return Value{}, LogErr("mismatched brackets in JSON data")
	}

	// JSONの形式に従ってscript.Valueを構築
	val, err := parseValue(data)
	if err != nil {
		return Value{}, err
	}

	// 構築されたscript.Valueをvにコピー
	return val, nil
}

// **********************************************************************
// Marshalの実装
// json.MarshalJSONを使ってしまわないよう、関数名を違うものにしておく
func (v Value) ToValueJSON() ([]byte, error) {
	var buf bytes.Buffer

	switch v.Type {
	case TypeNumber:
		buf.WriteString(`{"`)
		buf.WriteString(TypeNumberStr)
		buf.WriteString(`":`)
		buf.WriteString(Itoa(v.Num))
		buf.WriteByte('}')

	case TypeString, TypeProperty:
		switch v.Type {
		case TypeString:
			buf.WriteString(`{"`)
			buf.WriteString(TypeStringStr)
			buf.WriteString(`":`)
		case TypeProperty:
			buf.WriteString(`{"`)
			buf.WriteString(TypePropertyStr)
			buf.WriteString(`":`)
		}
		buf.WriteString(Quote(v.Str))
		buf.WriteByte('}')

	case TypeBool:
		if v.Bool() {
			buf.WriteString(`{"`)
			buf.WriteString(TypeBoolStr)
			buf.WriteString(`":true}`)
		} else {
			buf.WriteString(`{"`)
			buf.WriteString(TypeBoolStr)
			buf.WriteString(`":false}`)
		}
	case TypeList, TypeLitList:
		switch v.Type {
		case TypeList:
			buf.WriteString(`{"`)
			buf.WriteString(TypeListStr)
			buf.WriteString(`":`)
		case TypeLitList:
			buf.WriteString(`{"`)
			buf.WriteString(TypeLitListStr)
			buf.WriteString(`":`)
		}
		buf.WriteByte('[')
		for i, item := range v.List {
			if i > 0 {
				buf.WriteByte(',')
			}
			// 再帰的に呼び出される
			b, err := item.ToValueJSON()
			if err != nil {
				return nil, err
			}
			buf.Write(b)
		}
		buf.WriteString(`]}`)
	case TypeLitMap:
		buf.WriteString(`{"`)
		buf.WriteString(TypeLitMapStr)
		buf.WriteString(`":{`)
		i := 0
		for k, v := range v.Map {
			if i > 0 {
				buf.WriteByte(',')
			}
			// キーは文字列としてエスケープして書き込む
			buf.WriteString(Quote(k))
			buf.WriteByte(':')
			// 値は再帰的に呼び出される
			b, err := v.ToValueJSON()
			if err != nil {
				return nil, err
			}
			buf.Write(b)
			i++
		}
		buf.WriteString(`}}`)

	default:
		return nil, LogErr("unsupported Value type: %v", v.Type)
	}

	return buf.Bytes(), nil
}

// **********************************************************************
// Unmarshalの実装のヘルパー
// ==================================================
// parseValue. JSON形式(つまり文字列)のValueをValueに変換する
func parseValue(data []byte) (Value, error) {
	data = bytes.TrimSpace(data)

	// まずはブラケットを取り除く
	data = bytes.TrimPrefix(data, []byte("{"))
	data = bytes.TrimSuffix(data, []byte("}"))

	// キーと値を分割
	parts := bytes.SplitN(data, []byte(":"), 2)
	if len(parts) != 2 {
		return Value{}, LogErr("invalid JSON format: expected key-value pair, got \"%s\"", string(data))
	}

	key, err := Unquote(TrimSpace(string(parts[0])))
	if err != nil {
		return Value{}, LogErr("invalid JSON format: invalid key \"%s\": %v", string(parts[0]), err)
	}
	value := bytes.TrimSpace(parts[1])

	switch key {
	case TypeNumberStr:
		f, err := Atoi(string(value))
		if err != nil {
			return Value{}, LogErr("invalid JSON format: invalid number \"%s\": %v", string(value), err)
		}
		return NewNumber(f), nil
	case TypeStringStr:
		s, err := Unquote(string(value))
		if err != nil {
			return Value{}, LogErr("invalid JSON format: invalid string \"%s\": %v", string(value), err)
		}
		return NewString(s), nil
	case TypePropertyStr:
		s, err := Unquote(string(value))
		if err != nil {
			return Value{}, LogErr("invalid JSON format: invalid property \"%s\": %v", string(value), err)
		}
		return NewProperty(s), nil
	case TypeBoolStr:
		b, err := AtoBool(string(value))
		if err != nil {
			return Value{}, LogErr("invalid JSON format: invalid bool \"%s\": %v", string(value), err)
		}
		return NewBool(b), nil
	case TypeListStr:
		list, err := parseValueList(value)
		if err != nil {
			return Value{}, LogErr("invalid JSON format: invalid list \"%s\": %v", string(value), err)
		}
		return NewList(list), nil
	case TypeLitListStr:
		list, err := parseValueList(value)
		if err != nil {
			return Value{}, LogErr("invalid JSON format: invalid literal list \"%s\": %v", string(value), err)
		}
		return NewLitList(list), nil
	case TypeLitMapStr:
		m, err := parseValueMap(value)
		if err != nil {
			return Value{}, LogErr("invalid JSON format: invalid literal map \"%s\": %v", string(value), err)
		}
		return NewLitMap(m), nil
	default:
		return Value{}, LogErr("invalid JSON format: unknown type \"%s\"", key)
	}
}

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
func parseValueList(data []byte) ([]Value, error) {
	// 中身はValu型なので、{"Num":1} のような形式で入っているはず
	data = bytes.TrimSpace(data)

	// まずはブラケットを取り除く
	data = bytes.TrimPrefix(data, []byte("["))
	data = bytes.TrimSuffix(data, []byte("]"))

	var list []Value

	// Listの各要素は{}のペアなので、{}ごとに処理をしていく
	current := 0
	for current < len(data) {
		// 最初の{}のペアを見つける
		start, end := findStartEnd(data, current, '{', '}')
		if start == -1 || end == -1 {
			return nil, LogErr("invalid JSON format: expected list of objects, got \"%s\"", string(data))
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
func parseValueMap(data []byte) (map[string]Value, error) {
	// 中身はValu型なので、{"key":{"Num":1}} のような形式で入っているはず
	data = bytes.TrimSpace(data)

	// まずはブラケットを取り除く
	data = bytes.TrimPrefix(data, []byte("{"))
	data = bytes.TrimSuffix(data, []byte("}"))

	m := make(map[string]Value)

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

		key, err := Unquote(string(data[keyStart:keyEnd]))
		if err != nil {
			return nil, LogErr("invalid JSON format: invalid key at position %d: %v", keyStart, err)
		}

		// コロンを見つける
		// ----------------------------------------
		colonIndex := bytes.IndexByte(data[keyEnd:], ':')
		if colonIndex == -1 {
			return nil, LogErr("invalid JSON format: expected ':' after key \"%s\"", key)
		}

		// 値を見つける
		// ----------------------------------------
		valueStart, valueEnd := findStartEnd(data, keyEnd+colonIndex, '{', '}')
		if valueStart == -1 || valueEnd == -1 {
			return nil, LogErr("invalid JSON format: expected value object after key \"%s\"", key)
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
