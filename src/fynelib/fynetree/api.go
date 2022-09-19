package fynetree

import (
	"github.com/ddkwork/golibrary/mylog"
	"github.com/drognisep/fynehelpers/generation"
)

type (
	Interface interface { //todo fyne help 增加api生成tree，左右视图排版增加例子到lib
		NewTree(modelRoots ...generation.TreeModel) *BranchTree

		NewRoot(title string) *Object
		RootAddChild(root *Object, title string)

		NewBranch(title string) *Object
		BranchAddChild(branch *Object, title string)

		NewNode(title string) *Object
		NodeAddChild(node *Object, title string)
	}
	api struct{}
)

func (a *api) NewTree(modelRoots ...generation.TreeModel) *BranchTree {
	return NewBranchTree(modelRoots...)
}
func (a *api) NewRoot(title string) *Object { return a.NewBranch(title) }
func (a *api) RootAddChild(root *Object, title string) {
	root.title = title
	if !mylog.Error(root.AddChild(root)) {
		return
	}

}
func (a *api) NewBranch(title string) *Object {
	return &Object{
		mod:   new(generation.BaseTreeModel),
		title: title,
	}
}
func (a *api) BranchAddChild(branch *Object, title string) {
	branch.title = title
	if !mylog.Error(branch.AddChild(branch)) {
		return
	}
}
func (a *api) NewNode(title string) *Object { return a.NewBranch(title) }
func (a *api) NodeAddChild(node *Object, title string) {
	node.title = title
	if !mylog.Error(node.AddChild(node)) {
		return
	}
}
func New() Interface {
	return &api{}
}

//func NewBranch(title string) *object {
//	return &object{
//		mod:   new(generation.BaseTreeModel),
//		title: title,
//	}
//}
//
//func AddNode(title string) *object {
//	return NewBranch(title)
//}
//
//func NewTree(modelRoots ...generation.TreeModel) *BranchTree {
//	return NewBranchTree(modelRoots...)
//}
