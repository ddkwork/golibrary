package main

import (
	"fmt"
	"github.com/ddkwork/golibrary/mylog"
	"iter"
	"reflect"
	"strconv"
	"time"
)

//Table 是row布局，grid是rows和columns布局，网格是处理所有单元格的，row的话是通过结构体处理行单元格的
//column布局只有一个场景，复制列数据到剪切板，没有别的场景了。不需要支持

//布局说明:
//structView:key-value形式的结构体字段渲染成一行，然后flex垂直排列
//树形表格: key-value形式的结构体字段渲染成一行，然后flex水平排列，key用于每一列的标题，value用于每一行的每个单元格值
//两种场景都支持自定义格式化和反序列化字段，通过回调函数实现，
//对于序列化，回调返回空则使用默认的fmt包格式化字段值，否则使用回调返回值
//对于反序列化，回调返回nil则使用默认的反序列化方式，否则使用回调返回值更新字段值

// filed或者cell都是key-value键值对，对于form和structView，key是字段名，value是字段值，对于树形表格，key是列名，value是单元格值
// 树形节点编辑调用structView布局
// 适用场景：树形表格，structView，form 布局
type codec[T any] interface {
	MarshalRow(callback func(key string) (value string)) (cells iter.Seq2[string, string])
	UnmarshalRow(cells iter.Seq2[string, string], callback func(key string) (value any)) T
	MarshalFields(callback func(key string) (value string)) iter.Seq2[string, string]
	UnmarshalFields(fields iter.Seq2[string, string], callback func(key string) (value any)) T
}

var _ codec[any] = (*Struct[any])(nil)

type Struct[T any] struct {
	Data T
}

func (s *Struct[T]) MarshalRow(callback func(key string) (value string)) (cells iter.Seq2[string, string]) { //flex水平排列
	return s.MarshalFields(callback) //获得字段键值对后渲染树形表格的row是flex水平排列，节点编辑是structView的flex垂直排列
}
func (s *Struct[T]) UnmarshalRow(cells iter.Seq2[string, string], callback func(key string) (value any)) T { //树形节点编辑,flex垂直排列,布局成StructView
	return s.UnmarshalFields(cells, callback) //从StructView的每一行的键值对反序列化成结构体,结果是更新树形表格的rowCells
}

// MarshalFields 序列化字段,想象成structView的flex垂直布局的每一行，每行都是一个键值对，key是字段名，value是字段值,callback理解成structView的每一行的自定义序列化的callback
func (s *Struct[T]) MarshalFields(callback func(key string) (value string)) iter.Seq2[string, string] { //flex垂直排列
	return func(yield func(string, string) bool) {
		fields := reflect.VisibleFields(reflect.TypeOf(s.Data).Elem())
		value := reflect.ValueOf(s.Data).Elem()
		if len(fields) != value.NumField() {
			panic("wrong number of fields")
		}
		for i, field := range fields {
			if field.Tag.Get("table") == "-" || field.Tag.Get("json") == "-" {
				mylog.Trace("field is ignored: ", field.Name) // 用于树形表格序列化json保存到文件，标签table或json为-则忽略
				continue                                      //todo test
			}
			if !field.IsExported() {
				mylog.Trace("field name is not exported: ", field.Name) // 用于树形表格序列化json保存到文件，没有导出则json会失败
				continue
			}
			v := value.Field(i).Interface()
			if callback != nil {
				v2 := callback(field.Name)
				if v2 != "" {
					v = v2
				}
			}
			if !yield(field.Name, fmt.Sprint(v)) {
				return
			}
		}
	}
}

// UnmarshalFields 反序列化字段,fields理解成struckView或者form布局的所有行，每行都是一个键值对，key是字段名，value是字段值,callback理解成structView的每一行的自定义反序列化的callback
func (s *Struct[T]) UnmarshalFields(fields iter.Seq2[string, string], callback func(key string) (value any)) T {
//todo
//	decoded := reflect.New(reflect.TypeOf(s.Data).Elem()).Interface()

	valueOf := reflect.ValueOf(s.Data)

	// 关键修复点1：确保是指针类型且可寻址
	if valueOf.Kind() != reflect.Ptr || valueOf.Elem().Kind() != reflect.Struct {
		panic("Data must be a pointer to struct")
	}
	valueOf = valueOf.Elem()

	// 遍历输入字段
	for k, v := range fields {
		// 关键修复点2：正确获取字段反射对象
		field, ok := valueOf.Type().FieldByName(k)
		if !ok {
			continue // 忽略不存在的字段
		}

		// 关键修复点3：检查字段导出性和标签
		if !field.IsExported() ||
			field.Tag.Get("table") == "-" ||
			field.Tag.Get("json") == "-" {
			continue
		}

		fieldValue := valueOf.FieldByName(k)
		if !fieldValue.CanSet() {
			panic("field is not settable: " + k)
			continue
		}

		// 处理回调覆盖
		if callback != nil {
			if val := callback(k); val != nil {
				if rv := reflect.ValueOf(val); rv.Type().AssignableTo(fieldValue.Type()) {
					fieldValue.Set(rv)
					continue
				}
			}
		}

		// 如果回调返回 nil，则使用默认的反序列化方式
		// 根据字段类型将字符串转换为相应的类型
		//一般情况下，树形表格的元数据结构体字段类型不会太复杂，它也不会引入嵌套字段，因为层级是通过线性布局的字段构造的，所以这里只需要处理简单类型即可
		switch field.Type.Kind() { //todo 单元测试时间相关的类型，以及实现了stringer接口的字段类型
		case reflect.String:
			fieldValue.SetString(v)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			//fieldValue.SetInt(mylog.Check2(strconv.ParseInt(v, 10, 64))) //如果是其他进制和位数呢？这就需要回调函数了
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fieldValue.SetUint(mylog.Check2(strconv.ParseUint(v, 10, 64)))
		case reflect.Float32, reflect.Float64:
			fieldValue.SetFloat(mylog.Check2(strconv.ParseFloat(v, 64)))
		case reflect.Bool:
			fieldValue.SetBool(mylog.Check2(strconv.ParseBool(v)))
		case reflect.TypeFor[time.Duration]().Kind():
			fieldValue.Set(reflect.ValueOf(mylog.Check2(time.ParseDuration(v))))
		case reflect.TypeFor[time.Time]().Kind():
			fieldValue.Set(reflect.ValueOf(mylog.Check2(time.Parse(time.RFC3339, v))))
		default:
			//todo n.Data.PadTime = mylog.Check2(time.ParseDuration(value))
			// 此外，那些类型别名，结构体别名的玩意,需要足够的单元测试和模糊测试,性能测试
			mylog.Check("unsupported field type: " + field.Type.Kind().String()) //回调没有处理不支持字段类型，这里有没有case到，直接panic了
		}
	}
	//老外的做法:和返回一个new(T)没区别吧？
	//	reflect.ValueOf(&e.beforeData).Elem().Set(reflect.New(reflect.TypeOf(e.beforeData).Elem()))
	//	e.beforeData.CopyFrom(target)
	//
	//	reflect.ValueOf(&e.editorData).Elem().Set(reflect.New(reflect.TypeOf(e.editorData).Elem()))
	//	e.editorData.CopyFrom(target)
	return s.Data
}

//todo 重命名树形表格 ctx 自定义上下文菜单 CustomContextMenuItems
//   求和也使用反射自动完成，不要在每个表格里面都写相同的逻辑
//   结构体字段不应该重名，现在使用迭代器虽然节省内存，但是理论上应该使用有序map来出重，最后就是 tag - 和非导出字段等等一系列玩意一个封装一个valid函数统一检查
