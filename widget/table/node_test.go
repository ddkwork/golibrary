package table_test

import (
	"reflect"
	"testing"

	"github.com/ddkwork/golibrary/gen"
	"github.com/ddkwork/golibrary/widget/table"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream/cmd"
)

func TestName(t *testing.T) {
	t.Skip()
	cmd.CheckLoopvarAndNilPoint()
}

type doc struct {
	Method      string
	Description string
}

func makeMethods(methods ...string) (docs []doc) {
	for _, method := range methods {
		docs = append(docs, doc{
			Method:      method,
			Description: "",
		})
	}
	return
}

func TestDoc(t *testing.T) {
	gen.New().Enum("tableDoc", []string{
		"init",
		"walker",
		"ctx",
		"adder",
		"remover",
		"updater",
		"finder",
		"tableCtx",
	}, nil)
	root := table.New("doc", doc{})
	root.SetRowCallback(func(node *table.Node[doc]) (cell []string) {
		if node.Container() {
			node.MetaData.Method = node.Type
		}
		return []string{
			node.MetaData.Method,
			node.MetaData.Description,
		}
	})

	containers := table.NewContainerNodes[doc](table.InvalidTableDocKind.Keys()...)
	for i, container := range containers {
		root.AddChild(container)
		k := table.ConvertInteger2TableDocKind(i)

		// todo gen,如何正确的填充？列编辑模式，空下多行，输入双引号和逗号，最后再粘贴方法，空格就去掉了，
		//  方法列表的话让鼠标在node结构体上就会提取所有方法
		var _ = []string{ //new methods
			"Root() *Node[T]",
			"IsRoot() bool",
			"Walk(callback func(node *Node[T]))",
			"WalkContainer(callback func(node *Node[T]))",
			"Depth() int",
			"UUID() uuid.UUID",
			"Sort(cmp func(a T, b T) bool)",
			"Clone() (newNode *Node[T])",
			"CopyFrom(from *Node[T])",
			"ApplyTo(to *Node[T])",
			"Parent() *Node[T]",
			"SetParent(parent *Node[T])",
			"Container() bool",
			"HasChildren() bool",
			"SetChildren(children []*Node[T])",
			"clearUnusedFields()",
			"GetType() string",
			"SetType(t string)",
			"kind(base string) string",
			"Open() bool",
			"SetOpen(open bool)",
			"OpenAll()",
			"CloseAll()",
			"LenChildren() int",
			"LastChild() (lastChild *Node[T])",
			"IsLastChild() bool",
			"ResetChildren()",
			"AddContainerByData(typeKey string, data T) (newContainer *Node[T])",
			"AddChildByData(data T)",
			"AddChildByDatas(datas ...T)",
			"AddChild(child *Node[T])",
			"InsertChildItem(parentID uuid.UUID, data T) *Node[T]",
			"CreateItem(parent *Node[T], data T)",
			"InsertChildByID(id uuid.UUID) *Node[T]",
			"InsertChildByName(id uuid.UUID) *Node[T]",
			"InsertChildByIndex(id uuid.UUID) *Node[T]",
			"InsertChildByType(id uuid.UUID) *Node[T]",
			"InsertChildByPath(id uuid.UUID) *Node[T]",
			"RemoveByID(id uuid.UUID)",
			"RemoveSelf(id uuid.UUID)",
			"RemoveChildByID()",
			"RemoveChildByName()",
			"RemoveChildByIndex()",
			"RemoveChildByPath()",
			"UpdateChildByID(id uuid.UUID, data T)",
			"UpdateChildByName(id uuid.UUID, data T)",
			"UpdateChildByIndex(id uuid.UUID, data T)",
			"UpdateChildByPath(id uuid.UUID, data T)",
			"UpdateChildByType(id uuid.UUID, data T)",
			"SwapChild()",
			"MoveChildTo()",
			"FindChildByID(id uuid.UUID) *Node[T]",
			"FindChildByName(id uuid.UUID) *Node[T]",
			"FindChildByIndex(id uuid.UUID) *Node[T]",
			"FindChildByPath(id uuid.UUID) *Node[T]",
		}
		switch k {
		case table.InvalidTableDocKind:
		case table.InitKind:
			docs := makeMethods(
				"New[T any](typeKey string, data T) (root *Node[T])            ",
				"NewContainerNode[T any](typeKey string, data T) *Node[T]      ",
				"NewNode[T any](typeKey string, data T) *Node[T]               ",
				"Root() *Node[T]                                               ",
				"IsRoot() bool	                                               ",
			)
			container.AddChildByDatas(docs...)
		case table.WalkerKind:
			docs := makeMethods(
				"Walk(callback func(node *Node[T]))",
				"WalkContainer(callback func(node *Node[T]))",
			)
			container.AddChildByDatas(docs...)
		case table.CtxKind:
			docs := makeMethods(
				"Depth() int                            ",
				"UUID() uuid.UUID                                        ",
				"Sort(cmp func(a T, b T) bool)                           ",
				"Clone(newParent *Node[T], preserveID bool) *Node[T]     ",
				"CopyFrom(from *Node[T])                                 ",
				"ApplyTo(to *Node[T])                                    ",
				"Parent() *Node[T]                                       ",
				"SetParent(parent *Node[T])                              ",
				"Container() bool                                        ",
				"HasChildren() bool                                      ",
				"SetChildren(children []*Node[T])                        ",
				"clearUnusedFields()                                     ",
				"GetType() string                                        ",
				"SetType(t string)                                       ",
				"kind(base string) string                                ",
				"Open() bool                                             ",
				"SetOpen(open bool)                                      ",
				"OpenAll() bool                                          ",
				"CloseAll() bool                                         ",
				"LenChildren() int                                       ",
				"LastChild()                                             ",
				"IsLastChild() bool                                      ",
				"ResetChildren()                                         ",
			)
			container.AddChildByDatas(docs...)
		case table.AdderKind:
			docs := makeMethods(
				"AddContainer(container *Node[T])                        ",
				"AddContainerByData(typeKey string, data T)              ",
				"AddChildByData(typeKey string, data T)                  ",
				"AddChild(child *Node[T])                                ",
				"InsertChildItem(parentID uuid.UUID, data T) *Node[T]    ",
				"CreateItem(parent *Node[T], data T)                     ",
				"InsertChildByID(id uuid.UUID) *Node[T]                  ",
				"InsertChildByName(id uuid.UUID) *Node[T]                ",
				"InsertChildByIndex(id uuid.UUID) *Node[T]               ",
				"InsertChildByType(id uuid.UUID) *Node[T]                ",
				"InsertChildByPath(id uuid.UUID) *Node[T]                ",
			)
			container.AddChildByDatas(docs...)
		case table.RemoverKind:
			docs := makeMethods(
				"RemoveByID(id uuid.UUID)  ",
				"RemoveSelf(id uuid.UUID)  ",
				"RemoveChildByID()         ",
				"RemoveChildByName()       ",
				"RemoveChildByIndex()      ",
				"RemoveChildByPath()       ",
			)
			container.AddChildByDatas(docs...)
		case table.UpdaterKind:
			docs := makeMethods(
				"UpdateChildByID(id uuid.UUID, data T)   ",
				"UpdateChildByName(id uuid.UUID, data T) ",
				"UpdateChildByIndex(id uuid.UUID, data T)",
				"UpdateChildByPath(id uuid.UUID, data T) ",
				"UpdateChildByType(id uuid.UUID, data T) ",
				"SwapChild()                             ",
				"MoveChildTo()                           ",
			)
			container.AddChildByDatas(docs...)
		case table.FinderKind:
			docs := makeMethods(
				"FindChildByID(id uuid.UUID) *Node[T]   ",
				"FindChildByName(id uuid.UUID) *Node[T] ",
				"FindChildByIndex(id uuid.UUID) *Node[T]",
				"FindChildByPath(id uuid.UUID) *Node[T] ",
				"FindChildByType(id uuid.UUID) *Node[T] ",
			)
			container.AddChildByDatas(docs...)
		case table.TableCtxKind:
			docs := makeMethods(
				"SetRowCallback(rowCallback func(node *Node[T]) (cell []string))",
				"String() string                                              ",
				"Document() string                                            ",
				"format(node *Node[T], s *stream.Stream)                      ",
				"Enabled() bool                                               ",
				"SetSelectedFunc()                                            ",
				"GetPath()                                                    ",
				"SetAlign()                                                   ",
				"MouseHandler()                                               ",
			)
			container.AddChildByDatas(docs...)
		}
	}
	println(root.Document())
}

