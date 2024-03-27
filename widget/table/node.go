package table

import (
	"fmt"
	"reflect"
	"slices"
	"sort"
	"strings"

	"github.com/google/uuid"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
)

// 不支持接口约束，不过有了泛型就够了，虽然这个场景约束类型不适合，理论上也是不应该约束类型，因为每个项目都有很多不同的结构体
// 与其约束结构体的类型和约束结构体字段类型一样不合理，而让泛型接受任何结构体或者普通类型数据是刚需
// 以前用接口接受任意类型并约束方法的可见性，但是不通用，内部方法还要断言，换成泛型限制接收了任意类型没有了方法可见性控制，但通用了，
// 与之而来的方法重写也实验失败了，不管了，反正只要安全+好用，解决方案使用回调

//对于Encoding接口,node是被外部匿名嵌套套娃的，所以要在不能在本包实现，会无限循环导致堆栈溢出，是错误的逻辑，需要在外部包实现接口
//泛型的类型约束场景:普通函数，结构体只有一个字段的结构体类型约束以及结构体的方法类型约束
//泛型方法场景：通用树形结构，因为字段Data来自外部包，无论它是结构体还是普通类型的数据，
//每次实例化一个节点都限制了单次树形对象指针的外部类型data是唯一的类型
//PS：结构体字段永远不适合约束所有字段类型，理论上也是通的，因为字段的命名和个数不确定，是由外部包决定的

//外部方法使用注意：不要吧新参各种身份的node放在结构体字段内，因为每个node只代表自己，
//它有独立的父级，孩子集，类型，以及外部结构体或任意类型的真实数据

//了解以上好处后，不应该在任何场景考虑使用接口，而应该考虑泛型方法代替接口，因为它通用和安全，外部包不再需要断言类型
//接口属于鸭子类型的方法集，能限制可见性，而泛型方法属于限制鸭子类型的无需外部包断言的方法集，虽然不能限制可见性

