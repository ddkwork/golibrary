package golibrary
/*

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

// ================= 泛型约束定义 =================

// BinarySupported 定义binary包支持的类型
type BinarySupported interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 | uintptr |
		float32 | float64 |
		bool |
		~[4]byte | ~[8]byte | ~[16]byte | ~[32]byte // 支持固定大小数组
}

// BinaryField 封装结构体字段，提供编译时类型检查
type BinaryField[T BinarySupported] struct {
	value T
}

func (f *BinaryField[T]) Set(value T) {
	f.value = value
}

func (f *BinaryField[T]) Get() T {
	return f.value
}

func (f *BinaryField[T]) Size() int {
	var v T
	return binary.Size(v)
}

// ================= 泛型编解码器 =================

type BinaryCodec struct {
	order binary.ByteOrder
}

func NewBinaryCodec(order binary.ByteOrder) *BinaryCodec {
	return &BinaryCodec{order: order}
}

// EncodeStruct 编码整个结构体
func (c *BinaryCodec) EncodeStruct(data interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, c.order, data); err != nil {
		return nil, fmt.Errorf("编码失败: %w", err)
	}
	return buf.Bytes(), nil
}

// DecodeStruct 解码整个结构体
func (c *BinaryCodec) DecodeStruct(data []byte, target interface{}) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, c.order, target)
}

// EncodeField 编码单个字段（泛型方法）
func (c *BinaryCodec) EncodeField[T BinarySupported](field *BinaryField[T]) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, c.order, field.Get()); err != nil {
		return nil, fmt.Errorf("字段编码失败: %w", err)
	}
	return buf.Bytes(), nil
}

// DecodeField 解码单个字段（泛型方法）
func (c *BinaryCodec) DecodeField[T BinarySupported](data []byte, field *BinaryField[T]) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, c.order, &field.value)
}

// ================= 使用示例 =================

// 定义测试结构体使用编译时安全字段
type NetworkPacket struct {
	ID     BinaryField[uint32]
	Length BinaryField[uint16]
	Status BinaryField[uint8]  // 使用uint8表示布尔状态：1=true, 0=false
	Header BinaryField[[16]byte]
	// Data   BinaryField[[]byte] // 编译错误：[]byte不满足约束（不能编译通过）
}

func main() {
	// 实例化结构体
	packet := &NetworkPacket{}
	packet.ID.Set(12345)
	packet.Length.Set(1500)
	packet.Status.Set(1) // true
	packet.Header.Set([16]byte{'H', 'e', 'a', 'd', 'e', 'r', 'D', 'a', 't', 'a'})

	// 创建编解码器
	codec := NewBinaryCodec(binary.LittleEndian)

	// 编码整个结构体
	encoded, err := codec.EncodeStruct(packet)
	if err != nil {
		panic(err)
	}
	fmt.Printf("编码后的字节: % x\n", encoded)
	fmt.Printf("总大小: %d bytes\n\n", len(encoded))

	// 单个字段编码示例
	idData, err := codec.EncodeField(&packet.ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ID字段编码: % x\n", idData)

	// 解码整个结构体
	newPacket := &NetworkPacket{}
	if err := codec.DecodeStruct(encoded, newPacket); err != nil {
		panic(err)
	}

	fmt.Println("\n解码结果:")
	fmt.Printf("ID: %d\n", newPacket.ID.Get())
	fmt.Printf("Length: %d\n", newPacket.Length.Get())
	fmt.Printf("Status: %t\n", newPacket.Status.Get() == 1)
	fmt.Printf("Header: %s\n", string(newPacket.Header.Get()[:5])) // 只打印前5字节

	// ============ 非法类型演示 ============
	// 尝试创建包含不支持类型的字段
	type InvalidStruct struct {
		Slice BinaryField[[]byte] // 编译错误：[]byte不满足约束
	}

	// 直接尝试使用非法类型实例化 - 注释掉才能编译通过
	// var invalid BinaryField[[]byte]
}

*/