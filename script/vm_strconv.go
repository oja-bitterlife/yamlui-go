package script

import "strings"

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

	var sb strings.Builder
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
	var sb strings.Builder
	sb.WriteByte('"')
	for _, c := range s {
		switch c {
		case '\\':
			sb.WriteString(`\\`)
		case '"':
			sb.WriteString(`\"`)
		case '\n':
			sb.WriteString(`\n`)
		case '\r':
			sb.WriteString(`\r`)
		case '\t':
			sb.WriteString(`\t`)
		default:
			sb.WriteRune(c)
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

	var sb strings.Builder
	for i := 1; i < len(s)-1; i++ {
		c := s[i]
		if c == '\\' {
			if i+1 >= len(s)-1 {
				return "", LogErr("invalid escape sequence in quoted string: %s", s)
			}
			next := s[i+1]
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
			sb.WriteByte(c)
		}
	}

	return sb.String(), nil
}

func ToLower(s string) string {
	var sb strings.Builder
	for _, c := range s {
		if 'A' <= c && c <= 'Z' {
			sb.WriteRune(c + ('a' - 'A'))
		} else {
			sb.WriteRune(c)
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
