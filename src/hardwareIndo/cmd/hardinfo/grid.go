package hardinfo

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/fpabl0/sparky-go/swid"
)

type (
	grid struct {
		title         *widget.Label
		textFormField []*swid.TextFormField
	}
)

func newGrid() *grid {
	return &grid{
		title:         widget.NewLabel(""),
		textFormField: nil,
	}
}

func (g *grid) WithColumns() fyne.CanvasObject {
	columns := container.NewGridWithColumns(len(g.textFormField) + 1)
	columns.Add(g.title)
	for _, field := range g.textFormField {
		columns.Add(field)
	}
	return columns
}
func (g *grid) WithRows() fyne.CanvasObject {
	Rows := container.NewGridWithRows(len(g.textFormField) + 1)
	Rows.Add(g.title)
	for _, field := range g.textFormField {
		Rows.Add(field)
	}
	return Rows
}
func (g *grid) SetTitle(title string) {
	g.title.Text = title
	g.title.Alignment = fyne.TextAlignCenter
}

func (g *grid) SetLabels(labels ...string) {
	g.textFormField = make([]*swid.TextFormField, len(labels))
	for i, label := range labels {
		g.textFormField[i] = swid.NewTextFormField(label, "")
	}
}

func (g *grid) SetText(text ...string) {
	for i, s := range text {
		g.textFormField[i].SetText(s)
	}
}
