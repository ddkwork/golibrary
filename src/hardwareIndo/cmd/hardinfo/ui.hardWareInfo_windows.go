package hardinfo

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/ddkwork/golibrary/src/fynelib/mybutton"
	"github.com/ddkwork/golibrary/src/hardwareIndo"
)

type (
	Interface interface {
		CanvasObject(window fyne.Window) fyne.CanvasObject
	}
	object struct{}
)

func New() Interface { return &object{} }

func (o *object) CanvasObject(window fyne.Window) fyne.CanvasObject {
	macInfo := newGrid()
	macInfo.SetTitle("mac")
	macInfo.SetLabels("name", "guid", "row", "wmi", "align")

	cpu0 := newGrid()
	cpu0.SetTitle("cpu0")
	cpu0.SetLabels("eax", "ebx", "ecx", "edx")

	cpu1 := newGrid()
	cpu1.SetTitle("cpu1")
	cpu1.SetLabels("eax", "ebx", "ecx", "edx")

	ssd := newGrid()
	ssd.SetTitle("ssd")
	ssd.SetLabels("ModelNumber", "SerialNumber", "Version")

	read := mybutton.NewButtonWithIcon("read", nil, func() {
		h := hardwareIndo.New()
		if !h.MacInfo.Get() {
			return //todo show error info
		}
		macInfo.SetText(h.MacInfo.Description, "", h.MacInfo.Row, "", "")

		if !h.CpuInfo.Get() {
			return //todo show error info
		}
		cpu0.SetText(h.CpuInfo.FormatCpu0()...)
		cpu1.SetText(h.CpuInfo.FormatCpu1()...)

		if !h.SsdInfo.Get() {
			return //todo show error info
		}
		ssd.SetText(h.SsdInfo.ModelNumber, h.SsdInfo.SerialNumber, h.SsdInfo.Version)
	})
	hook := mybutton.NewButtonWithIcon("hook", nil, func() {
	})
	bypass := widget.NewCheck("ByPass baobao software", func(b bool) {

	})
	buttons := container.NewGridWithColumns(5, read, hook, bypass)
	box := container.NewGridWithColumns(2, container.NewVBox(
		ssd.WithRows(),
		cpu0.WithColumns(),
		cpu1.WithColumns(),
	), macInfo.WithRows())
	return container.NewVBox(box, buttons)
}
