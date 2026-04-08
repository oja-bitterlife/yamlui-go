package script

// **********************************************************************
// strings的な機能の自前実装
// ==================================================
// strings.Builderのような
type strBuilder struct {
	buf []byte
}

func NewStrBuilder(capacity int) *strBuilder {
	return &strBuilder{
		buf: make([]byte, 0, capacity),
	}
}

func (sb *strBuilder) WriteByte(b byte) error {
	sb.buf = append(sb.buf, b)
	return nil
}

func (sb *strBuilder) String() string {
	return string(sb.buf)
}

// ==================================================
// Prefix/Suffix
func HasPrefix(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	return s[:len(prefix)] == prefix
}

func HasSuffix(s, suffix string) bool {
	if len(s) < len(suffix) {
		return false
	}
	return s[len(s)-len(suffix):] == suffix
}

func TrimPrefix(s, prefix string) string {
	if !HasPrefix(s, prefix) {
		return s
	}
	return s[len(prefix):]
}

func TrimSuffix(s, suffix string) string {
	if !HasSuffix(s, suffix) {
		return s
	}
	return s[:len(s)-len(suffix)]
}

// **********************************************************************
// strconv的な機能の自前実装
// ==================================================
// 文字列の前後のスペースを削除する
func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func TrimSpace(s string) string {
	start := 0
	for start < len(s) && isSpace(s[start]) {
		start++
	}

	end := len(s) - 1
	for end >= 0 && isSpace(s[end]) {
		end--
	}

	if start > end {
		return ""
	}
	return s[start : end+1]
}

// ==================================================
// Int
func Itoa(num int) string {
	if num == 0 {
		return "0"
	}

	sb := NewStrBuilder(11) // intの最大値は10桁 + 符号

	if num < 0 {
		sb.WriteByte('-')
		num = -num
	}

	digits := []byte{}
	for num > 0 {
		digit := byte(num % 10)
		digits = append(digits, '0'+digit)
		num /= 10
	}

	// digitsは逆順なので、正しい順序でsbに追加する
	for i := len(digits) - 1; i >= 0; i-- {
		sb.WriteByte(digits[i])
	}

	return sb.String()
}

func IsNumeric(s string) bool {
	s = TrimSpace(s)
	if s == "" {
		return false
	}

	if s[0] == '-' || s[0] == '+' {
		s = s[1:]
	}

	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return false
		}
	}

	return true
}

func Atoi(s string) (int, error) {
	s = TrimSpace(s)
	if s == "" {
		return 0, LogErr("imput string is empty")
	}

	sign := 1
	switch s[0] {
	case '-':
		sign = -1
		s = s[1:]
	case '+':
		s = s[1:]
	}

	num := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return 0, LogErr("invalid number: %s", s)
		}
		num = num*10 + int(c-'0')
	}

	return sign * num, nil
}

// ==================================================
// string
// Quoteは、文字列をエスケープしてダブルクオートで囲む関数
func Quote(s string) string {
	sb := NewStrBuilder(len(s) + 8) // エスケープで増える分を考慮して少し余裕を持たせる

	sb.WriteByte('"')
	bytes := []byte(s)
	for _, b := range bytes {
		switch b {
		case '\\':
			sb.WriteByte('\\')
			sb.WriteByte('\\')
		case '"':
			sb.WriteByte('\\')
			sb.WriteByte('"')
		case '\n':
			sb.WriteByte('\\')
			sb.WriteByte('n')
		case '\r':
			sb.WriteByte('\\')
			sb.WriteByte('r')
		case '\t':
			sb.WriteByte('\\')
			sb.WriteByte('t')
		default:
			sb.WriteByte(b)
		}
	}
	sb.WriteByte('"')

	return sb.String()
}

// Unquoteは、Quoteでエスケープされた文字列を元に戻す関数
func Unquote(s string) (string, error) {
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return "", LogErr("invalid quoted string: %s", s)
	}

	sb := NewStrBuilder(len(s)) // 元の長さと同じか小さくなる

	bytes := []byte(s)
	for i := 1; i < len(bytes)-1; i++ {
		b := bytes[i]
		if b == '\\' {
			if i+1 >= len(bytes)-1 {
				return "", LogErr("invalid escape sequence in quoted string: %s", s)
			}
			next := bytes[i+1]
			switch next {
			case '\\':
				sb.WriteByte('\\')
			case '"':
				sb.WriteByte('"')
			case 'n':
				sb.WriteByte('\n')
			case 'r':
				sb.WriteByte('\r')
			case 't':
				sb.WriteByte('\t')
			default:
				return "", LogErr("invalid escape sequence in quoted string: \\%c", next)
			}
			i++
		} else {
			sb.WriteByte(b)
		}
	}

	return sb.String(), nil
}

// ToLowerは、ASCIIの大文字を小文字に変換する関数
func ToLower(s string) string {
	sb := NewStrBuilder(len(s))

	bytes := []byte(s)
	for _, b := range bytes {
		if 'A' <= b && b <= 'Z' {
			sb.WriteByte(b + ('a' - 'A'))
		} else {
			sb.WriteByte(b)
		}
	}

	return sb.String()
}

// ToUpperは、ASCIIの小文字を大文字に変換する関数
func ToUpper(s string) string {
	sb := NewStrBuilder(len(s))

	bytes := []byte(s)
	for _, b := range bytes {
		if 'a' <= b && b <= 'z' {
			sb.WriteByte(b - ('a' - 'A'))
		} else {
			sb.WriteByte(b)
		}
	}

	return sb.String()
}

// ==================================================
// bool
func AtoBool(s string) (bool, error) {
	s = ToLower(TrimSpace(s))
	switch s {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, LogErr("invalid boolean string: %s", s)
	}
}
