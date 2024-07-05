package gen

import (
	"unicode"
	"unicode/utf8"
)

type FileInfo struct {
	*File

	allEnums    []*enumInfo
	allMessages []*messageInfo

	allEnumsByPtr         map[*enumInfo]int
	allMessagesByPtr      map[*messageInfo]int
	allMessageFieldsByPtr map[*messageInfo]*structFields

	needRawDesc bool
}

type structFields struct {
	count      int
	unexported map[int]string
}

func (sf *structFields) append(name string) {
	if r, size := utf8.DecodeRuneInString(name); !unicode.IsUpper(r) {
		size = size
		if sf.unexported == nil {
			sf.unexported = make(map[int]string)
		}
		sf.unexported[sf.count] = name
	}
	sf.count++
}

func newFileInfo(file *File) *FileInfo {
	f := &FileInfo{File: file}

	var walkMessages func([]*Message, func(*Message))
	walkMessages = func(messages []*Message, f func(*Message)) {
		for _, m := range messages {
			f(m)
			walkMessages(m.Messages, f)
		}
	}
	initEnumInfos := func(enums []*Enum) {
		for _, enum := range enums {
			f.allEnums = append(f.allEnums, newEnumInfo(f, enum))
		}
	}
	initMessageInfos := func(messages []*Message) {
		for _, message := range messages {
			f.allMessages = append(f.allMessages, newMessageInfo(f, message))
		}
	}
	initExtensionInfos := func(extensions []*Extension) {}
	initEnumInfos(f.Enums)
	initMessageInfos(f.Messages)
	initExtensionInfos(f.Extensions)
	walkMessages(f.Messages, func(m *Message) {
		initEnumInfos(m.Enums)
		initMessageInfos(m.Messages)
		initExtensionInfos(m.Extensions)
	})

	if len(f.allEnums) > 0 {
		f.allEnumsByPtr = make(map[*enumInfo]int)
		for i, e := range f.allEnums {
			f.allEnumsByPtr[e] = i
		}
	}
	if len(f.allMessages) > 0 {
		f.allMessagesByPtr = make(map[*messageInfo]int)
		f.allMessageFieldsByPtr = make(map[*messageInfo]*structFields)
		for i, m := range f.allMessages {
			f.allMessagesByPtr[m] = i
			f.allMessageFieldsByPtr[m] = new(structFields)
		}
	}

	return f
}

type enumInfo struct {
	*Enum

	genJSONMethod    bool
	genRawDescMethod bool
}

func newEnumInfo(f *FileInfo, enum *Enum) *enumInfo {
	e := &enumInfo{Enum: enum}
	e.genJSONMethod = true
	e.genRawDescMethod = true
	return e
}

type messageInfo struct {
	*Message

	genRawDescMethod  bool
	genExtRangeMethod bool

	isTracked bool
	hasWeak   bool
}

func newMessageInfo(f *FileInfo, message *Message) *messageInfo {
	m := &messageInfo{Message: message}
	m.genRawDescMethod = true
	m.genExtRangeMethod = true
	return m
}
