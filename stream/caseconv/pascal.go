package caseconv

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

func ToPascal(str string) string {
	chunks := chunk(str)
	for i, c := range chunks {
		chunks[i] = cases.Title(language.English).String(c)
	}
	return strings.Join(chunks, "")
}