//todo 制作一个树形表格或者思维导图清晰的解释为什么要这么使用以及注意事项
// interface or method name here
/*
  ───────────────────────────────────────────────────────────────────────────────────────────────────
  │Method                                                                           │Description
  ───────────────────────────────────────────────────────────────────────────────────────────────────
  │doc_container                                                                    │
  ├───init_container                                                                │
  │   ├───New[T any](typeKey string, data T) (root *Node[T])                        │
  │   ├───NewContainerNode[T any](typeKey string, data T) *Node[T]                  │
  │   ├───NewNode[T any](typeKey string, data T) *Node[T]                           │
  │   ├───Root() *Node[T]                                                           │
  │   └───IsRoot() bool	                                                            │
  ├───walker_container                                                              │
  │   ├───Walk(callback func(node *Node[T]))                                        │
  │   └───WalkContainer(callback func(node *Node[T]))                               │
  ├───ctx_container                                                                 │
  │   ├───Depth() int                                                               │
  │   ├───UUID() uuid.UUID                                                          │
  │   ├───Sort(cmp func(a T, b T) bool)                                             │
  │   ├───Clone(newParent *Node[T], preserveID bool) *Node[T]                       │
  │   ├───CopyFrom(from *Node[T])                                                   │
  │   ├───ApplyTo(to *Node[T])                                                      │
  │   ├───Parent() *Node[T]                                                         │
  │   ├───SetParent(parent *Node[T])                                                │
  │   ├───Container() bool                                                          │
  │   ├───HasChildren() bool                                                        │
  │   ├───SetChildren(children []*Node[T])                                          │
  │   ├───clearUnusedFields()                                                       │
  │   ├───GetType() string                                                          │
  │   ├───SetType(t string)                                                         │
  │   ├───kind(base string) string                                                  │
  │   ├───Open() bool                                                               │
  │   ├───SetOpen(open bool)                                                        │
  │   ├───OpenAll() bool                                                            │
  │   ├───CloseAll() bool                                                           │
  │   ├───LenChildren() int                                                         │
  │   ├───LastChild()                                                               │
  │   ├───IsLastChild() bool                                                        │
  │   └───ResetChildren()                                                           │
  ├───adder_container                                                               │
  │   ├───AddContainer(container *Node[T])                                          │
  │   ├───AddContainerByData(typeKey string, data T)                                │
  │   ├───AddChildByData(typeKey string, data T)                                    │
  │   ├───AddChild(child *Node[T])                                                  │
  │   ├───InsertChildItem(parentID uuid.UUID, data T) *Node[T]                      │
  │   ├───CreateItem(parent *Node[T], data T)                                       │
  │   ├───InsertChildByID(id uuid.UUID) *Node[T]                                    │
  │   ├───InsertChildByName(id uuid.UUID) *Node[T]                                  │
  │   ├───InsertChildByIndex(id uuid.UUID) *Node[T]                                 │
  │   ├───InsertChildByType(id uuid.UUID) *Node[T]                                  │
  │   └───InsertChildByPath(id uuid.UUID) *Node[T]                                  │
  ├───remover_container                                                             │
  │   ├───RemoveByID(id uuid.UUID)                                                  │
  │   ├───RemoveSelf(id uuid.UUID)                                                  │
  │   ├───RemoveChildByID()                                                         │
  │   ├───RemoveChildByName()                                                       │
  │   ├───RemoveChildByIndex()                                                      │
  │   └───RemoveChildByPath()                                                       │
  ├───updater_container                                                             │
  │   ├───UpdateChildByID(id uuid.UUID, data T)                                     │
  │   ├───UpdateChildByName(id uuid.UUID, data T)                                   │
  │   ├───UpdateChildByIndex(id uuid.UUID, data T)                                  │
  │   ├───UpdateChildByPath(id uuid.UUID, data T)                                   │
  │   ├───UpdateChildByType(id uuid.UUID, data T)                                   │
  │   ├───SwapChild()                                                               │
  │   └───MoveChildTo()                                                             │
  ├───finder_container                                                              │
  │   ├───FindChildByID(id uuid.UUID) *Node[T]                                      │
  │   ├───FindChildByName(id uuid.UUID) *Node[T]                                    │
  │   ├───FindChildByIndex(id uuid.UUID) *Node[T]                                   │
  │   ├───FindChildByPath(id uuid.UUID) *Node[T]                                    │
  │   └───FindChildByType(id uuid.UUID) *Node[T]                                    │
  └───tableCtx_container                                                            │
  │   ├───SetRowCallback(rowCallback func(node *Node[T]) (cell []string))           │
  │   ├───String() string                                                           │
  │   ├───Document() string                                                         │
  │   ├───format(node *Node[T], s *stream.Stream)                                   │
  │   ├───Enabled() bool                                                            │
  │   ├───SetSelectedFunc()                                                         │
  │   ├───GetPath()                                                                 │
  │   ├───SetAlign()                                                                │
  │   └───MouseHandler()                                                            │
*/
type (
	Node[T any] struct { //ContainerBase
		ID       uuid.UUID `json:"-" tableview:"-"`
		MetaData T
		Type     string     `json:"type"`
		IsOpen   bool       `json:"-" tableview:"-"` // Container only
		Children []*Node[T] // Container only
		parent   *Node[T]
		root     *Node[T]
		TableCtx[T]
	}
	TableCtx[T any] struct {
		rowCallback func(node *Node[T]) (cell []string)
	}
)

const ContainerKeyPostfix = "_container"

// init 初始化
func New[T any](typeKey string, data T) (root *Node[T]) { //this is a narytree root
	node := NewContainerNode(typeKey, data)
	node.root = node
	return node
}
func (n *Node[T]) Root() *Node[T] { return n.root }
func (n *Node[T]) IsRoot() bool   { return n.parent == nil }

//func (n *Node[T]) FindRoot() (root *Node[T]) {
//	n.Walk(func(node *Node[T]) {
//		if node.parent == nil {
//			root = node
//		}
//	})
//	return
//}

func NewContainerNode[T any](typeKey string, data T) *Node[T] { return newNode(typeKey, true, data) }
func NewContainerNodes[T any](typeKeys ...string) (containers []*Node[T]) {
	containers = make([]*Node[T], 0)
	for _, key := range typeKeys {
		var data T //it is zero value
		containers = append(containers, NewContainerNode(key, data))
	}
	return
}
func NewNode[T any](data T) *Node[T] { return newNode("", false, data) }
func newNode[T any](typeKey string, isContainer bool, data T) *Node[T] {
	if isContainer {
		typeKey += ContainerKeyPostfix
	}
	return &Node[T]{
		ID:       NewUUID(),
		MetaData: data,
		Type:     typeKey,
		IsOpen:   isContainer,
		Children: make([]*Node[T], 0),
		parent:   nil,
	}
}

func NewUUID() uuid.UUID {
	id, err := uuid.NewRandom()
	if !mylog.Error(err) {
		return uuid.UUID{}
	}
	return id
}

// 深度遍历和广度遍历

