package main

import (
	"iter"
	"reflect"
	"sync"
)

func TransposeMatrix[T any](rootRows [][]T) iter.Seq2[int, iter.Seq[T]] {
	return func(yield func(int, iter.Seq[T]) bool) {
		if len(rootRows) == 0 {
			return
		}
		numCols := len(rootRows[0])

		// 校验所有行长度一致（关键步骤）
		for _, row := range rootRows[1:] {
			if len(row) != numCols {
				panic("irregular matrix shape") // 或返回错误
			}
		}

		// 生成列迭代器
		for j := range numCols {
			colSeq := func(yieldElem func(T) bool) {
				for i := range rootRows {
					if !yieldElem(rootRows[i][j]) {
						return
					}
				}
			}
			if !yield(j, colSeq) {
				return
			}
		}
	}
}

// 泛型矩阵转置迭代器
func TransposeSeq2[T any](rootRows [][]T) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		rows := len(rootRows)
		if rows == 0 {
			return
		}
		cols := len(rootRows[0])
		// 遍历原矩阵的列作为转置后的行
		for j := range cols {
			// 遍历原矩阵的行作为转置后的列
			for i := range rows {
				// 传递转置后的行索引 j 和元素值 rootRows[i][j]
				if !yield(j, rootRows[i][j]) {
					return // 若 yield 返回 false，提前终止迭代
				}
			}
		}
	}
}

func ConcurrentTranspose[T any](rootRows [][]T) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		next, stop := iter.Pull2(TransposeSeq2(rootRows))
		defer stop()
		for {
			j, v, ok := next()
			if !ok {
				return
			}
			if !yield(j, v) {
				return
			}
		}
	}
}

func SparseTranspose[T any](rootRows [][]T) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for _, row := range rootRows {
			for j, v := range row {
				if !isZero(v) { // 自定义零值判断，稀疏矩阵的零值判断模拟 //https://github.com/BaseMax/SparseMatrixLinkedListGo
					if !yield(j, v) {
						return
					}
				}
			}
		}
	}
}

// isZero 使用反射来判断一个值是否为零值
func isZero[T any](v T) bool {
	val := reflect.ValueOf(v)
	zero := reflect.Zero(val.Type())
	return val.Interface() == zero.Interface()
}

// Zeroable 是一个接口，要求实现 IsZero 方法来判断零值
type Zeroable interface {
	IsZero() bool
}

// SparseTransposeCustom 是一个泛型函数，接受一个二维切片作为输入，并返回一个生成器函数
func SparseTransposeCustom[T Zeroable](rootRows [][]T) func(yield func(int, T) bool) {
	return func(yield func(int, T) bool) {
		for _, row := range rootRows {
			for j, v := range row {
				if !v.IsZero() { // 使用自定义的零值判断
					if !yield(j, v) {
						return
					}
				}
			}
		}
	}
}

func TransposeSeq3[T any](rootRows [][]T) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for _, row := range rootRows {
			for j, v := range row {
				if !yield(j, v) { // j 是列索引，作为转置后的行键
					return
				}
			}
		}
	}
}

func BlockTransposeSeq2[T any](rootRows [][]T, blockSize int) iter.Seq2[int, T] { //分块提高命中率，理论效果是减少耗时
	return func(yield func(int, T) bool) {
		rows := len(rootRows)
		if rows == 0 {
			return
		}
		cols := len(rootRows[0])
		for jb := 0; jb < cols; jb += blockSize { // 按块遍历
			for ib := 0; ib < rows; ib += blockSize {
				for j := jb; j < min(jb+blockSize, cols); j++ { // 处理当前块内的转置
					for i := ib; i < min(ib+blockSize, rows); i++ {
						if !yield(j, rootRows[i][j]) {
							return
						}
					}
				}
			}
		}
	}
}

func FlatTranspose[T any](rootRows []T, rows, cols int) []T {
	result := make([]T, rows*cols)
	for i := range rows {
		for j := range cols {
			result[j*rows+i] = rootRows[i*cols+j]
		}
	}
	return result
}
func FlatTranspose9[T any](data []T, rows, cols int) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for j := range cols {
			for i := range rows {
				if !yield(j, data[i*cols+j]) {
					return
				} // 原矩阵按行存储
			}
		}
	}
}

func ConcurrentBlockTranspose[T any](rootRows [][]T) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		var wg sync.WaitGroup
		blockSize := 16
		rows := len(rootRows)
		cols := len(rootRows[0])
		for jb := 0; jb < cols; jb += blockSize {
			wg.Add(1)
			go func(jStart int) {
				defer wg.Done()
				for j := jStart; j < min(jStart+blockSize, cols); j++ {
					for i := range rows {
						if !yield(j, rootRows[i][j]) {
							return
						}
					}
				}
			}(jb)
		}
		wg.Wait()
	}
}
