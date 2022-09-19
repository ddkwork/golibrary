package canvasobjectapi

import (
	"fyne.io/fyne/v2"
)

type Interface interface {
	CanvasObject(window fyne.Window) fyne.CanvasObject
}