func (n *Node[T]) Walk(callback func(node *Node[T])) { //this method can not be call reaped
	callback(n)
	for _, child := range n.Children {
		child.Walk(callback)
	}
}
func (n *Node[T]) WalkContainer(callback func(node *Node[T])) { //this method can not be call reaped
	queue := []*Node[T]{n}
	for len(queue) > 0 {
		node := queue[0]  //出栈
		queue = queue[1:] //每回调一次相当于出栈一次，出栈就要干掉当前的elem。 堆栈只有承上的特性，而树形则有承上启下的特性，上至parent，下至children
		callback(node)
		for _, child := range node.Children {
			queue = append(queue, child) //入栈
		}
	}
}

func (n *Node[T]) Depth() int {
	count := 0
	p := n.parent
	for p != nil {
		count++
		p = p.parent
	}
	return count
}
func (n *Node[T]) UUID() uuid.UUID { return n.ID }
func (n *Node[T]) Sort(cmp func(a, b T) bool) {
	sort.SliceStable(n.Children, func(i, j int) bool {
		return cmp(n.Children[i].MetaData, n.Children[j].MetaData)
	})
	for _, child := range n.Children {
		child.Sort(cmp)
	}
}
func (n *Node[T]) Clone() (newNode *Node[T]) {
	if n.Container() {
		return NewContainerNode(n.Type, n.MetaData)
	}
	return NewNode(n.MetaData)
}
func (n *Node[T]) CopyFrom(from *Node[T])    { *n = *from }
func (n *Node[T]) ApplyTo(to *Node[T])       { *to = *n }
func (n *Node[T]) Parent() *Node[T]          { return n.parent }   //todo export
func (n *Node[T]) SetParent(parent *Node[T]) { n.parent = parent } //todo remove
func (n *Node[T]) Container() bool           { return strings.HasSuffix(n.Type, ContainerKeyPostfix) }
func (n *Node[T]) HasChildren() bool         { return n.Container() && n.LenChildren() > 0 }
func (n *Node[T]) SetChildren(children []*Node[T]) {
	n.Children = children
}
func (n *Node[T]) clearUnusedFields() {
	if !n.Container() {
		n.Children = nil
		n.IsOpen = false
	}
}
func (n *Node[T]) GetType() string  { return n.Type }
func (n *Node[T]) SetType(t string) { n.Type = t }
func (n *Node[T]) kind(base string) string { //todo delete
	if n.Container() {
		return fmt.Sprintf("%s Container", base)
	}
	return base
}
func (n *Node[T]) Open() bool        { return n.IsOpen && n.Container() }
func (n *Node[T]) SetOpen(open bool) { n.IsOpen = open && n.Container() }
func (n *Node[T]) OpenAll() {
	n.WalkContainer(func(node *Node[T]) {
		if node.Container() {
			node.SetOpen(true)
		}
	})
}
func (n *Node[T]) CloseAll() {
	n.WalkContainer(func(node *Node[T]) {
		if node.Container() {
			node.SetOpen(false)
		}
	})
}
func (n *Node[T]) LenChildren() int {
	return len(n.Children)
}
func (n *Node[T]) LastChild() (lastChild *Node[T]) {
	if n.IsRoot() {
		return n.Children[len(n.Children)-1]
	}
	return n.parent.Children[len(n.parent.Children)-1]
}
func (n *Node[T]) IsLastChild() bool {
	return n.LastChild() == n
}
func (n *Node[T]) ResetChildren() { n.Children = n.Children[:0] }

// AddContainerByData 增 todo
func (n *Node[T]) AddContainerByData(typeKey string, data T) (newContainer *Node[T]) { //我们需要返回新的容器节点用于递归填充它的孩子节点，用例是explorer文件资源管理器
	newContainer = NewContainerNode(typeKey, data)
	n.AddChild(newContainer)
	return
}
func (n *Node[T]) AddChildByData(data T) { n.AddChild(NewNode(data)) }
func (n *Node[T]) AddChildByDatas(datas ...T) {
	for _, data := range datas {
		n.AddChild(NewNode(data))
	}
}
func (n *Node[T]) AddChild(child *Node[T]) {
	child.parent = n
	n.Children = append(n.Children, child)
}

