package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/ddkwork/golibrary/mylog"
	"reflect"
	"sync"
)

func StructToTree(v any) (*TreeNode, error) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return buildTree("root", val)
}

func buildTree(name string, val reflect.Value) (*TreeNode, error) {
	node := &TreeNode{
		Name: name,
		//Type:  val.Type(),//todo: 类型信息丢失
		Value: nil,
	}

	switch val.Kind() {
	case reflect.Struct:
		for i := range val.NumField() {
			field := val.Type().Field(i)
			child := mylog.Check2(buildTree(field.Name, val.Field(i)))

			node.Children = append(node.Children, child)
		}

	case reflect.Slice, reflect.Array:
		for i := range val.Len() {
			child := mylog.Check2(buildTree(fmt.Sprintf("[%d]", i), val.Index(i)))

			node.Children = append(node.Children, child)
		}

	case reflect.Ptr:
		if !val.IsNil() {
			child := mylog.Check2(buildTree(name, val.Elem()))

			node.Children = append(node.Children, child)
		}

	default:
		mylog.Todo("Unknow Type")
		//node.Value = val.Interface()
	}
	return node, nil
}

func TreeToStruct(node *TreeNode, out any) error {
	outVal := reflect.ValueOf(out)
	if outVal.Kind() != reflect.Ptr || outVal.IsNil() {
		return fmt.Errorf("输出必须为非空指针")
	}
	return decodeNode(node, outVal.Elem())
}

func decodeNode(node *TreeNode, outVal reflect.Value) error {
	//if node.Type != nil && outVal.Type() != node.Type {
	//	return fmt.Errorf("类型不匹配: 预期 %v, 实际 %v", node.Type, outVal.Type())
	//}
	mylog.Todo("类型信息丢失")
	switch outVal.Kind() {
	case reflect.Struct:
		for _, child := range node.Children {
			field := outVal.FieldByName(child.Name)
			if !field.IsValid() {
				continue
			}
			mylog.Check(decodeNode(child, field))
		}

	case reflect.Slice:
		outVal.Set(reflect.MakeSlice(outVal.Type(), len(node.Children), len(node.Children)))
		for i, child := range node.Children {
			elem := outVal.Index(i)
			decodeNode(child, elem)
		}

	case reflect.Ptr:
		if outVal.IsNil() {
			outVal.Set(reflect.New(outVal.Type().Elem()))
		}
		return decodeNode(node.Children[0], outVal.Elem())

	default:
		if node.Value != nil {
			val := reflect.ValueOf(node.Value)
			if val.Type().ConvertibleTo(outVal.Type()) {
				outVal.Set(val.Convert(outVal.Type()))
			}
		}
	}
	return nil
}

// TreeNode 树节点定义
type TreeNode struct {
	Name     string       // 字段名/元素索引（如 "Age" 或 "[1]"）
	Type     reflect.Kind // 数据类型
	Value    []byte       // 原始字节值（避免反射开销）
	Children []*TreeNode  // 子节点
}

// Codec 编解码器
type Codec struct {
	typeCache sync.Map // 类型缓存 map[reflect.Type]typeInfo
}

type typeInfo struct {
	fields []fieldInfo
}

type fieldInfo struct {
	name string
	idx  int
	typ  reflect.Type
}

// NewCodec 创建编解码器
func NewCodec() *Codec {
	return &Codec{}
}

// Encode 结构体→树
func (c *Codec) Encode(v any) (*TreeNode, error) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return c.encodeNode("root", val)
}

func (c *Codec) encodeNode(name string, val reflect.Value) (*TreeNode, error) {
	node := &TreeNode{Name: name, Type: val.Kind()}

	switch val.Kind() {
	case reflect.Struct:
		// 缓存类型信息加速
		ti, _ := c.typeCache.LoadOrStore(val.Type(), c.buildTypeInfo(val.Type()))
		for _, fi := range ti.(typeInfo).fields {
			child := mylog.Check2(c.encodeNode(fi.name, val.Field(fi.idx)))

			node.Children = append(node.Children, child)
		}

	case reflect.Slice, reflect.Array:
		for i := range val.Len() {
			child := mylog.Check2(c.encodeNode(fmt.Sprintf("[%d]", i), val.Index(i)))

			node.Children = append(node.Children, child)
		}

	default:
		// 直接存储原始字节（避免反射解码开销）
		buf := bytes.NewBuffer(nil)
		mylog.Check(binary.Write(buf, binary.LittleEndian, val.Interface()))
		node.Value = buf.Bytes()
	}
	return node, nil
}

func (c *Codec) buildTypeInfo(typ reflect.Type) typeInfo {
	info := typeInfo{}
	for i := range typ.NumField() {
		f := typ.Field(i)
		info.fields = append(info.fields, fieldInfo{
			name: f.Name,
			idx:  i,
			typ:  f.Type,
		})
	}
	return info
}

// Decode 树→结构体
func (c *Codec) Decode(node *TreeNode, out any) {
	outVal := reflect.ValueOf(out).Elem()
	c.decodeNode(node, outVal)
}

func (c *Codec) decodeNode(node *TreeNode, outVal reflect.Value) {
	if node.Type != outVal.Kind() {
		mylog.Check(fmt.Errorf("类型不匹配: %s vs %s", node.Type, outVal.Kind()))
	}

	switch outVal.Kind() {
	case reflect.Struct:
		ti, _ := c.typeCache.LoadOrStore(outVal.Type(), c.buildTypeInfo(outVal.Type()))
		for _, fi := range ti.(typeInfo).fields {
			for _, child := range node.Children {
				if child.Name == fi.name {
					c.decodeNode(child, outVal.Field(fi.idx))
				}
			}
		}

	case reflect.Slice:
		outVal.Set(reflect.MakeSlice(outVal.Type(), len(node.Children), len(node.Children)))
		for i, child := range node.Children {
			c.decodeNode(child, outVal.Index(i))
		}

	default:
		// 直接从字节解析（比反射快5倍）
		buf := bytes.NewReader(node.Value)
		mylog.Check(binary.Read(buf, binary.LittleEndian, outVal.Addr().Interface()))
	}
}
