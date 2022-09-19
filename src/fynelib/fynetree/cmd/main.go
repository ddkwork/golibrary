package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/ddkwork/golibrary/src/fynelib/fyneTheme"
	"github.com/ddkwork/golibrary/src/fynelib/fynetree"
	"github.com/drognisep/fynehelpers/generation"
)

func main() {
	tree := fynetree.New()
	root := tree.NewRoot("root{")
	root.AddChild(tree.NewNode("1"))
	root.AddChild(tree.NewNode("2"))

	branch1 := tree.NewBranch("branch1{")
	branch1.AddChild(tree.NewNode("1"))
	branch1.AddChild(tree.NewNode("2"))
	branch1.AddChild(tree.NewNode("3"))
	branch1.AddChild(tree.NewNode("}"))

	branch2 := tree.NewBranch("branch2{")
	branch2.AddChild(tree.NewNode("1"))
	branch2.AddChild(tree.NewNode("2"))
	branch2.AddChild(tree.NewNode("3"))
	branch2.AddChild(tree.NewNode("}"))

	root.AddChild(branch1)
	root.AddChild(branch2)
	root.AddChild(tree.NewNode("3"))
	root.AddChild(tree.NewNode("4"))
	root.AddChild(tree.NewNode("5"))
	root.AddChild(tree.NewNode("}"))

	newTree := tree.NewTree(root)
	newTree.OnTapped = func(id widget.TreeNodeID, model generation.TreeModel, event *fyne.PointEvent) {
		fmt.Println(id + " " + model.DisplayString())
	}

	a := app.NewWithID("com.rows.app")
	a.SetIcon(nil)
	fyneTheme.Dark()
	w := a.NewWindow("app")
	w.Resize(fyne.NewSize(640, 480))
	w.SetMaster()
	w.CenterOnScreen()
	newTree.OpenAllBranches() //在上面有空指针，不理解

	w.SetContent(newTree)
	w.ShowAndRun()
}