func Test_mock(t *testing.T) {
	type (
		obj struct {
			Index int
			Name  string
		}
	)
	o := obj{
		Index: 9,
		Name:  "ppp",
	}
	root := table.New("", o)

	child1 := table.NewNode(o)
	child2 := table.NewNode(o)
	child3 := table.NewNode(o)

	root.AddChild(child1)
	root.AddChild(child2)
	root.AddChild(child3)

	grandchild1 := table.NewNode(o)
	grandchild2 := table.NewNode(o)

	child1.AddChild(grandchild1)
	child1.AddChild(grandchild2)

	println(root.String())

	root.Walk(func(node *table.Node[obj]) { //深度遍历
		mylog.Struct(node.MetaData)
	})

	root.WalkContainer(func(node *table.Node[obj]) { //广度遍历
		mylog.Struct(node.MetaData)
	})

	root.Sort(func(a, b obj) bool {
		return a.Index < b.Index
	})

	println(root.String())

	root.RemoveByID(child2.ID)

	println(root.String())

	root.UpdateChildByID(grandchild1.ID, o)

	println(root.String())
}

func Test_mock2(t *testing.T) {
	type field struct {
		Name   string
		Number int
		Depth  int
		K      reflect.Kind
		Value  any
	}
	root := table.New("", field{
		Name:   "x",
		Number: 0,
		Depth:  0,
		K:      0,
		Value:  nil,
	})
	Binary1 := table.NewNode(field{
		Name:   "y",
		Number: 0,
		Depth:  0,
		K:      0,
		Value:  "game/system/session/info",
	})
	root.AddChild(Binary1)

	println(root.String())
}
