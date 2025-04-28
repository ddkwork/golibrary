package main

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"
)

type UserForm struct {
	Name    string    `form:"name"`
	Age     int       `form:"age"`
	Salary  float64   `form:"salary"`
	Created time.Time `form:"created"`
}

func main() {
	user := UserForm{
		Name:    "张三",
		Age:     30,
		Salary:  15000.50,
		Created: time.Now(),
	}

	// 定义自定义格式化规则
	fieldFmtCallback := map[string]FieldFmtCallback{
		"Salary": func(v any) string {
			return fmt.Sprintf("￥%.0f元", v.(float64))
		},
	}

	// 编码到map
	m := EncodeToForm(&user, "form", fieldFmtCallback)
	fmt.Printf("编码结果: %#v\n", m)

	// 解码到结构体
	newUser := &UserForm{}
	newFieldFmtCallback := map[string]func(string) (any, error){
		"Salary": func(s string) (any, error) {
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return nil, err
			}
			return f, nil
		},
	}
	DecodeFromForm(map[string]string{
		"name":    "李四",
		"age":     "25",
		"salary":  "20000.5",
		"created": "2023-04-22",
	}, newUser, "form", newFieldFmtCallback)
	fmt.Printf("解码结果: %+v", newUser)
}

func TestEncodeDecode(t *testing.T) {
	testCases := []struct {
		name             string
		input            any
		form             map[string]string
		fieldFmtCallback map[string]FieldFmtCallback
		expected         map[string]string
		decodeCallback   map[string]func(string) (any, error)
	}{
		{
			name: "基本类型转换",
			input: &struct {
				Name string `form:"name"`
				Age  int    `form:"age"`
			}{
				Name: "test",
				Age:  20,
			},
			expected: map[string]string{
				"name": "test",
				"age":  "20",
			},
		},
		{
			name: "时间类型格式化",
			input: &struct {
				Date time.Time `form:"date"`
			}{
				Date: time.Date(2023, 4, 22, 0, 0, 0, 0, time.UTC),
			},
			fieldFmtCallback: map[string]FieldFmtCallback{
				"Date": func(v any) string {
					return v.(time.Time).Format("2006/01/02")
				},
			},
			expected: map[string]string{
				"date": "2023/04/22",
			},
			decodeCallback: map[string]func(string) (any, error){
				"Date": func(s string) (any, error) {
					t, err := time.Parse("2006/01/02", s)
					if err != nil {
						return nil, err
					}
					return t, nil
				},
			},
		},
		{
			name: "自定义格式化和解码",
			input: &struct {
				Salary float64 `form:"salary"`
			}{
				Salary: 15000.50,
			},
			fieldFmtCallback: map[string]FieldFmtCallback{
				"Salary": func(v any) string {
					return fmt.Sprintf("%.2f", v.(float64))
				},
			},
			expected: map[string]string{
				"salary": "15000.50",
			},
			decodeCallback: map[string]func(string) (any, error){
				"Salary": func(s string) (any, error) {
					f, err := strconv.ParseFloat(s, 64)
					if err != nil {
						return nil, err
					}
					return f, nil
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 测试编码
			encoded := EncodeToForm(tc.input, "form", tc.fieldFmtCallback)
			if !reflect.DeepEqual(encoded, tc.expected) {
				t.Errorf("编码失败: 期望 %v，实际 %v", tc.expected, encoded)
			}

			// 测试解码
			decoded := reflect.New(reflect.TypeOf(tc.input).Elem()).Interface()
			DecodeFromForm(encoded, decoded, "form", tc.decodeCallback)
			if !reflect.DeepEqual(decoded, tc.input) {
				t.Errorf("解码失败: 输入 %v，输出 %v", tc.input, decoded)
			}
		})
	}
}
