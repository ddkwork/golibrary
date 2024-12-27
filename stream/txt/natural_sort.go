package txt

import (
	"slices"
)

func NaturalLess(s1, s2 string, caseInsensitive bool) bool {
	return NaturalCmp(s1, s2, caseInsensitive) < 0
}

func NaturalCmp(s1, s2 string, caseInsensitive bool) int {
	i1 := 0
	i2 := 0
	for i1 < len(s1) && i2 < len(s2) {
		c1 := s1[i1]
		c2 := s2[i2]
		d1 := c1 >= '0' && c1 <= '9'
		d2 := c2 >= '0' && c2 <= '9'
		switch {
		case d1 != d2:
			if d1 {
				return -1
			}
			return 1
		case !d1:

			if caseInsensitive {
				if c1 >= 'a' && c1 <= 'z' {
					c1 -= 'a' - 'A'
				}
				if c2 >= 'a' && c2 <= 'z' {
					c2 -= 'a' - 'A'
				}
			}
			if c1 != c2 {
				if c1 < c2 {
					return -1
				}
				return 1
			}
			i1++
			i2++
		default:

			for i1 < len(s1) && s1[i1] == '0' {
				i1++
			}
			for i1 < len(s1) && s1[i1] == '0' {
				i1++
			}
			for i2 < len(s2) && s2[i2] == '0' {
				i2++
			}

			nz1, nz2 := i1, i2
			for i1 < len(s1) && s1[i1] >= '0' && s1[i1] <= '9' {
				i1++
			}
			for i2 < len(s2) && s2[i2] >= '0' && s2[i2] <= '9' {
				i2++
			}

			if len1, len2 := i1-nz1, i2-nz2; len1 != len2 {
				if len1 < len2 {
					return -1
				}
				return 1
			}

			if nr1, nr2 := s1[nz1:i1], s2[nz2:i2]; nr1 != nr2 {
				if nr1 < nr2 {
					return -1
				}
				return 1
			}

			if nz1 != nz2 {
				if nz1 < nz2 {
					return -1
				}
				return 1
			}
		}

	}

	switch {
	case len(s1) == len(s2):
		if caseInsensitive {
			return NaturalCmp(s1, s2, false)
		}
		return 0
	case len(s1) < len(s2):
		return -1
	default:
		return 1
	}
}

func SortStringsNaturalAscending(in []string) {
	slices.SortFunc(in, func(a, b string) int { return NaturalCmp(a, b, true) })
}

func SortStringsNaturalDescending(in []string) {
	slices.SortFunc(in, func(a, b string) int { return NaturalCmp(b, a, true) })
}
