package widget

import (
	"github.com/ddkwork/hyperdbgui/ux/widget/doctable"
	"github.com/ddkwork/hyperdbgui/ux/widget/tabbar"
	"github.com/richardwilkes/unison"
)

type (
	Interface interface {
		tabbar.Interface
		doctable.Interface
		NewMenus(window *unison.Window, initializer func(unison.Menu))
		NewPopMenus(window *unison.Window, initializer func(unison.Menu))
		OpenWith()
		DropFiles()
	}
	object struct{}
)
