package mylog

import (
	"fmt"
)

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

func (k kind) String() string {
	fnFmtLevel := func(l string) string { return fmt.Sprintf("%8s ->", l) }
	switch k {
	case hexKind:
		return fnFmtLevel("Hex")
	case hexDumpKind:
		return fnFmtLevel("HexDump")
	case jsonKind:
		return fnFmtLevel("Json")
	case structKind:
		return fnFmtLevel("Struct")
	case infoKind:
		return fnFmtLevel("Info")
	case traceKind:
		return fnFmtLevel("Trace")
	case successKind:
		return fnFmtLevel("Success")
	case warningKind:
		return fnFmtLevel("Warning")
	case errorKind:
		return fnFmtLevel("Error")
	}
	return ""
}

const (
	colorFormat = "\x1b[1m\x1b[%dm%s\x1b[0m"
)

func (l *log) printColorBody() {
	ColorBody := ""
	switch l.kind {
	case hexKind:
		ColorBody = fmt.Sprintf(colorFormat, FgHiBlue, l.row)
	case hexDumpKind:
		ColorBody = fmt.Sprintf(colorFormat, FgHiBlue, l.row)
	case jsonKind:
		ColorBody = fmt.Sprintf(colorFormat, FgHiBlue, l.row)
	case structKind:
		ColorBody = fmt.Sprintf(colorFormat, FgHiBlue, l.row)
	case infoKind:
		ColorBody = fmt.Sprintf(colorFormat, FgHiCyan, l.row)
	case traceKind:
		ColorBody = fmt.Sprintf(colorFormat, FgHiMagenta, l.row)
	case errorKind:
		ColorBody = fmt.Sprintf(colorFormat, FgHiRed, l.row)
	case warningKind:
		ColorBody = fmt.Sprintf(colorFormat, FgHiYellow, l.row)
	case successKind:
		ColorBody = fmt.Sprintf(colorFormat, FgHiGreen, l.row)
	}
	if l.debug {
		fmt.Println(ColorBody)
	}
}

// "github.com/fatih/color"
// attribute defines a single SGR Code

type attribute int

// const escape = "\x1b"

// Base attributes
//const (
//	Reset attribute = iota
//	Bold
//	Faint
//	Italic
//	Underline
//	BlinkSlow
//	BlinkRapid
//	ReverseVideo
//	Concealed
//	CrossedOut
//)

//const (
//	ResetBold attribute = iota + 22
//	ResetItalic
//	ResetUnderline
//	ResetBlinking
//	_
//	ResetReversed
//	ResetConcealed
//	ResetCrossedOut
//)

//var mapResetAttributes map[attribute]attribute = map[attribute]attribute{
//	Bold:         ResetBold,
//	Faint:        ResetBold,
//	Italic:       ResetItalic,
//	Underline:    ResetUnderline,
//	BlinkSlow:    ResetBlinking,
//	BlinkRapid:   ResetBlinking,
//	ReverseVideo: ResetReversed,
//	Concealed:    ResetConcealed,
//	CrossedOut:   ResetCrossedOut,
//}

// Foreground text colors
//const (
//	FgBlack attribute = iota + 30
//	FgRed
//	FgGreen
//	FgYellow
//	FgBlue
//	FgMagenta
//	FgCyan
//	FgWhite
//
//	// used internally for 256 and 24-bit coloring
//	foreground
//)

// Foreground Hi-Intensity text colors
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

//
//// Background text colors
//const (
//	BgBlack attribute = iota + 40
//	BgRed
//	BgGreen
//	BgYellow
//	BgBlue
//	BgMagenta
//	BgCyan
//	BgWhite
//
//	// used internally for 256 and 24-bit coloring
//	background
//)
//
//// Background Hi-Intensity text colors
//const (
//	BgHiBlack attribute = iota + 100
//	BgHiRed
//	BgHiGreen
//	BgHiYellow
//	BgHiBlue
//	BgHiMagenta
//	BgHiCyan
//	BgHiWhite
//)
