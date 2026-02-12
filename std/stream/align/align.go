package align

import (
	"strings"
	"unicode"

	"github.com/ddkwork/golibrary/types"
)

// StringWidth PS: 对齐要求
// 1. minWidth =50
// 2.对于goland需要consolas.ttf，设置--》editor-->font-->consolas
// 3. fmt.Printf(" 姓名：%s \t 分数:%d\n", formattedName, score) //因为特殊字符的视觉宽度真的就是以像素为单位的，加一个\t
// 4. gio代码编辑器字体设置似乎没生效，得用consolas.ttf，logview对不齐也是字体没生效的原因
func StringWidth[T types.Number](s string) (width T) {
	for _, char := range s {
		if isChinese(char) {
			width += 2 // 中文字符宽度
		} else {
			width += 1 // 英文和其他字符宽度
		}
	}
	return
}

func isChinese(r rune) bool {
	return unicode.Is(unicode.Han, r)
}

func formatString(s string) string {
	const minWidth = 50 // 总宽度可以根据需要调整，设小了不行
	currentWidth := StringWidth[int](s)
	padding := minWidth - currentWidth
	if padding > 0 {
		return strings.Repeat(".", padding) + s
		return s + strings.Repeat(" ", padding)
	}
	return s
}
