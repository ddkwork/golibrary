package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/ddkwork/golibrary/src/fynelib/fyneTheme"
	"github.com/ddkwork/golibrary/src/fynelib/mybutton"
)

func main() {
	a := app.New()
	fyneTheme.Dark()
	w := a.NewWindow("Hello")
	hello := mybutton.NewButtonWithIcon("文本应该加粗吗", nil, func() {

	})
	w.SetContent(container.NewVBox(hello))
	w.Resize(fyne.NewSize(300, 300))
	w.CenterOnScreen()
	w.ShowAndRun()
}
