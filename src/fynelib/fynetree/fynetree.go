package fynetree

import (
	"fyne.io/fyne/v2"
	"github.com/drognisep/fynehelpers/generation"
)

type (
	_Interface interface {
		generation.TreeModel
	}
	Object struct {
		mod   *generation.BaseTreeModel
		title string
	}
)

func (o *Object) DisplayIcon() fyne.Resource {
	return o.mod.DisplayIcon()
}

func (o *Object) DisplayString() string {
	return o.title
}

func (o *Object) Children() []generation.TreeModel {
	return o.mod.Children()
}

func (o *Object) AddChild(model generation.TreeModel) error {
	return o.mod.AddChild(model)
}

func (o *Object) AddChildAt(i int, model generation.TreeModel) error {
	return o.mod.AddChildAt(i, model)
}

func (o *Object) RemoveChild() generation.TreeModel {
	return o.mod.RemoveChild()
}

func (o *Object) RemoveChildAt(i int) generation.TreeModel {
	return o.mod.RemoveChildAt(i)
}
