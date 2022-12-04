package doctable

import (
	"github.com/ddkwork/golibrary/skiaLib/widget/tabbar"
	"github.com/richardwilkes/unison"
)

type (
	Interface interface {
		unison.Dockable
		unison.TabCloser
	}
	object struct{ tabbar.Interface }
)

func New(title, tip string, background unison.Ink) Interface {
	return &object{Interface: tabbar.New(title, tip, background)}
}

func (o *object) MayAttemptClose() bool { return true }

func (o *object) AttemptClose() bool {
	if dc := unison.Ancestor[*unison.DockContainer](o); dc != nil {
		dc.Close(o)
		return true
	}
	return false
}
