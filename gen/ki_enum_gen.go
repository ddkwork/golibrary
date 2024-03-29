package gen

import (
	"cogentcore.org/core/gi"
	"cogentcore.org/core/icons"
	"golang.org/x/exp/constraints"
)

// Code generated by GeneratedFile enum - DO NOT EDIT.

type KiKind byte

const (
	SliceViewBaseKind KiKind = iota
	WidgetBaseKind
	InvalidKiKind
)

func ConvertInteger2KiKind[T constraints.Integer](v T) KiKind {
	return KiKind(v)
}

func (k KiKind) AssertKind(kinds string) KiKind {
	for _, kind := range k.Kinds() {
		if kinds == kind.String() {
			return kind
		}
	}
	return InvalidKiKind
}

func (k KiKind) ChooserItem() []gi.ChooserItem {
	chooserItems := make([]gi.ChooserItem, 0)
	for _, kind := range k.Kinds() {
		chooserItems = append(chooserItems, gi.ChooserItem{
			Value:           kind.String(),
			Label:           "",
			Icon:            icons.Icon(k.SvgFileName()),
			Tooltip:         kind.String(),
			Func:            nil,
			SeparatorBefore: false,
		})
	}
	return chooserItems
}

func (k KiKind) String() string {
	switch k {
	case SliceViewBaseKind:
		return "SliceViewBase"
	case WidgetBaseKind:
		return "WidgetBase"
	default:
		return "InvalidKiKind"
	}
}

func (k KiKind) Keys() []string {
	return []string{
		"SliceViewBase",
		"WidgetBase",
		"InvalidKiKind",
	}
}

func (k KiKind) Kinds() []KiKind {
	return []KiKind{
		SliceViewBaseKind,
		WidgetBaseKind,
		InvalidKiKind,
	}
}

func (k KiKind) SvgFileName() string {
	switch k {
	case SliceViewBaseKind:
		return "SliceViewBase"
	case WidgetBaseKind:
		return "WidgetBase"
	default:
		return "InvalidKiKind"
	}
}
