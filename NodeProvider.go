package golibrary

import (
	"fmt"

	"github.com/google/uuid"
)

type NodeProvider[T any] interface {
	fmt.Stringer
	UUID() uuid.UUID
	Clone() T
	Kind() string
	Container() bool
	Parent() T
	SetParent(parent T)
	HasChildren() bool
	NodeChildren() []T
	SetChildren(children []T)
	Enabled() bool
	Open() bool
	Depth() int
	SetOpen(open bool)
	// CellData(columnID int, data *CellData)
	FillWithNameableKeys(m map[string]string)
	ApplyNameableKeys(m map[string]string)
}

type EditorData[T any] interface {
	CopyFrom(from T)
	ApplyTo(to T)
}

func AsNode[T any](in T) NodeProvider[T] {
	return any(in).(NodeProvider[T])
}
