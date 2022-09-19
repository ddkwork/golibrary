package main

import (
	"encoding/json"
	"flag"
	"fyne.io/fyne/v2/app"
	"github.com/ddkwork/golibrary/src/fynelib/terminal/cmd/fyneterm/data"
	"github.com/ddkwork/golibrary/src/fynelib/terminal/cmd/fyneterm/fyneterm"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"os"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "Show terminal debug messages")
	flag.Parse()

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustParseMessageFileBytes(fyneterm.ResourceActiveFrJson.Content(), fyneterm.ResourceActiveFrJson.Name())
	bundle.MustParseMessageFileBytes(fyneterm.ResourceActiveRuJson.Content(), fyneterm.ResourceActiveRuJson.Name())
	fyneterm.Localizer = i18n.NewLocalizer(bundle, os.Getenv("LANG"))

	a := app.New()
	a.SetIcon(data.Icon)
	th := fyneterm.NewTermTheme()
	a.Settings().SetTheme(th)
	w := fyneterm.NewTerminalWindow(a, th, debug)
	w.ShowAndRun()
}
