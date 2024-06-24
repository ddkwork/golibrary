package txt

func FirstN(s string, n int) string {
	if n < 1 {
		return ""
	}
	r := []rune(s)
	if n > len(r) {
		return s
	}
	return string(r[:n])
}

func LastN(s string, n int) string {
	if n < 1 {
		return ""
	}
	r := []rune(s)
	if n > len(r) {
		return s
	}
	return string(r[len(r)-n:])
}

func Truncate(s string, count int, keepFirst bool) string {
	var result string
	if keepFirst {
		result = FirstN(s, count)
	} else {
		result = LastN(s, count)
	}
	if result != s {
		if keepFirst {
			result += "…"
		} else {
			result = "…" + result
		}
	}
	return result
}
