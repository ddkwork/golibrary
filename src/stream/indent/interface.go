package indent

import (
	"bytes"
	"strings"
)

type (
	Interface interface {
		Left(src string) string
		Right(src string) string
		Json(src string) string
		NewLine(src *bytes.Buffer) error //cyclicImport for not mycheck error
		MakeSpace(src string) string
	}
	object struct{}
)

func New() Interface { return &object{} }

func (o *object) MakeSpace(src string) string {
	return strings.Repeat(" ", len(src))
}
