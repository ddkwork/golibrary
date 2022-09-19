package grid

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/fpabl0/sparky-go/swid"
)

type (
	Interface interface {
		WithColumns() fyne.CanvasObject
		WithRows() fyne.CanvasObject
		SetTitle(title string)
		SetLabels(labels ...string)
		SetText(text ...string)
	}
	object struct {
		title         *widget.Label
		textFormField []*swid.TextFormField
	}
)

func New() Interface { return &object{} }

func newGrid() *object {
	return &object{
		title:         widget.NewLabel(""),
		textFormField: nil,
	}
}

func (o *object) WithColumns() fyne.CanvasObject {
	columns := container.NewGridWithColumns(len(o.textFormField) + 1)
	columns.Add(o.title)
	for _, field := range o.textFormField {
		columns.Add(field)
	}
	return columns
}
func (o *object) WithRows() fyne.CanvasObject {
	Rows := container.NewGridWithRows(len(o.textFormField) + 1)
	Rows.Add(o.title)
	for _, field := range o.textFormField {
		Rows.Add(field)
	}
	return Rows
}
func (o *object) SetTitle(title string) {
	o.title.Text = title
	o.title.Alignment = fyne.TextAlignCenter
}

func (o *object) SetLabels(labels ...string) {
	o.textFormField = make([]*swid.TextFormField, len(labels))
	for i, label := range labels {
		o.textFormField[i] = swid.NewTextFormField(label, "")
	}
}

func (o *object) SetText(text ...string) {
	for i, s := range text {
		o.textFormField[i].SetText(s)
	}
}
