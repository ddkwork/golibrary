package widget

import (
	"cogentcore.org/core/gi"
	"cogentcore.org/core/goosi/driver/desktop"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func NewWindowRunAndWait(b *gi.Body, DropCallback func(names []string)) {
	w := b.NewWindow().Run()
	if w == nil || w.MainMgr == nil || w.MainMgr.RenderWin == nil {
		return
	}
	win := w.MainMgr.RenderWin.GoosiWin
	ww, ok := win.(*desktop.Window)
	if ok {
		ww.Glw.SetDropCallback(func(w *glfw.Window, names []string) {
			if DropCallback != nil {
				DropCallback(names)
			}
		})
	}
	w.Wait()
}
