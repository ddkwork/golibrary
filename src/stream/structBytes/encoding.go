package structBytes // import "go.dedis.ch/fixbuf"

import (
	"fmt"
	"strings"
)

//反射精华：value包含了type和value，type就只包含type，value就是结构体对象

//思路：保存结构体成员的数据结构可以定义为:妈的，终于有了思路，填充了这个map之后可以直接用这个google包编解码
//D:\CodeBackup\good\HardInfo\protoreflect 这个命名清晰一点
//泛型是传递类型，反射也是传递类型，反射只要是性能和编译期间检查不完善，泛型的话不用case类型了，不容易编写和阅读，但事实可以减少代码量
//这方面反射的可读性感觉也强不到哪里去，不理解反射就算代码写正确了还是看得一脸懵逼
//func sd()  {
//	type  structFalid struct {
//		FalidType reflect.Type
//		FalidValue reflect.value
//		structCTX reflect.StructField
//	}
//	a:= map[reflect.Type]structFalid{}
//}

//以这个结构体来理解发射
//type obj struct { obj就是value
//	a int  ----》value的左边是type int，右边是type int +value:obj.a=9999，就是时候一个数据首先有了类型，然后赋值
//}
//map结构为什么适合保存成员，还要研究下别人的代码,可能是value部分，就是右边的值部分的字段有index，所以可以排序或者顺序操作。。。

//终于把反射理解进去了一点，反思了一下感觉泛型不适合直接把结构体编解码，因为pb从结构体编解码要case对应的字段类型执行对应类型的编解码函数，
//泛型的话，实现的形参就写死了类型，那么传递结构体的时候需要把每个字段对应的类型编解码函数传递，根本不能一次性传递一个完整的填充好成员数据的结构体
//总的来说泛型避免了业务实现case类型却不能一次性操作所有类型:interface

func info(depth int, format string, a ...any) {
	fmt.Print(strings.Repeat("  ", depth))
	fmt.Printf(format, a...)
}
