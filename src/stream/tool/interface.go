package tool

import (
	"github.com/ddkwork/golibrary/src/stream/indent"
	"github.com/ddkwork/golibrary/src/stream/swap"
	"github.com/ddkwork/golibrary/src/stream/tool/file"
	"github.com/ddkwork/golibrary/src/stream/tool/flot"
	"github.com/ddkwork/golibrary/src/stream/tool/gzip"
	"github.com/ddkwork/golibrary/src/stream/tool/path"
	"github.com/ddkwork/golibrary/src/stream/tool/platform"
	"github.com/ddkwork/golibrary/src/stream/tool/random"
	"github.com/ddkwork/golibrary/src/stream/tool/regexp"
	"github.com/ddkwork/golibrary/src/stream/tool/strconv"
	"github.com/ddkwork/golibrary/src/stream/tool/time"
	"github.com/ddkwork/golibrary/src/stream/tool/unicode"
	"github.com/ddkwork/golibrary/src/stream/tool/version"
)

type (
	Interface interface {
		VerSion() version.Interface
		Path() path.Interface
		File() file.Interface
		Flot() flot.Interface
		Gzip() gzip.Interface
		Random() random.Interface
		Swap() swap.Interface
		Time() time.Interface
		Platform() platform.Interface
		Indent() indent.Interface
		Regexp() regexp.Interface
		Strconv() strconv.Interface
		Unicode() unicode.Interface
	}
	object struct {
		verSion       version.Interface
		path          path.Interface
		file          file.Interface
		flot          flot.Interface
		gzip          gzip.Interface
		random        random.Interface
		swap          swap.Interface
		time          time.Interface
		platform      platform.Interface
		indent        indent.Interface
		regexp        regexp.Interface
		strconv       strconv.Interface
		unicodeString unicode.Interface
	}
)

func (o *object) Unicode() unicode.Interface   { return o.unicodeString }
func (o *object) Strconv() strconv.Interface   { return o.strconv }
func (o *object) Regexp() regexp.Interface     { return o.regexp }
func (o *object) VerSion() version.Interface   { return o.verSion }
func (o *object) Path() path.Interface         { return o.path }
func (o *object) File() file.Interface         { return o.file }
func (o *object) Flot() flot.Interface         { return o.flot }
func (o *object) Gzip() gzip.Interface         { return o.gzip }
func (o *object) Random() random.Interface     { return o.random }
func (o *object) Swap() swap.Interface         { return o.swap }
func (o *object) Time() time.Interface         { return o.time }
func (o *object) Platform() platform.Interface { return o.platform }
func (o *object) Indent() indent.Interface     { return o.indent }
func New() Interface {
	return &object{
		verSion:       version.New(),
		path:          path.New(),
		file:          file.New(),
		flot:          flot.New(),
		gzip:          gzip.New(),
		random:        random.New(),
		swap:          swap.New(),
		time:          time.New(),
		platform:      platform.New(),
		indent:        indent.New(),
		regexp:        regexp.New(),
		strconv:       strconv.New(),
		unicodeString: unicode.New(),
	}
}

var Default = New()

func UnicodeString() unicode.Interface { return Default.Unicode() }
func Strconv() strconv.Interface       { return Default.Strconv() }
func Regexp() regexp.Interface         { return Default.Regexp() }
func VerSion() version.Interface       { return Default.VerSion() }
func Path() path.Interface             { return Default.Path() }
func File() file.Interface             { return Default.File() }
func Flot() flot.Interface             { return Default.Flot() }
func Gzip() gzip.Interface             { return Default.Gzip() }
func Random() random.Interface         { return Default.Random() }
func Swap() swap.Interface             { return Default.Swap() }
func Time() time.Interface             { return Default.Time() }
func Platform() platform.Interface     { return Default.Platform() }
func Indent() indent.Interface         { return Default.Indent() }
