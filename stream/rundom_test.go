package stream

import (
	"fmt"
	"strconv"
	"testing"
)

func TestAlignString(t *testing.T) {
	fmt.Println(strconv.Quote(AlignString("中文SetHan═╬═dles(ha电═╬═饭锅电饭锅ndles []Handle)", 55)))
	fmt.Println(strconv.Quote(AlignString("Handlesjk═╬═js 看见你地方df() []Handf的 dle", 55)))
	fmt.Println(strconv.Quote(AlignString("en═╬═flish", 55)))
}

func TestIsDirDeep1(t *testing.T) {
	println(IsDirDeep1("pkg\\cpp2go\\cpp"))
	println(IsDirDeep1(".git"))
}
