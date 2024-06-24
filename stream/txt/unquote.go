package txt

import (
	"unicode/utf8"
)

func UnquoteBytes(text []byte) []byte {
	if len(text) > 1 {
		if ch, _ := utf8.DecodeRune(text); ch == '"' {
			if ch, _ = utf8.DecodeLastRune(text); ch == '"' {
				text = text[1 : len(text)-1]
			}
		}
	}
	return text
}

func Unquote(text string) string {
	if len(text) > 1 {
		if ch, _ := utf8.DecodeRuneInString(text); ch == '"' {
			if ch, _ = utf8.DecodeLastRuneInString(text); ch == '"' {
				text = text[1 : len(text)-1]
			}
		}
	}
	return text
}
