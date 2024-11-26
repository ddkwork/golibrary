package txt

import "unicode"

type VowelChecker func(rune) bool

func IsVowel(ch rune) bool {
	if unicode.IsUpper(ch) {
		ch = unicode.ToLower(ch)
	}
	switch ch {
	case 'a', 'à', 'á', 'â', 'ä', 'æ', 'ã', 'å', 'ā',
		'e', 'è', 'é', 'ê', 'ë', 'ē', 'ė', 'ę',
		'i', 'î', 'ï', 'í', 'ī', 'į', 'ì',
		'o', 'ô', 'ö', 'ò', 'ó', 'œ', 'ø', 'ō', 'õ',
		'u', 'û', 'ü', 'ù', 'ú', 'ū':
		return true
	default:
		return false
	}
}

func IsVowely(ch rune) bool {
	if unicode.IsUpper(ch) {
		ch = unicode.ToLower(ch)
	}
	switch ch {
	case 'y', 'ÿ':
		return true
	default:
		return IsVowel(ch)
	}
}
