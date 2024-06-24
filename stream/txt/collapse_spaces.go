package txt

import "strings"

func CollapseSpaces(in string) string {
	var buffer strings.Builder
	lastWasSpace := false
	for i, r := range in {
		if r == ' ' {
			if !lastWasSpace {
				if i != 0 {
					buffer.WriteByte(' ')
				}
				lastWasSpace = true
			}
		} else {
			buffer.WriteRune(r)
			lastWasSpace = false
		}
	}
	str := buffer.String()
	if lastWasSpace && len(str) > 0 {
		str = str[:len(str)-1]
	}
	return str
}