// InsertChildItem 下面这些增加方法没什么用，这种应该在外部的结构体遍历
func (n *Node[T]) InsertChildItem(parentID uuid.UUID, data T) *Node[T] { //todo remove 改为ByData
	parent := n.FindChildByID(parentID)
	if parent == nil {
		return n
	}
	child := NewNode(data)
	parent.AddChild(child)
	return child
}
func (n *Node[T]) CreateItem(parent *Node[T], data T) { //todo remove  改为ByData //其实就是一个nodeNode的事，只是说Item在不停的业务场景中把node形象化，可读性更好
	child := NewNode(data)
	parent.AddChild(child)
}
func (n *Node[T]) InsertChildByID(id uuid.UUID) *Node[T] { //todo remove 改为ByData
	//slices.Insert(n.Children,0,1)
	//TODO implement me
	panic("implement me")
}
func (n *Node[T]) InsertChildByName(id uuid.UUID) *Node[T] { //todo remove 改为ByData
	//TODO implement me
	panic("implement me")
}
func (n *Node[T]) InsertChildByIndex(id uuid.UUID) *Node[T] { //todo remove 改为ByData
	//TODO implement me
	panic("implement me")
}
func (n *Node[T]) InsertChildByType(id uuid.UUID) *Node[T] { //todo remove 改为ByData
	//TODO implement me
	panic("implement me")
}
func (n *Node[T]) InsertChildByPath(id uuid.UUID) *Node[T] { //todo remove 改为ByData
	//TODO implement me
	panic("implement me")
}

// 删
func (n *Node[T]) RemoveByID(id uuid.UUID) { //todo merge  改为ByData
	for i, child := range n.Children {
		if child.ID == id {
			n.Children = slices.Delete(n.Children, i, i+1)
			break
		}
	}
}
func (n *Node[T]) RemoveSelf(id uuid.UUID) { //todo remove 改为ByData
	//TODO implement me
	panic("implement me")
}
func (n *Node[T]) RemoveChildByID() { //todo remove 改为ByData
	//TODO implement me
	panic("implement me")
}
func (n *Node[T]) RemoveChildByName() { //todo remove 改为ByData
	//TODO implement me
	panic("implement me")
}
func (n *Node[T]) RemoveChildByIndex() { //todo remove 改为ByData
	//TODO implement me
	panic("implement me")
}
func (n *Node[T]) RemoveChildByPath() { //todo remove 改为ByData
	//TODO implement me
	panic("implement me")
}

// 改
func (n *Node[T]) UpdateChildByID(id uuid.UUID, data T) {
	node := n.FindChildByID(id)
	if node != nil {
		node.MetaData = data
	}
}
func (n *Node[T]) UpdateChildByName(id uuid.UUID, data T) { //todo remove 改为ByData
	//TODO implement me
	panic("implement me")
}
func (n *Node[T]) UpdateChildByIndex(id uuid.UUID, data T) { //todo remove 改为ByData
	//TODO implement me
	panic("implement me")
}
func (n *Node[T]) UpdateChildByPath(id uuid.UUID, data T) { //todo remove 改为ByData
	//TODO implement me
	panic("implement me")
}
func (n *Node[T]) UpdateChildByType(id uuid.UUID, data T) { //todo remove 改为ByData
	//TODO implement me
	panic("implement me")
}
func (n *Node[T]) SwapChild() { //todo remove
	//TODO implement me
	panic("implement me")
}
func (n *Node[T]) MoveChildTo() { //todo remove
	//TODO implement me
	panic("implement me")
}

// 查
func (n *Node[T]) FindChildByID(id uuid.UUID) *Node[T] {
	if n.ID == id {
		return n
	}
	for _, child := range n.Children {
		found := child.FindChildByID(id)
		if found != nil {
			return found
		}
	}
	return nil
}
func (n *Node[T]) FindChildByName(id uuid.UUID) *Node[T] { //todo remove 改为ByData
	//TODO implement me
	panic("implement me")
}
func (n *Node[T]) FindChildByIndex(id uuid.UUID) *Node[T] { //todo remove 改为ByData
	//TODO implement me
	panic("implement me")
}
func (n *Node[T]) FindChildByPath(id uuid.UUID) *Node[T] { //todo remove 改为ByData
	//TODO implement me
	panic("implement me")
}
func (n *Node[T]) FindChildByType(id uuid.UUID) *Node[T] { //todo remove 改为ByData
	//TODO implement me
	panic("implement me")
}

// SetRowCallback 表格相关
// container节点自定义单元格内容
// container节点Sum，比如gb mb kb等
// 此外，应注意排序方法是否处理了单元格的单位造成排序错误，比如gb mb kb等
func (t *TableCtx[T]) SetRowCallback(rowCallback func(node *Node[T]) (cell []string)) {
	t.rowCallback = rowCallback
}

