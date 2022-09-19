package mybutton

import (
	"encoding/hex"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/src/fynelib/fyneTheme"

	"image/color"
)

const (
	// CloseCursor is the mouse cursor that indicates a close action
	CloseCursor desktop.StandardCursor = iota + desktop.VResizeCursor // add to the end of the fyne list
)

type closeButton struct {
	widget.Button
	bg *canvas.Rectangle
}

func (c *closeButton) Cursor() desktop.Cursor { return CloseCursor }
func (c *closeButton) MouseIn(*desktop.MouseEvent) {
	fyneTheme.MouseIn()
	c.Refresh()
	//https://github.com/microsoft/PowerToys/releases
	b, err := hex.DecodeString("dedefa00")
	if !mylog.Error(err) {
		return
	}
	c.bg.FillColor = color.NRGBA{
		R: b[0],
		G: b[1],
		B: b[2],
		A: b[3],
	}
	c.bg.Refresh()
}
func (c *closeButton) MouseMoved(*desktop.MouseEvent) {}
func (c *closeButton) MouseOut() {
	fyneTheme.MouseOut()
	c.Refresh()
	c.bg.FillColor = color.Transparent
	c.bg.Refresh()
}
func NewButtonWithIcon(label string, icon fyne.Resource, tapped func()) fyne.CanvasObject {
	b := &closeButton{
		Button: widget.Button{
			DisableableWidget: widget.DisableableWidget{},
			Text:              label,
			Icon:              icon,
			Importance:        widget.LowImportance,
			Alignment:         0,
			IconPlacement:     0,
			OnTapped:          tapped,
		},
		bg: canvas.NewRectangle(color.Transparent),
	}
	b.ExtendBaseWidget(b)
	//return b
	return container.NewMax(b.bg, b)
}
