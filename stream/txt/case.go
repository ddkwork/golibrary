package txt

import (
	"unicode"
	"unicode/utf8"
)

func ToCamelCase(in string) string {
	runes := []rune(in)
	out := make([]rune, 0, len(runes))
	up := true
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == '_' {
			up = true
		} else {
			if up {
				r = unicode.ToUpper(r)
				up = false
			}
			out = append(out, r)
		}
	}
	return string(out)
}

func ToCamelCaseWithExceptions(in string, exceptions *AllCaps) string {
	out := ToCamelCase(in)
	pos := 0
	runes := []rune(out)
	rr := RuneReader{}
	for {
		rr.Src = runes[pos:]
		rr.Pos = 0
		matches := exceptions.regex.FindReaderIndex(&rr)
		if len(matches) == 0 {
			break
		}
		for i := matches[0] + 1; i < matches[1]; i++ {
			runes[pos+i] = unicode.ToUpper(runes[pos+i])
		}
		pos += matches[0] + 1
	}
	return string(runes)
}

func ToSnakeCase(in string) string {
	runes := []rune(in)
	out := make([]rune, 0, 1+len(runes))
	for i := 0; i < len(runes); i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < len(runes) && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}
	return string(out)
}

func FirstToUpper(in string) string {
	if in == "" {
		return in
	}
	r, size := utf8.DecodeRuneInString(in)
	if r == utf8.RuneError {
		return in
	}
	return string(unicode.ToUpper(r)) + in[size:]
}

func FirstToLower(in string) string {
	if in == "" {
		return in
	}
	r, size := utf8.DecodeRuneInString(in)
	if r == utf8.RuneError {
		return in
	}
	return string(unicode.ToLower(r)) + in[size:]
}
