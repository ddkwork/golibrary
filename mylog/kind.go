package mylog

import (
	"fmt"

	"github.com/fatih/color"
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

func (o *object) printColorBody() {
	ColorBody := ""
	switch o.kind {
	case hexKind:
		ColorBody = fmt.Sprintf(colorFormat, color.FgHiBlue, o.body)
	case hexDumpKind:
		ColorBody = fmt.Sprintf(colorFormat, color.FgHiBlue, o.body)
	case jsonKind:
		ColorBody = fmt.Sprintf(colorFormat, color.FgHiBlue, o.body)
	case structKind:
		ColorBody = fmt.Sprintf(colorFormat, color.FgHiBlue, o.body)
	case infoKind:
		ColorBody = fmt.Sprintf(colorFormat, color.FgHiCyan, o.body)
	case traceKind:
		ColorBody = fmt.Sprintf(colorFormat, color.FgHiMagenta, o.body)
	case errorKind:
		ColorBody = fmt.Sprintf(colorFormat, color.FgHiRed, o.body)
	case warningKind:
		ColorBody = fmt.Sprintf(colorFormat, color.FgHiYellow, o.body)
	case successKind:
		ColorBody = fmt.Sprintf(colorFormat, color.FgHiGreen, o.body)
	}
	if o.debug {
		fmt.Println(ColorBody)
	}
}
