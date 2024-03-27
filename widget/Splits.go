package widget

import (
	"cogentcore.org/core/colors"
	"cogentcore.org/core/gi"
	"cogentcore.org/core/ki"
	"cogentcore.org/core/styles"
)

func NewVSplits(parent ki.Ki) *gi.Splits { // Horizontal and vertical
	splits := gi.NewSplits(parent)
	splits.Style(func(s *styles.Style) {
		s.Direction = styles.Column
		s.Background = colors.C(colors.Scheme.SurfaceContainerLow)
	})
	return splits
}
