package mylog

import "fmt"

type kind int

const (
	hexKind kind = iota
	hexDumpKind
	jsonKind
	structKind
	infoKind
	traceKind
	successKind
	warningKind
	errorKind
)

var kindStrings = map[kind]string{
	hexKind:     " HexV ",
	hexDumpKind: " Dump ",
	jsonKind:    " Json ",
	structKind:  " Stru ",
	infoKind:    " Info ",
	traceKind:   " Trac ",
	successKind: " Succ ",
	warningKind: " Warn ",
	errorKind:   " Erro ",
}

func (k kind) String() string {
	return kindStrings[k]
}

const (
	colorFormat = "\x1b[1m\x1b[%dm%s\x1b[0m"
)

func (l *log) printColorBody(s string) {
	ColorBody := ""
	switch l.kind {
	case hexKind:
		ColorBody = fmt.Sprintf(colorFormat, FgHiBlue, s)
	case hexDumpKind:
		ColorBody = fmt.Sprintf(colorFormat, FgHiBlue, s)
	case jsonKind:
		ColorBody = fmt.Sprintf(colorFormat, FgHiBlue, s)
	case structKind:
		ColorBody = fmt.Sprintf(colorFormat, FgHiBlue, s)
	case infoKind:
		ColorBody = fmt.Sprintf(colorFormat, FgHiCyan, s)
	case traceKind:
		ColorBody = fmt.Sprintf(colorFormat, FgHiMagenta, s)
	case errorKind:
		ColorBody = fmt.Sprintf(colorFormat, FgHiRed, s)
	case warningKind:
		ColorBody = fmt.Sprintf(colorFormat, FgHiYellow, s)
	case successKind:
		ColorBody = fmt.Sprintf(colorFormat, FgHiGreen, s)
	}
	if l.debug {
		fmt.Print(ColorBody)
	}
}

type attribute int

const (
	FgHiBlack attribute = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)
