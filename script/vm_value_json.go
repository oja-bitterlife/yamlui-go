package script

import (
	"bytes"
	"errors"
	"strconv"
)

// **********************************************************************
// Marshal/Unmarshal
// ==================================================
// json.Marshaler インターフェースの実装
func (v Value) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	switch v.Type {
	case TypeNumber:
		buf.WriteString(`{"`)
		buf.WriteString(TypeNumberStr)
		buf.WriteString(`":`)
		// 数値を直接書き込み
		b := buf.AvailableBuffer()
		b = strconv.AppendFloat(b, v.Num, 'f', 4, 64)
		buf.Write(b)
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
		// strconv.AppendQuote を使えば、AvailableBuffer に直接
		// エスケープ済みの文字列 ("hello" など) を書き込めます！
		b := buf.AvailableBuffer()
		b = strconv.AppendQuote(b, v.Str)
		buf.Write(b)
		buf.WriteByte('}')

	case TypeBool:
		if v.Bool {
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
			b, err := item.MarshalJSON()
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
			b := buf.AvailableBuffer()
			b = strconv.AppendQuote(b, k)
			buf.Write(b)
			buf.WriteByte(':')
			// 値は再帰的に呼び出される
			b, err := v.MarshalJSON()
			if err != nil {
				return nil, err
			}
			buf.Write(b)
			i++
		}
		buf.WriteString(`}}`)

	default:
		return nil, errors.New("cannot marshal unknown type: " + strconv.Itoa(int(v.Type)))
	}

	return buf.Bytes(), nil
}

// ==================================================
// json.Unmarshaler インターフェースの実装
func (v *Value) parseValue(data []byte) (Value, error) {
	data = bytes.TrimSpace(data)

	// まずはブラケットを取り除く
	data = bytes.TrimPrefix(data, []byte("{"))
	data = bytes.TrimSuffix(data, []byte("}"))

	// キーと値を分割
	parts := bytes.SplitN(data, []byte(":"), 2)
	if len(parts) != 2 {
		return Value{}, errors.New("invalid JSON format: expected key-value pair")
	}

	key, err := strconv.Unquote(string(bytes.TrimSpace(parts[0])))
	if err != nil {
		return Value{}, err
	}
	value := bytes.TrimSpace(parts[1])

	switch string(key) {
	case TypeNumberStr:
		f, err := strconv.ParseFloat(string(value), 64)
		if err != nil {
			return Value{}, err
		}
		return NewNumber(f), nil
	case TypeStringStr:
		s, err := strconv.Unquote(string(value))
		if err != nil {
			return Value{}, err
		}
		return NewString(s), nil
	case TypePropertyStr:
		s, err := strconv.Unquote(string(value))
		if err != nil {
			return Value{}, err
		}
		return NewProperty(s), nil
	case TypeBoolStr:
		b, err := strconv.ParseBool(string(value))
		if err != nil {
			return Value{}, err
		}
		return NewBool(b), nil
	case TypeListStr:
		list, err := v.parseList(value)
		if err != nil {
			return Value{}, err
		}
		return NewList(list), nil
	case TypeLitListStr:
		list, err := v.parseList(value)
		if err != nil {
			return Value{}, err
		}
		return NewLitList(list), nil
	case TypeLitMapStr:
		m, err := v.parseMap(value)
		if err != nil {
			return Value{}, err
		}
		return NewLitMap(m), nil
	default:
		return Value{}, errors.New("invalid JSON format: unknown key \"" + string(key) + "\"")
	}
}

func (v *Value) findStartEnd(data []byte, start int, openChar, closeChar byte) (int, int) {
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

func (v *Value) parseList(data []byte) ([]Value, error) {
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
		start, end := v.findStartEnd(data, current, '{', '}')
		if start == -1 || end == -1 {
			return nil, errors.New("invalid JSON format: unmatched brackets in list")
		}

		// {}の中身を解析
		val, err := v.parseValue(data[start:end])
		if err != nil {
			return nil, err
		}
		list = append(list, val)

		// 次の{}を探すために、現在の位置を更新
		current = end
	}

	return list, nil
}
func (v *Value) parseMap(data []byte) (map[string]Value, error) {
	// 中身はValu型なので、{"key":{"Num":1}} のような形式で入っているはず
	data = bytes.TrimSpace(data)

	// まずはブラケットを取り除く
	data = bytes.TrimPrefix(data, []byte("{"))
	data = bytes.TrimSuffix(data, []byte("}"))

	m := NewValueMap()

	// キーと値のペアを処理
	current := 0
	for current < len(data) {
		// キーを見つける
		// シングルクオートはJSONでは使えないらしいのでダブルクオートのみ
		// ----------------------------------------
		keyStart, keyEnd := v.findStartEnd(data, current, '"', '"')
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
		valueStart, valueEnd := v.findStartEnd(data, keyEnd+colonIndex, '{', '}')
		if valueStart == -1 || valueEnd == -1 {
			return nil, errors.New("invalid JSON format: unmatched brackets in value")
		}

		val, err := v.parseValue(data[valueStart:valueEnd])
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

func (v *Value) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)

	if len(data) == 0 {
		return errors.New("empty JSON data")
	}

	// まずはブラケットの整合性をチェック
	if bytes.HasPrefix(data, []byte("{")) && !bytes.HasSuffix(data, []byte("}")) ||
		bytes.HasPrefix(data, []byte("[")) && !bytes.HasSuffix(data, []byte("]")) {
		return errors.New("invalid JSON format: mismatched brackets")
	}

	// JSONの形式に従ってValueを構築
	val, err := v.parseValue(data)
	if err != nil {
		return err
	}

	// 構築されたValueをvにコピー
	*v = val
	return nil
}
