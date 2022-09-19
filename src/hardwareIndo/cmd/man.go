package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/ddkwork/golibrary/src/fynelib/fyneTheme"
	"github.com/ddkwork/golibrary/src/hardwareIndo"
	"github.com/ddkwork/golibrary/src/hardwareIndo/cmd/hardinfo"
)

func main() {
	a := app.NewWithID("com.rows.app")
	//a.SetIcon(nil)
	fyneTheme.Dark()
	w := a.NewWindow("hardInfo")
	//w.Resize(fyne.NewSize(140, 580))
	//w.SetMaster()
	w.CenterOnScreen()
	h := hardinfo.New()
	w.SetContent(h.CanvasObject(w))
	w.ShowAndRun()
}

func test() {
	h := hardwareIndo.New()
	if !h.SsdInfo.Get() { //todo bug cpu pkg init
		return
	}
	if !h.CpuInfo.Get() {
		return
	}
	if !h.MacInfo.Get() {
		return
	}
}
