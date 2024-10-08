package align

import (
	"fmt"
	"gioui.org/unit"
	"strings"
	"unicode"
)

func isChinese(r rune) bool {
	return unicode.Is(unicode.Han, r)
}

func StringWidth(s string) (width unit.Dp) {
	for _, char := range s {
		if isChinese(char) {
			width += 2 // 中文字符宽度
		} else {
			width += 1 // 英文和其他字符宽度
		}
	}
	return
}

func formatString(s string) string {
	const minWidth = 50 // 总宽度可以根据需要调整，设小了不行
	currentWidth := StringWidth(s)
	padding := minWidth - currentWidth
	if padding > 0 {
		return strings.Repeat(".", int(padding)) + s
		return s + strings.Repeat(" ", int(padding))
	}
	return s
}

func main() {
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
		//fmt.Printf("name：%s \t scores:%d\n", formatString(name), score)
		fmt.Printf("姓名：%s \t 分数:%d\n", formatString(name), score) //因为特殊字符的视觉宽度真的就是以像素为单位的，加一个\t
	}
}

//PS: 对齐要求
//1. minWidth =50
//2.对于goland需要consolas.ttf，设置--》editor-->font-->consolas
//3. fmt.Printf(" 姓名：%s \t 分数:%d\n", formattedName, score) //因为特殊字符的视觉宽度真的就是以像素为单位的，加一个\t
//4. gio代码编辑器字体设置似乎没生效，得用consolas.ttf，logview对不齐也是字体没生效的原因
//5. todo merge into mylog package
