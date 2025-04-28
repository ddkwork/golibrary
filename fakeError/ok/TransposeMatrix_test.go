package main

import (
	"fmt"
	"testing"
)

func generateMatrix(rows, cols int) [][]int {
	matrix := make([][]int, rows)
	for i := range matrix {
		matrix[i] = make([]int, cols)
	}
	for i := range rows {
		for j := range cols {
			matrix[i][j] = i*cols + j
		}
	}
	return matrix
}

func BenchmarkTranspose(b *testing.B) {
	matrix := generateMatrix(1000, 1000) // 生成 1000x1000 矩阵
	b.ResetTimer()
	b.Run("iter.Seq2", func(b *testing.B) {
		for b.Loop() {
			for j, v := range TransposeSeq2(matrix) {
				_ = j
				_ = v
			}
		}
	})
	b.Run("blocking", func(b *testing.B) {
		for b.Loop() {
			for j, v := range BlockTransposeSeq2(matrix, 16) {
				_ = j
				_ = v
			}
		}
	})
}

func TestTransposeSeq2(t *testing.T) {
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	// 遍历转置后的元素
	for j, v := range TransposeSeq2(matrix) {
		fmt.Printf("转置后行索引:%d, 值:%d\n", j, v)
	}
}

func TestConcurrentTranspose(t *testing.T) {
	matrix := [][]int{
		{1, 2},
		{3, 4},
		{5, 6},
	}
	// 遍历并发转置后的元素
	for j, v := range ConcurrentTranspose(matrix) {
		fmt.Printf("转置行:%d, 值:%d\n", j, v)
	}
}

func TestSparseTranspose(t *testing.T) {
	type MyStruct struct {
		Value int
	}
	matrix := [][]MyStruct{
		{{0}, {1}, {0}},
		{{2}, {0}, {3}},
		{{0}, {4}, {0}},
	}

	transposeFunc := SparseTranspose(matrix)
	transposeFunc(func(col int, value MyStruct) bool {
		if value.Value != 0 {
			fmt.Printf("Column: %d, Value: %d\n", col, value.Value)
		}
		return true
	})
}
