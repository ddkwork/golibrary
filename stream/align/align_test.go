package align

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	data := map[string]int{
		"张三":           95,
		"李四":           85,
		"John":         90,
		"Jane Doe":     88,
		"Jane Doe 李四":  88,
		"0xb76 也可以 12": 88,
		"张三tt":         85,
		"李四qq":         90,
		"firstEnd xor 0x72B8 删除第一个字节": 75,
		"firstEnd xor  删除第一个字节":       75,
	}
	for name, score := range data {
		// fmt.Printf("name：%s \t scores:%d\n", formatString(name), score)
		fmt.Printf("姓名：%s \t 分数:%d\n", formatString(name), score) // 因为特殊字符的视觉宽度真的就是以像素为单位的，加一个\t
	}
}
