package widget

import (
	"cogentcore.org/core/gi"
	"cogentcore.org/core/ki"
	"cogentcore.org/core/states"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/units"
)

func NewButton(parent ki.Ki, name ...string) *gi.Button {
	button := gi.NewButton(parent, name...)
	button.OnWidgetAdded(func(w gi.Widget) {
		if w.PathFrom(button) == "parts/icon" {
			w.Style(func(s *styles.Style) {
				s.Min.Set(units.Dp(22)) // 工具栏大图标
			})
		}
		if lb, ok := w.(*gi.Label); ok {
			lb.Style(func(s *styles.Style) {
				if button.StateIs(states.Hovered) {
					s.Font.Size = units.Dp(17)
				} else {
					s.Font.Size = units.Dp(14)
				}
			})
		}
	})
	return button
}
