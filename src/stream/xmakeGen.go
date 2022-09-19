package stream

import (
	"strconv"
	"strings"
)

func (s *Stream) WriteXMakeBody(key string, values ...string) {
	isNewLineKey := len(values) > 1
	if strings.HasPrefix(values[0], "wdk") {
		isNewLineKey = false
	}
	s.WriteString(key)
	s.WriteString("(")
	if isNewLineKey {
		s.NewLine()
	}
	for i, value := range values {
		if isNewLineKey {
			s.WriteString("\t")
		}
		s.WriteString(strconv.Quote(value))
		if key == "add_includedirs" {
			s.WriteString(",{public=true}") //for deps add
		}
		if i+1 < len(values) {
			s.WriteString(",")
			if isNewLineKey {
				s.NewLine()
			}
		}
	}
	if isNewLineKey {
		s.NewLine()
	}
	s.WriteString(")")
	s.NewLine()
}
