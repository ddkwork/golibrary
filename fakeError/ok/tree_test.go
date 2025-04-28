package main

import (
	"fmt"
	"testing"
)

func BenchmarkCodec(b *testing.B) {
	type Data struct {
		ID    int
		Names []string
		Coord struct{ X, Y float64 }
	}

	codec := NewCodec()
	data := Data{
		ID:    1001,
		Names: []string{"A", "B", "C"},
		Coord: struct{ X, Y float64 }{X: 1.1, Y: 2.2},
	}

	b.Run("Encode", func(b *testing.B) {
		for b.Loop() {
			codec.Encode(&data)
		}
	})

	b.Run("Decode", func(b *testing.B) {
		tree, _ := codec.Encode(&data)
		var out Data
		for b.Loop() {
			codec.Decode(tree, &out)
		}
	})
}

func TestTree(t *testing.T) {
	user := User{
		ID:   1001,
		Name: "Alice",
		Address: struct {
			City    string
			ZipCode string
		}{City: "Beijing", ZipCode: "100000"},
		Devices: []string{"Phone", "Laptop"},
	}

	// 编码为树形结构
	tree, _ := StructToTree(&user)
	fmt.Printf("Tree Nodes: %d\n", countNodes(tree)) // 输出节点总数

	// 解码回结构体
	var decoded User
	_ = TreeToStruct(tree, &decoded)
	fmt.Printf("Decoded User: %+v\n", decoded)
}

func countNodes(node *TreeNode) int {
	count := 1
	for _, child := range node.Children {
		count += countNodes(child)
	}
	return count
}

type User struct {
	ID      int
	Name    string
	Address struct {
		City    string
		ZipCode string
	}
	Devices []string
}
