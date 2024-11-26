package txt

import (
	"fmt"
	"strings"

	"golang.org/x/exp/constraints"
)

func Comma[T constraints.Integer | constraints.Float](value T) string {
	return CommaFromStringNum(fmt.Sprintf("%v", value))
}

func CommaFromStringNum(s string) string {
	var buffer strings.Builder
	if strings.HasPrefix(s, "-") {
		buffer.WriteByte('-')
		s = s[1:]
	}
	parts := strings.Split(s, ".")
	i := 0
	needComma := false
	if len(parts[0])%3 != 0 {
		i += len(parts[0]) % 3
		buffer.WriteString(parts[0][:i])
		needComma = true
	}
	for ; i < len(parts[0]); i += 3 {
		if needComma {
			buffer.WriteByte(',')
		} else {
			needComma = true
		}
		buffer.WriteString(parts[0][i : i+3])
	}
	if len(parts) > 1 {
		buffer.Write([]byte{'.'})
		buffer.WriteString(parts[1])
	}
	return buffer.String()
}
