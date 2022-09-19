package cpp2go

import (
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/src/clang"

	"github.com/ddkwork/golibrary/src/stream"
	"github.com/ddkwork/golibrary/src/stream/tool"
	goScanner "go/scanner"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// ╔════╤══════════════════════════════════════════════╤══════════════════════════════════════════════════════════════╤══════════════════════════════════════════════════════════╤═════════╤══════╗
// ║ ID ║                     api                      ║                           function                           ║                           note                           ║ chinese ║ todo ║
// ╠════╪══════════════════════════════════════════════╪══════════════════════════════════════════════════════════════╪══════════════════════════════════════════════════════════╪═════════╪══════╣
// ║ 1  ║ Translate(root string) (ok bool)             ║ call Scan(root string) (ok bool)                             ║                                                          ║         ║      ║
// ╠════╪══════════════════════════════════════════════╪══════════════════════════════════════════════════════════════╪══════════════════════════════════════════════════════════╪═════════╪══════╣
// ║ 2  ║ RemoveComment(root string) (ok bool)         ║ Scan(root string) (ok bool)                                  ║                                                          ║         ║      ║
// ╠════╪══════════════════════════════════════════════╪══════════════════════════════════════════════════════════════╪══════════════════════════════════════════════════════════╪═════════╪══════╣
// ║ 3  ║ scan() (ok bool)                             ║ walk root dir and filter file ext for scan lexer             ║                                                          ║         ║      ║
// ╠════╪══════════════════════════════════════════════╪══════════════════════════════════════════════════════════════╪══════════════════════════════════════════════════════════╪═════════╪══════╣
// ║ 4  ║ generateNoCommentFile(body string) (ok bool) ║ generate No Comment File for checking                        ║                                                          ║         ║      ║
// ╠════╪══════════════════════════════════════════════╪══════════════════════════════════════════════════════════════╪══════════════════════════════════════════════════════════╪═════════╪══════╣
// ║ 5  ║ translate() (ok bool)                        ║ FindAllBlock,makeBlock,translateBlock,bindBlockType,handleP- ║                                                          ║         ║      ║
// ║    ║                                              ║ kgOrApiName,generateGoCodes,reset scanner                    ║                                                          ║         ║      ║
// ╠════╪══════════════════════════════════════════════╪══════════════════════════════════════════════════════════════╪══════════════════════════════════════════════════════════╪═════════╪══════╣
// ║ 6  ║ reset()                                      ║ reset scanner ctx when finished every file convert work      ║                                                          ║         ║      ║
// ╠════╪══════════════════════════════════════════════╪══════════════════════════════════════════════════════════════╪══════════════════════════════════════════════════════════╪═════════╪══════╣
// ║ 7  ║ FindAllBlock() (ok bool)                     ║ MergeLineElem,FindTypedefs,FindEnums,FindDefines,FindExtern- ║ first merge every line text elem lexer to line text, and ║         ║      ║
// ║    ║                                              ║ s,FindMethods                                                ║ reset allLines when cpp type was found evey time         ║         ║      ║
// ╠════╪══════════════════════════════════════════════╪══════════════════════════════════════════════════════════════╪══════════════════════════════════════════════════════════╪═════════╪══════╣
// ║ 8  ║ MergeLineElem() (ok bool)                    ║ Merge Line Elem (every lexer word) in to line text           ║                                                          ║         ║      ║
// ╠════╪══════════════════════════════════════════════╪══════════════════════════════════════════════════════════════╪══════════════════════════════════════════════════════════╪═════════╪══════╣
// ║ 9  ║ FindTypedefs()                               ║ Find Typedefs                                                ║ it will be include struct or point reType,enum,          ║         ║      ║
// ╠════╪══════════════════════════════════════════════╪══════════════════════════════════════════════════════════════╪══════════════════════════════════════════════════════════╪═════════╪══════╣
// ║ 10 ║ FindEnums()                                  ║ Find Enums                                                   ║                                                          ║         ║      ║
// ╠════╪══════════════════════════════════════════════╪══════════════════════════════════════════════════════════════╪══════════════════════════════════════════════════════════╪═════════╪══════╣
// ║ 11 ║ FindDefines()                                ║ Find Defines                                                 ║                                                          ║         ║      ║
// ╠════╪══════════════════════════════════════════════╪══════════════════════════════════════════════════════════════╪══════════════════════════════════════════════════════════╪═════════╪══════╣
// ║ 12 ║ FindExterns()                                ║ Find Externs                                                 ║                                                          ║         ║      ║
// ╠════╪══════════════════════════════════════════════╪══════════════════════════════════════════════════════════════╪══════════════════════════════════════════════════════════╪═════════╪══════╣
// ║ 13 ║ FindMethods()                                ║ Find Methods                                                 ║                                                          ║         ║      ║
// ╠════╪══════════════════════════════════════════════╪══════════════════════════════════════════════════════════════╪══════════════════════════════════════════════════════════╪═════════╪══════╣
// ║ 14 ║ ReSetAllLines(block []LineInfo)              ║ when founded struct,method etc reset lines to null           ║                                                          ║         ║      ║
// ╠════╪══════════════════════════════════════════════╪══════════════════════════════════════════════════════════════╪══════════════════════════════════════════════════════════╪═════════╪══════╣
// ║ 15 ║ makeBlock() (ok bool)                        ║ all line merge in to block,for example all struct,all define ║                                                          ║         ║      ║
// ╠════╪══════════════════════════════════════════════╪══════════════════════════════════════════════════════════════╪══════════════════════════════════════════════════════════╪═════════╪══════╣
// ║ 16 ║ translateBlock() (ok bool)                   ║ range all bock and convert to go for get go object name,type ║                                                          ║         ║      ║
// ║    ║                                              ║ etc                                                          ║                                                          ║         ║      ║
// ╠════╪══════════════════════════════════════════════╪══════════════════════════════════════════════════════════════╪══════════════════════════════════════════════════════════╪═════════╪══════╣
// ║ 17 ║ bindBlockType() (ok bool)                    ║ bind all cpp type to go type                                 ║                                                          ║         ║      ║
// ╠════╪══════════════════════════════════════════════╪══════════════════════════════════════════════════════════════╪══════════════════════════════════════════════════════════╪═════════╪══════╣
// ║ 18 ║ handlePkgOrApiName() (ok bool)               ║ check go syntax                                              ║                                                          ║         ║      ║
// ╠════╪══════════════════════════════════════════════╪══════════════════════════════════════════════════════════════╪══════════════════════════════════════════════════════════╪═════════╪══════╣
// ║ 19 ║ generateGoCodes() (ok bool)                  ║ last work                                                    ║                                                          ║         ║      ║
// ╚════╧══════════════════════════════════════════════╧══════════════════════════════════════════════════════════════╧══════════════════════════════════════════════════════════╧═════════╧══════╝
type (
	word struct {
		Line  int
		Text  string
		token token.Token
	}
	LineInfo struct {
		word   //token is last word
		tokens []token.Token
		elems  []string
	}
	fileInfo struct {
		LineInfo
		path     string
		words    []word
		allLines []LineInfo //set api is MergeLineElem
		typedefs []LineInfo
		enums    []LineInfo
		defines  []LineInfo
		externs  []LineInfo
		methods  []LineInfo
	}
)

func (f *fileInfo) FindTypedefs() {
	//todo typedef RFLAGS *PRFLAGS;
	//isDefineStart := false
	//isDefineEnd := false
	isDef := false
	name := ""
	for i := 0; i < len(f.allLines); i++ {
		lineInfo := f.allLines[i]
		split := strings.Split(lineInfo.Text, " ")
		for i, s2 := range split {
			switch {
			case s2 == typedefKInd && split[i+1] == structKInd:
				name = split[2]
				isDef = true
				break
			case s2 == "}" && strings.Contains(name, split[1]):
				isDef = false
				f.typedefs = append(f.typedefs, lineInfo)
				break
			}
		}
		if isDef {
			f.typedefs = append(f.typedefs, lineInfo)
		}
	}
	f.ReSetAllLines(f.typedefs)
}
func (f *fileInfo) FindEnums() {
	//f.enums = enums
	//f.ReSetAllLines(f.enums)
}
func (f *fileInfo) FindDefines() {
	//#define CONFIG_FILE_NAME L"config.ini"
	//#define VMM_DRIVER_NAME "hprdbghv"
	for i := 0; i < len(f.allLines); i++ {
		lineInfo := f.allLines[i]
		if strings.Contains(lineInfo.Text, defineKInd) {
			block := f.allLines[i:]
			for _, s := range block {
				f.defines = append(f.defines, lineInfo)
				//mylog.Trace(fmt.Sprint(s.Line), s)
				if !strings.Contains(s.Text, `\`) {
					break
				}
			}
		}
	}
	f.ReSetAllLines(f.defines)
}
func (f *fileInfo) FindExterns() {
	for i := 0; i < len(f.allLines); i++ {
		lineInfo := f.allLines[i]
		lineTextElems := strings.Split(lineInfo.Text, " ")
		for _, elem := range lineTextElems {
			switch {
			case elem == externKInd:
				//mylog.Trace(fmt.Sprint(lineInfo.Line), lineInfo.Text)
				f.externs = append(f.externs, lineInfo)
			}
		}
	}
	f.ReSetAllLines(f.externs)
}
func (f *fileInfo) FindMethods() {
	isArgEnd := false
	isApiStart := false
	//apiStart := 0
	name := ""
	for i := 0; i < len(f.allLines); i++ {
		lineInfo := f.allLines[i]
		lineTextElems := strings.Split(lineInfo.Text, " ")
		for j, textElem := range lineTextElems {
			if name == "" {
				switch {
				case lineTextElems[0] == `#`:
					continue
				case lineTextElems[0] == `typedef`:
					continue
				}
				if !f.isApi(textElem, methodArgStartKInd) {
					continue
				}
				name = lineTextElems[j-1]
				//apiStart = j - 1
				//apiStart = apiStart
			}
			if name == "" && len(lineTextElems) < 3 {
				continue
			}
			if !isArgEnd {
				if !f.isApi(textElem, methodArgEndKInd) {
					continue
				}
				isArgEnd = true
			}
			if !isApiStart {
				if !f.isApi(textElem, methodStartKInd) {
					continue
				}
				isApiStart = true
			}
			f.methods = append(f.methods, lineInfo)
			if lineInfo.token == token.EOF {
				//mylog.Warning("name", strconv.Quote(name)+" //"+fmt.Sprint(block[0].Line))
				//mylog.Warning("start", block[0].Text+" //"+fmt.Sprint(block[0].Line))
				//mylog.Warning("end", strconv.Quote(lineInfo.Text)+" //"+fmt.Sprint(lineInfo.Line)+"\n")
				//for _, info := range apiBlock {
				//	mylog.Info(fmt.Sprint(info.Line), info.Text)
				//}
				f.methods = append(f.methods, lineInfo)
				return
			}
		}
	}
	f.ReSetAllLines(f.methods)
}
func (f *fileInfo) isApi(textElem, substr string) (ok bool) {
	switch {
	case strings.Contains(textElem, substr):
		return true
	case strings.Contains(textElem, `#`):
		return
	case strings.Contains(textElem, `#if`):
		return
	default:
		return
	}
}
func (f *fileInfo) ReSetAllLines(block []LineInfo) {
	for i, info := range block {
		if info.Line == f.allLines[i].Line {
			f.allLines[i] = LineInfo{
				word: word{
					Line:  0,
					Text:  "",
					token: 0,
				},
				tokens: f.tokens[:0],
				elems:  f.elems[:0],
			}
		}
	}
}

type (
	Ir interface {
		translateBlock() (ok bool)
		bindBlockType() (ok bool)
		handlePkgOrApiName() (ok bool)
		generateGoCodes() (ok bool) //with format
	}
	Scanner interface {
		Translate(root string) (ok bool)
		RemoveComment(root string) (ok bool)
	}
	scanner struct {
		fileInfo
		isTranslate bool
		stream      *stream.Stream
		paths       []string
	}
)

func newScanner() Scanner {
	return &scanner{
		fileInfo: fileInfo{
			LineInfo: LineInfo{
				word: word{
					Line:  0,
					Text:  "",
					token: 0,
				},
				tokens: make([]token.Token, 0),
				elems:  make([]string, 0),
			},
			path:     "",
			words:    make([]word, 0),
			allLines: make([]LineInfo, 0),
			typedefs: make([]LineInfo, 0),
			enums:    make([]LineInfo, 0),
			defines:  make([]LineInfo, 0),
			externs:  make([]LineInfo, 0),
			methods:  make([]LineInfo, 0),
		},
		isTranslate: false,
		stream:      stream.New(),
		paths:       make([]string, 0),
	}
}
func (s *scanner) Translate(root string) (ok bool) {
	s.isTranslate = true
	return s.Scan(root)
}
func (s *scanner) RemoveComment(root string) (ok bool) { return s.Scan(root) }
func (s *scanner) Scan(root string) (ok bool) {
	return mylog.Error(filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".clang-format" {
			if !mylog.Error(os.Remove(path)) {
				return err
			}
		}
		s.path = path
		if !s.hasExt() {
			return err
		}
		if !s.scan() {
			return err
		}
		s.paths = append(s.paths, path)
		return nil
	}))
}
func (s *scanner) scan() (ok bool) {
	b, err := os.ReadFile(s.path)
	if !mylog.Error(err) {
		return
	}
	s.stream.Reset()
	if !mylog.Error2(s.stream.Write(b)) {
		return
	}
	var sc goScanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), s.stream.Len())
	sc.Init(file, s.stream.Bytes(), nil /* no error handler */, goScanner.ScanComments)
	body := s.stream.String()
	s.stream.Reset()
	Append := func() {
		s.words = append(s.words, s.word)
		s.tokens = append(s.tokens, s.token)
	}
	for {
		pos, tok, lit := sc.Scan()
		s.token = tok
		if tok == token.EOF {
			Append()
			break
		}
		switch {
		case tok == token.COMMENT:
			count := strings.Count(lit, "\n")
			repeat := strings.Repeat("\n", count)
			body = strings.Replace(body, lit, repeat, 1)
			//fmt.Printf("%s\t%s\t%q\n", fset.Position(pos), tok, lit)
		default:
			s.Line = fset.Position(pos).Line
			s.Text = lit
			if lit == "" {
				s.Text = tok.String()
			}
			Append()
		}
	}
	if !s.generateNoCommentFile(body) {
		return
	}
	if s.isTranslate {
		return s.translate()
	}
	c := clang.New()
	if !c.WriteClangFormatBody(s.noCommentPath()) {
		return
	}
	return c.Format(s.path)
}
func (s *scanner) generateNoCommentFile(body string) (ok bool) {
	dstDir := ""
	lines, b := tool.File().ToLines(body)
	if !b {
		return
	}
	for _, line := range lines {
		switch s.isTranslate {
		case true:
			dstDir = s.bindingPath()
			s.stream.WriteStringLn(line)
		case false:
			dstDir = s.noCommentPath()
			if strings.TrimSpace(line) == "" {
				continue
			}
			s.stream.WriteStringLn(line)
			switch {
			case line == "}":
				if !mylog.Error2(s.stream.WriteString("\n")) {
					return
				}
			}
		}
	}
	dst := filepath.Join(dstDir, s.path)
	if s.isTranslate {
		ext := filepath.Ext(dst)
		dst = dst[:len(dst)-len(ext)] + s.backExt()
	}
	if !tool.File().WriteTruncate(dst, s.stream.String()[:s.stream.Len()-1]) {
		return
	}
	s.stream.Reset()
	return true
}
func (s *scanner) translate() (ok bool) {
	if !s.FindAllBlock() {
		return
	}
	if !s.makeBlock() { //todo add object for make block
		return
	}
	if !s.translateBlock() {
		return
	}
	if !s.bindBlockType() {
		return
	}
	if !s.handlePkgOrApiName() {
		return
	}
	if !s.generateGoCodes() {
		return
	}
	s.reset()
	return true
}
func (s *scanner) reset() {
	*s = scanner{
		fileInfo: fileInfo{
			LineInfo: LineInfo{
				word: word{
					Line:  0,
					Text:  "",
					token: 0,
				},
				tokens: s.tokens[:0],
				elems:  s.elems[:0],
			},
			path:     "",
			words:    s.words[:0],
			allLines: s.allLines[:0],
			typedefs: s.typedefs[:0],
			enums:    s.enums[:0],
			defines:  s.defines[:0],
			externs:  s.externs[:0],
			methods:  s.methods[:0],
		},
		isTranslate: s.isTranslate,
		stream:      s.stream,
		paths:       s.paths,
	}
}

