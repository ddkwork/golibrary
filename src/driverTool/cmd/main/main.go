package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/ddkwork/golibrary/src/driverTool"
	"github.com/ddkwork/golibrary/src/fynelib/fyneTheme"
)

func main() {
	a := app.NewWithID("com.rows.app")
	a.SetIcon(nil)
	fyneTheme.Dark()
	w := a.NewWindow("windows nt driver tool")
	w.Resize(fyne.NewSize(1040, 480))
	w.SetMaster()
	w.CenterOnScreen()
	w.SetContent(driverTool.New().CanvasObject(w))
	w.ShowAndRun()
}
