package txt

import (
	"strings"
)

func Wrap(prefix, text string, maxColumns int) string {
	var buffer strings.Builder
	for i, line := range strings.Split(text, "\n") {
		if i != 0 {
			buffer.WriteByte('\n')
		}
		buffer.WriteString(prefix)
		avail := maxColumns - len(prefix)
		for j, token := range strings.Fields(line) {
			if j != 0 {
				if 1+len(token) > avail {
					buffer.WriteByte('\n')
					buffer.WriteString(prefix)
					avail = maxColumns - len(prefix)
				} else {
					buffer.WriteByte(' ')
					avail--
				}
			}
			buffer.WriteString(token)
			avail -= len(token)
		}
	}
	return buffer.String()
}
