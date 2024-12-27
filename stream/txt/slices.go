package txt

import "strings"

func StringSliceToMap(slice []string) map[string]bool {
	m := make(map[string]bool, len(slice))
	for _, str := range slice {
		m[str] = true
	}
	return m
}

func MapToStringSlice(m map[string]bool) []string {
	s := make([]string, 0, len(m))
	for str := range m {
		s = append(s, str)
	}
	return s
}

func CloneStringSlice(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	out := make([]string, len(in))
	copy(out, in)
	return out
}

func RunesEqual(left, right []rune) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}

func CaselessSliceContains(slice []string, target string) bool {
	for _, one := range slice {
		if strings.EqualFold(one, target) {
			return true
		}
	}
	return false
}