const (
	typedefKInd        = "typedef"
	structKInd         = "struct"
	enumKInd           = "enum"
	defineKInd         = "define"
	externKInd         = "extern"
	methodArgStartKInd = "("
	methodArgEndKInd   = ")"
	methodStartKInd    = "{"

	//replaced by token.eof
	methodEndKInd  = "}"  // .cpp .c .cxx .cpp files etc
	methodEnd2KInd = ");" // .h file
)

func (s *scanner) FindAllBlock() (ok bool) {
	if !s.MergeLineElem() {
		return
	}
	return true
	s.FindTypedefs()
	s.FindEnums()
	s.FindDefines()
	s.FindExterns()
	s.FindMethods()
	//todo clean slow lines and log them
	return true
}
func (s *scanner) MergeLineElem() (ok bool) {
	Append := func(w word) {
		s.token = w.token
		s.elems = append(s.elems, w.Text) //next
		s.Text = strings.Join(s.elems, " ")
		s.allLines = append(s.allLines, s.LineInfo) //todo  luan xu  bug
		s.elems = s.elems[:0]
	}
	isLast := func(i int) bool { return i+1 == len(s.words) }
	for i, w := range s.words {
		s.Line = w.Line
		Append(w)
		nextLIne := 0
		if isLast(i) {
			nextLIne = w.Line
		} else {
			nextLIne = s.words[i+1].Line
		}
		if s.Line != nextLIne {
			s.elems = s.elems[:len(s.elems)-1]
		}
	}
	return true
}
func (f *fileInfo) makeBlock() (ok bool) { //todo make slice
	return true
}
func (f *fileInfo) translateBlock() (ok bool) {
	//for {
	//	element := a.list.Front()
	//	if element == nil {
	//		return true
	//	}
	//	switch element.Value.(type) {
	//	case typedefStructsTYpe:
	//		typedefStructs := element.Value.(typedefStructsTYpe)
	//		for _, typedefStruct := range typedefStructs {
	//			for _, info := range typedefStruct.LineInfos {
	//				mylog.Trace(fmt.Sprint(info.Line), info.Text)
	//				if !a.bindBlockType() {
	//					return
	//				}
	//				if !a.handlePkgOrApiName() {
	//					return
	//				}
	//				if !a.generateGoCodes() {
	//					return
	//				}
	//			}
	//		}
	//	case typedefEnumsTYpe:
	//		typedefEnums := element.Value.(typedefEnumsTYpe)
	//		for _, typedefEnum := range typedefEnums {
	//			for _, info := range typedefEnum.LineInfos {
	//				mylog.Trace(fmt.Sprint(info.Line), info.Text)
	//				if !a.bindBlockType() {
	//					return
	//				}
	//				if !a.handlePkgOrApiName() {
	//					return
	//				}
	//				if !a.generateGoCodes() {
	//					return
	//				}
	//			}
	//		}
	//	}
	//	return true
	//}
	return true
}
func (s *scanner) bindBlockType() (ok bool) {
	//TODO implement me
	panic("implement me")
}
func (s *scanner) handlePkgOrApiName() (ok bool) {
	//TODO implement me
	panic("implement me")
}
func (s *scanner) generateGoCodes() (ok bool) {
	//TODO implement me
	panic("implement me")
}

// ////////////////////////////////////////////////////////////
func (s *scanner) bindingPath() string   { return "binding" }
func (s *scanner) noCommentPath() string { return "noComment" }
func (s *scanner) backExt() string       { return ".back" }
func (s *scanner) Exts() []string {
	return []string{
		".c",
		".cc",
		".cpp",
		".cppm",
		".h",
		".hh",
		".hpp",
		".ixx",
		".cs",
		".go",
		//".back",
	}
}
func (s *scanner) hasExt() (ok bool) {
	for _, e := range s.Exts() {
		if e == filepath.Ext(s.path) {
			return true
		}
	}
	return
}
