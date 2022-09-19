package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/ddkwork/golibrary/src/fynelib/notes"
)

func main() {
	a := app.NewWithID("xyz.andy.notes")
	a.Settings().SetTheme(&notes.MyTheme{})
	w := a.NewWindow("Notes")

	list := &notes.Notelist{Pref: a.Preferences()}
	list.Load()
	notesUI := &notes.Ui{Notes: list}

	w.SetContent(notesUI.LoadUI())
	notesUI.RegisterKeys(w)

	w.Resize(fyne.NewSize(400, 320))
	w.ShowAndRun()
}
