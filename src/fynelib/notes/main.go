package notes

import (
	"fmt"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Ui struct {
	current *note
	Notes   *Notelist

	content *widget.Entry
	list    *widget.List
}

func (u *Ui) addNote() {
	newNote := u.Notes.add()
	u.setNote(newNote)
}

func (u *Ui) setNote(n *note) {
	u.content.Unbind()
	if n == nil {
		u.content.SetText(u.placeholderContent())
		return
	}
	u.current = n
	u.content.Bind(n.content)
	u.list.Refresh()
}

func (u *Ui) buildList() *widget.List {
	l := widget.NewList(
		func() int {
			return len(u.Notes.notes)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Title")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			l := obj.(*widget.Label)
			n := u.Notes.notes[id]
			l.Bind(n.title())
		})

	l.OnSelected = func(id widget.ListItemID) {
		n := u.Notes.notes[id]
		u.setNote(n)
	}

	return l
}

func (u *Ui) removeCurrentNote() {
	u.Notes.remove(u.current)
	if len(u.Notes.notes) > 0 {
		u.setNote(u.Notes.notes[0])
	} else {
		u.setNote(nil)
	}
	u.list.Refresh()
}

func (u *Ui) LoadUI() fyne.CanvasObject {
	u.content = widget.NewMultiLineEntry()
	u.content.SetText(u.placeholderContent())

	u.list = u.buildList()

	if len(u.Notes.notes) > 0 {
		u.setNote(u.Notes.notes[0])
		u.list.Select(0)
	}

	bar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			u.addNote()
		}),
		widget.NewToolbarAction(theme.ContentRemoveIcon(), func() {
			u.removeCurrentNote()
		}),
	)

	side := fyne.NewContainerWithLayout(layout.NewBorderLayout(bar, nil, nil, nil),
		bar, container.NewVScroll(u.list))

	return newAdaptiveSplit(side, u.content)
}

func (u *Ui) RegisterKeys(w fyne.Window) {
	shortcut := &desktop.CustomShortcut{KeyName: fyne.KeyN, Modifier: desktop.ControlModifier}
	if runtime.GOOS == "darwin" {
		shortcut.Modifier = desktop.SuperModifier
	}

	w.Canvas().AddShortcut(shortcut, func(_ fyne.Shortcut) {
		u.addNote()
	})
}

func (u *Ui) placeholderContent() string {
	text := "Welcome!\nTap '+' in the toolbar to add a note."
	if fyne.CurrentDevice().HasKeyboard() {
		modifier := "ctrl"
		if runtime.GOOS == "darwin" {
			modifier = "cmd"
		}
		text += fmt.Sprintf("\n\nOr use they keyboard shortcut %s+N.", modifier)
	}
	return text
}

func main() {
	a := app.NewWithID("xyz.andy.notes")
	a.Settings().SetTheme(&MyTheme{})
	w := a.NewWindow("Notes")

	list := &Notelist{Pref: a.Preferences()}
	list.Load()
	notesUI := &Ui{Notes: list}

	w.SetContent(notesUI.LoadUI())
	notesUI.RegisterKeys(w)

	w.Resize(fyne.NewSize(400, 320))
	w.ShowAndRun()
}
