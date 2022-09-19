package assets

import (
	_ "embed"
	"fyne.io/fyne/v2"
)

//go:embed msapp/logo.png
var logo []byte

//go:embed   icons/system-file-manager.png
var systemFileManager []byte

//go:embed   icons/folder.png
var folder []byte

//go:embed icons/preferences-system-search.png
var search []byte

//go:embed error.jpeg
var error []byte

//go:embed log.png
var log []byte

//go:embed about.png
var about []byte

//go:embed terminal.png
var terminal []byte

//go:embed icons/system-software-install.png
var software []byte

type pngs struct {
	App               *fyne.StaticResource
	Error             *fyne.StaticResource
	Log               *fyne.StaticResource
	Terminal          *fyne.StaticResource
	About             *fyne.StaticResource
	SystemFileManager *fyne.StaticResource
	Folder            *fyne.StaticResource
	Search            *fyne.StaticResource
	Software          *fyne.StaticResource
}

var Pngs = pngs{
	App:               fyne.NewStaticResource("logo.png", logo),
	Error:             fyne.NewStaticResource("error.jpeg", error),
	Log:               fyne.NewStaticResource("log.jpeg", log),
	Terminal:          fyne.NewStaticResource("terminal.jpeg", terminal),
	About:             fyne.NewStaticResource("about.jpeg", about),
	SystemFileManager: fyne.NewStaticResource("system-file-manager.png", systemFileManager),
	Folder:            fyne.NewStaticResource("folder.png", folder),
	Search:            fyne.NewStaticResource("preferences-system-search.png", search),
	Software:          fyne.NewStaticResource("software.png", software),
}
