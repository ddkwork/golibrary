package widget

import (
	"cogentcore.org/core/gi"
)

func NewWindowRunAndWait(b *gi.Body, DropCallback func(names []string)) {
	b.RunMainWindow()
}
