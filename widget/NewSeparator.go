package widget

import (
	"cogentcore.org/core/gi"
	"cogentcore.org/core/ki"
	"cogentcore.org/core/styles"
)

func NewSeparatorWithLabel(title string, parent ki.Ki) {
	gi.NewLabel(parent).SetText(title).Style(func(s *styles.Style) { s.Align.Self = styles.Center })
	gi.NewSeparator(parent)
}