func (n *Node[T]) String() string { //这里都是最后才调用,它每次都是顶层的，能得到root，所以不用调用findRoot
	s := stream.New("")
	n.format(n.root, s)
	return s.String()
}
func (n *Node[T]) Document() string { //这里都是最后才调用,它每次都是顶层的，能得到root，所以不用调用findRoot
	s := stream.New("")
	s.WriteStringLn("// interface or method name here")
	s.WriteStringLn("/*")
	lines, ok := stream.New(n.String()).ToLines()
	if !ok {
		return ""
	}
	for _, line := range lines {
		s.WriteStringLn("  " + line)
	}
	s.WriteStringLn("*/")
	return s.String()
}
func (n *Node[T]) format(node *Node[T], s *stream.Stream) {
	var header []string //第一列为树形层级列
	fields := reflect.VisibleFields(reflect.TypeOf(n.MetaData))
	for _, field := range fields {
		if field.Name == "Node" {
			break
		}
		header = append(header, field.Name)
	}

	//1 渲染层级，递归格式化树节点
	const (
		indent          = "│   "
		childPrefix     = "├───"
		lastChildPrefix = "└───"
	)

	type table struct { //todo 将此结构合并到node支持动态插入行并刷新列宽，并发测试等
		columns []struct {
			width int
			data  []string
		}
		rows [][]string
	}
	t := table{
		columns: []struct {
			width int
			data  []string
		}{},
		rows: [][]string{},
	}
	t.rows = make([][]string, 1)

	node.Walk(func(node *Node[T]) { // 从根节点开始遍历
		HierarchicalColumn := ""

		depth := node.Depth() - 1 // 根节点深度为0，每一层向下递增
		for i := 0; i < depth; i++ {
			//s.WriteString(indent) // 添加缩进
			HierarchicalColumn += indent
		}

		if node.IsRoot() { // 添加节点前缀
			HierarchicalColumn = "│" + HierarchicalColumn
		} else if node.parent != nil && !node.IsLastChild() {
			HierarchicalColumn += childPrefix

		} else if node.parent != nil && node.IsLastChild() {
			HierarchicalColumn += lastChildPrefix

		}
		t.rows[0] = header //表头为第一行

		for _, cell := range header { // 设置第一行为表头并计算列宽
			t.columns = append(t.columns, struct {
				width int
				data  []string
			}{width: len(cell)})
		}

		//2 渲染行，添加节点数据
		//s.WriteString(fmt.Sprintf("[%v] Type: %v, Open: %v, MetaData: %v\n", node.UUID(), node.Type, node.IsOpen, node.MetaData))
		if n.rowCallback != nil {
			cells := n.rowCallback(node) //获取每行的单元格数据
			if len(cells) == 0 {         //快速测试模式，业务模型还没建立好，树形还没准备好久跑单元测试的情况
				return
			}
			cells[0] = HierarchicalColumn + cells[0] //
			//mylog.Error(HierarchicalColumn)
			HierarchicalColumn = ""
			t.rows = append(t.rows, cells) // 否则，添加行并更新列宽
			for i, cell := range cells {
				if len(cell) > t.columns[i].width {
					t.columns[i].width = len(cell)
				}
			}
		}
	})

	for index, row := range t.rows {
		if index == 0 {
			fnFmtHeader := func() (h string) {
				for i, cell := range row {
					if i < len(t.columns)-1 {
						if i == 0 {
							indentStr := fmt.Sprintf("│%-*s ", t.columns[i].width-1, cell) //为什么要-1？
							h += indentStr
						} else {
							indentStr := fmt.Sprintf("│%-*s ", t.columns[i].width, cell)
							h += indentStr
						}
					}
				}
				return
			}
			fmtHeader := fnFmtHeader()
			s.WriteStringLn(strings.Repeat("─", len(fmtHeader))) //这个表头的矩形有点糟糕
			s.WriteStringLn(fmtHeader)
			s.WriteString(strings.Repeat("─", len(fmtHeader))) //todo 这里为什么不能换行？
			s.NewLine()
			continue
		}
		for i, cell := range row {
			if i < len(t.columns)-1 {
				indentStr := fmt.Sprintf("%-*s ", t.columns[i].width, cell) //层级列已经有了层级文本了，不需要填充
				if i > 0 {
					indentStr = fmt.Sprintf("│%-*s ", t.columns[i].width, cell) //层级列之外需要列分隔符
				}
				s.WriteString(indentStr)
			}
		}
		s.NewLine()
	}

}
func (n *Node[T]) Enabled() bool    { return true } //todo 测试gui表格需要？
func (n *Node[T]) SetSelectedFunc() {}
func (n *Node[T]) SetIndent()       {}
func (n *Node[T]) GetPath()         {} //展示层级路径 //todo remove
func (n *Node[T]) SetAlign()        {} //设置单元格对齐方式
func (n *Node[T]) MouseHandler()    {}
