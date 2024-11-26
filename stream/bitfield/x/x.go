package main

import (
	"fmt"
)

type BitField struct {
	data [16]byte // 支持到16字节
}

// 从给定的无符号整数设置从偏移位置开始的指定字节数
func (b *BitField) SetBytesFromUint(offset int, byteCount int, value uint64) {
	if offset+byteCount > 128 {
		panic("Field exceeds the maximum size of 128 bits")
	}

	// 逐字节设置
	for i := 0; i < byteCount; i++ {
		b.data[offset+i] = uint8(value >> (8 * (byteCount - 1 - i))) // 从高字节到低字节
	}
}

// 从偏移位置获取指定字节数的无符号整数值
func (b *BitField) GetUintFromBytes(offset int, byteCount int) uint64 {
	if offset+byteCount > 128 {
		panic("Field exceeds the maximum size of 128 bits")
	}
	var result uint64 = 0

	// 逐字节获取
	for i := 0; i < byteCount; i++ {
		result |= uint64(b.data[offset+i]) << (8 * (byteCount - 1 - i)) // 从高字节到低字节
	}
	return result
}

func main() {
	var bf BitField

	// 设置三个字节的值为 23
	bf.SetBytesFromUint(10, 3, 0x000017) // 第10字节开始，填充3个字节

	// 获取并输出字段值
	field1 := bf.GetUintFromBytes(10, 3) // 获取3个字节的值

	fmt.Printf("Field1: %b\n", field1)     // 输出二进制
	fmt.Printf("Field1: %d\n", field1)     // 输出十进制
	fmt.Printf("Field1: 0x%06X\n", field1) // 输出十六进制
}
