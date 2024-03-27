package widget

import (
	"cogentcore.org/core/ki"
	"cogentcore.org/core/texteditor"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/safeType"
)

// Editor https://swapoff.org/chroma/playground/
type Editor[T safeType.Type] struct {
	*texteditor.Editor
	*texteditor.Buf
	data T // not used
}

func NewEditor[T safeType.Type](parent ki.Ki, data T, tooltip string, name ...string) *Editor[T] { // data形参是为了后续的实例化接受泛型，不用指明要实例化什么类型
	if len(name) > 0 {
		tooltip += " " + name[0]
	}
	return &Editor[T]{
		Editor: texteditor.NewEditor(parent, name...).SetTooltip(tooltip),
		Buf:    texteditor.NewBuf(),
	}
}

func (e *Editor[T]) SetLogBody() {
	e.Buf.Hi.Style = "manni"
	e.Buf.Hi.Lang = "Go"
	e.Buf.SetText([]byte(mylog.Body()))
	e.Editor.SetBuf(e.Buf)
}

func (e *Editor[T]) SetData(data T) {
	e.Buf.Hi.Style = "manni" // todo mock goland color
	e.Buf.Hi.Lang = "Go"     // Nasm
	e.Buf.SetText(safeType.New(data).Bytes())
	e.Editor.SetBuf(e.Buf)
}

func (e *Editor[T]) SetLanguage(lang string) { e.Buf.Hi.Lang = lang }
func (e *Editor[T]) SetNasmLanguage()        { e.SetLanguage("Nasm") }
func (e *Editor[T]) SetGoLanguage()          { e.SetLanguage("Go") }
func (e *Editor[T]) SetCLanguage()           { e.SetLanguage("C") }
func (e *Editor[T]) SetRustLanguage()        { e.SetLanguage("Rust") }
func (e *Editor[T]) SetBashLanguage()        { e.SetLanguage("Bash") }
func (e *Editor[T]) SetCSharpLanguage()      { e.SetLanguage("CSharp") }
func (e *Editor[T]) SetJavaLanguage()        { e.SetLanguage("Java") }
func (e *Editor[T]) SetPythonLanguage()      { e.SetLanguage("Python") } // todo

//	// Code is a programming language file
//	AnyCode
//	Ada
//	Bash
//	Csh
//	C // includes C++
//	CSharp
//	D
//	Diff
//	Eiffel
//	Erlang
//	Forth
//	Fortran
//	FSharp
//	Go
//	Haskell
//	Java
//	JavaScript
//	Lisp
//	Lua
//	Makefile
//	Mathematica
//	Matlab
//	ObjC
//	OCaml
//	Pascal
//	Perl
//	Php
//	Prolog
//	Python
//	R
//	Ruby
//	Rust
//	Scala
//	Tcl
