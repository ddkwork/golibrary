package txt

import (
	_ "embed"
	"fmt"
	"regexp"
	"strings"

	"github.com/ddkwork/golibrary/mylog"
)

var stdAllCaps string

var StdAllCaps = MustNewAllCaps(strings.Split(NormalizeLineEndings(stdAllCaps), "\n")...)

type AllCaps struct {
	regex *regexp.Regexp
}

func NewAllCaps(in ...string) (*AllCaps, error) {
	var buffer strings.Builder
	for _, str := range in {
		if buffer.Len() > 0 {
			buffer.WriteByte('|')
		}
		buffer.WriteString(FirstToUpper(strings.ToLower(str)))
	}
	r := mylog.Check2(regexp.Compile(fmt.Sprintf("(%s)(?:$|[A-Z])", buffer.String())))

	return &AllCaps{regex: r}, nil
}

func MustNewAllCaps(in ...string) *AllCaps {
	return mylog.Check2(NewAllCaps(in...))
}
