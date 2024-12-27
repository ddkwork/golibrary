package txt

import "strings"

func IsTruthy(in string) bool {
	in = strings.ToLower(in)
	return in == "1" || in == "true" || in == "yes" || in == "on"
}
