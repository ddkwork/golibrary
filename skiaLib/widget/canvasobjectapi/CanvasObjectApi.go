package canvasobjectapi

import (
	"github.com/richardwilkes/unison"
)

type Interface interface {
	CanvasObject(window *unison.Window) *unison.Panel
}
