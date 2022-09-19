package driverTool

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/ddkwork/golibrary/src/driverTool/driver"
	"github.com/ddkwork/golibrary/src/fynelib/canvasobjectapi"

	"github.com/fpabl0/sparky-go/swid"
	"io/fs"
	"path/filepath"
)

type (
	Interface interface {
		canvasobjectapi.Interface
		Driver() *driver.Object
		SetUnloadVmmTapped(unloadVmmTapped func())
		SetLoadVmmTapped(loadVmmTapped func())
	}
	object struct {
		drivers         []string
		driver          *driver.Object
		loadVmmTapped   func()
		unloadVmmTapped func()
	}
)

func (o *object) SetUnloadVmmTapped(unloadVmmTapped func()) { o.unloadVmmTapped = unloadVmmTapped }
func (o *object) SetLoadVmmTapped(loadVmmTapped func())     { o.loadVmmTapped = loadVmmTapped }
func (o *object) Driver() *driver.Object                    { return o.driver }
func New() Interface {
	return &object{
		drivers: make([]string, 0),
		driver:  driver.NewObject(),
	}
}

func (o *object) CanvasObject(window fyne.Window) fyne.CanvasObject {
	o.WalkAllDriverPath("")
	path := swid.NewSelectFormField("path", "", o.drivers)
	link := swid.NewTextFormField("link", "")
	path.OnChanged = func(s string) {
		if o.driver.DeviceName == "" {
			ext := filepath.Ext(s)
			base := filepath.Base(s)
			base = base[:len(base)-len(ext)]
			link.SetText(base)
			o.driver.DeviceName = base
		} else {
			link.SetText(o.driver.DeviceName)
		}
	}
	ioCode := swid.NewTextFormField("ioCode", "")
	load := widget.NewButton("load", func() {
		if !o.driver.Load(path.Selected()) {
			return
		}
	})
	unload := widget.NewButton("unload", func() {
		if !o.driver.Unload() {
			return
		}
	})

	loadVmm := widget.NewButton("loadVmm", func() { //todo check vm status for pass bsod
		if o.loadVmmTapped == nil {
			return
		}
		o.loadVmmTapped()
		if o.driver.Status == 0 {
			return
		}
	})
	unloadVmm := widget.NewButton("unloadVmm", func() {
		if o.unloadVmmTapped == nil {
			return
		}
		o.unloadVmmTapped()
		if o.driver.Status == 0 {
			return
		}
	})

	errCode := swid.NewTextFormField("errCode", "")
	ntstatus := swid.NewTextFormField("ntstatus", "")
	hresult := swid.NewTextFormField("hresult", "")
	winerror := swid.NewTextFormField("winerror", "")

	reload := swid.NewTextFormField("reload path", "")
	reload.OnChanged = func(s string) {
		o.drivers = o.drivers[:0]
		path.Options = path.Options[:0]
		o.WalkAllDriverPath(s)
		path.Options = o.drivers
	}
	form := container.NewGridWithColumns(1,
		reload,
		path,
		link,
		ioCode,
		errCode,
		ntstatus,
		hresult,
		winerror,
		container.NewGridWithColumns(2, load, unload),
		container.NewGridWithColumns(2, loadVmm, unloadVmm),
	)
	split := container.NewHSplit(form, mycheck.CanvasObject(window))
	split.Offset = 0.4
	return split
}

func (o *object) WalkAllDriverPath(root string) bool {
	newRoot := root
	if root == "" {
		newRoot = "."
	}
	i := 0
	return mylog.Error(filepath.Walk(newRoot, func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ".sys" {
			i++
			o.drivers = append(o.drivers, path)
		}
		return nil
	}))
}
