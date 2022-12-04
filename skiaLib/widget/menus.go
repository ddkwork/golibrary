package widget

import "github.com/richardwilkes/unison"

func NewMenus(window *unison.Window, initializer func(unison.Menu)) {
	unison.DefaultMenuFactory().BarForWindow(window, initializer)
}
