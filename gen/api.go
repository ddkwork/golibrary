package gen

type Interface interface {
	P(v ...any)
	Enum(kindName string, kinds []string, tooltip []string)
	FileAction()
	ReadTemplates(path string)
}
