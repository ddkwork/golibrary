package fyneTheme

import (
	_ "embed"
	"fyne.io/fyne/v2"
	"sync"
)

//go:embed ttf/HarmonyOS_Sans_SC_Light.ttf
var ttfBuf []byte

var HarmonyOsSansScBoldTtf = &fyne.StaticResource{
	StaticName:    "HarmonyOS_Sans_SC_Bold.ttf",
	StaticContent: ttfBuf,
}

type (
	myTheme struct {
		current any
		lock    sync.Mutex
		dark    *dark
		light   *light
	}
)

func newTheme() *myTheme {
	return &myTheme{
		current: nil,
		lock:    sync.Mutex{},
		dark:    &dark{SizeNameText: 14},
		light:   &light{SizeNameText: 14},
	}
}

var defaults = newTheme()

func Dark() *myTheme {
	t := defaults
	t.current = t.dark
	fyne.CurrentApp().Settings().SetTheme(t.dark)
	return t
}
func Light() *myTheme {
	t := defaults
	t.current = t.light
	fyne.CurrentApp().Settings().SetTheme(t.light)
	return t
}

func (t *myTheme) kind() any {
	_, ok := t.current.(*dark)
	if ok {
		return t.current.(*dark)
	}
	return t.current.(*light)
}

//var lock = new(sync.Mutex)

func MouseIn() {
	switch defaults.kind().(type) {
	case *dark:
		defaults.lock.Lock()
		defaults.dark.SizeNameText = 15
		defaults.dark.Bold = true
		fyne.CurrentApp().Settings().SetTheme(defaults.dark)
		defaults.lock.Unlock()
	case *light:
		defaults.lock.Lock()
		defaults.light.SizeNameText = 15
		defaults.light.Bold = true
		fyne.CurrentApp().Settings().SetTheme(defaults.light)
		defaults.lock.Unlock()
	}
}
func MouseOut() {
	switch defaults.kind().(type) {
	case *dark:
		defaults.lock.Lock()
		defaults.dark.SizeNameText = 14
		defaults.dark.Bold = false
		fyne.CurrentApp().Settings().SetTheme(defaults.dark)
		defaults.lock.Unlock()
	case *light:
		defaults.lock.Lock()
		defaults.light.SizeNameText = 14
		defaults.light.Bold = false
		fyne.CurrentApp().Settings().SetTheme(defaults.light)
		defaults.lock.Unlock()
	}
}
