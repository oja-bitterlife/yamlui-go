package convert

import (
	"bytes"
	"strconv"

	"github.com/oja-bitterlife/yamlui-go/script"
)

// **********************************************************************
// 構文解析してscript.ValueのASTを作る
func Compile(src string) (script.Value, error) {
	return parse(src)
}

// **********************************************************************
// Marshalの実装
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
		return nil, script.LogErr("unsupported Value type: %v", v.Type)
	}

	return buf.Bytes(), nil
}

// **********************************************************************
// Unmarshalの実装
// json.UnmarshalJSONを使ってしまわないよう、関数名を違うものにしておく
func FromJSON(data []byte) (script.Value, error) {
	data = bytes.TrimSpace(data)

	if len(data) == 0 {
		return script.Value{}, script.LogErr("empty JSON data")
	}

	// まずはブラケットの整合性をチェック
	if bytes.HasPrefix(data, []byte("{")) && !bytes.HasSuffix(data, []byte("}")) ||
		bytes.HasPrefix(data, []byte("[")) && !bytes.HasSuffix(data, []byte("]")) {
		return script.Value{}, script.LogErr("mismatched brackets in JSON data")
	}

	// JSONの形式に従ってscript.Valueを構築
	val, err := parseValue(data)
	if err != nil {
		return script.Value{}, err
	}

	// 構築されたscript.Valueをvにコピー
	return val, nil
}
