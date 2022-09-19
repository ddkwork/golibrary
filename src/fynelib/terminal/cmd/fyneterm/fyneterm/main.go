//go:generate fyne bundle -o translation.go ../../translation/

package fyneterm

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"github.com/ddkwork/golibrary/src/fynelib/terminal"
	"github.com/ddkwork/golibrary/src/fynelib/terminal/cmd/fyneterm/data"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"image/color"
)

const termOverlay = fyne.ThemeColorName("termOver")

var Localizer *i18n.Localizer

func setupListener(t *terminal.Terminal, w fyne.Window) {
	listen := make(chan terminal.Config)
	go func() {
		for {
			config := <-listen

			if config.Title == "" {
				w.SetTitle(termTitle())
			} else {
				w.SetTitle(termTitle() + ": " + config.Title)
			}
		}
	}()
	t.AddListener(listen)
}

func termTitle() string {
	return Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Title",
			Other: "Fyne Terminal",
		},
	})
}

func guessCellSize() fyne.Size {
	cell := canvas.NewText("M", color.White)
	cell.TextStyle.Monospace = true

	return cell.MinSize()
}

func NewTerminalWindow(a fyne.App, th fyne.Theme, debug bool) fyne.Window {
	w := a.NewWindow(termTitle())
	w.SetPadded(false)

	bg := canvas.NewRectangle(theme.BackgroundColor())
	img := canvas.NewImageFromResource(data.FyneLogo)
	img.FillMode = canvas.ImageFillContain
	over := canvas.NewRectangle(th.Color(termOverlay, a.Settings().ThemeVariant()))

	ch := make(chan fyne.Settings)
	go func() {
		for {
			<-ch

			bg.FillColor = theme.BackgroundColor()
			bg.Refresh()
			over.FillColor = th.Color(termOverlay, a.Settings().ThemeVariant())
			over.Refresh()
		}
	}()
	a.Settings().AddChangeListener(ch)

	t := terminal.New()
	t.SetDebug(debug)
	setupListener(t, w)
	w.SetContent(container.NewMax(bg, img, over, t))

	cellSize := guessCellSize()
	w.Resize(fyne.NewSize(cellSize.Width*80, cellSize.Height*24))
	w.Canvas().Focus(t)

	t.AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyN, Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift},
		func(_ fyne.Shortcut) {
			w := NewTerminalWindow(a, th, debug)
			w.Show()
		})
	go func() {
		err := t.RunLocalShell()
		if err != nil {
			fyne.LogError("Failure in terminal", err)
		}
		w.Close()
	}()

	return w
}
