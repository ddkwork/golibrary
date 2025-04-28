package main

import (
	"fmt"
	"github.com/ddkwork/golibrary/mylog"
	"reflect"
	"strconv"
	"time"
)

type FieldFmtCallback func(any) string

func EncodeToForm(data any, tagName string, fieldFmtCallback map[string]FieldFmtCallback) map[string]string {
	fields := make(map[string]string)
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		mylog.Check("input must be struct")
	}

	t := v.Type()
	for i := range v.NumField() {
		field := t.Field(i)
		tagVal := field.Tag.Get(tagName)
		if tagVal == "" || tagVal == "-" {
			continue
		}

		fieldVal := v.Field(i).Interface()

		// 应用自定义格式化
		if formatter, ok := fieldFmtCallback[field.Name]; ok {
			fields[tagVal] = formatter(fieldVal)
		} else {
			// 默认格式化逻辑
			switch val := fieldVal.(type) {
			case int, int32, int64:
				fields[tagVal] = fmt.Sprintf("%d", val)
			case float32, float64:
				fields[tagVal] = fmt.Sprintf("%.2f", val)
			case time.Time:
				fields[tagVal] = val.Format("2006-01-02")
			default:
				fields[tagVal] = fmt.Sprintf("%v", val)
			}
		}
	}
	return fields
}
func DecodeFromForm(fields map[string]string, dataPtr any, tagName string, fieldFmtCallback map[string]func(string) (any, error)) {
	v := reflect.ValueOf(dataPtr)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		mylog.Check("input must be struct pointer")
	}
	v = v.Elem()

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tagVal := field.Tag.Get(tagName)
		if tagVal == "" || tagVal == "-" {
			continue
		}

		strVal, exists := fields[tagVal]
		if !exists {
			continue
		}

		// 应用自定义格式化
		if formatter, ok := fieldFmtCallback[field.Name]; ok {
			val, err := formatter(strVal)
			if err != nil {
				mylog.Check(err)
				continue
			}
			switch v.Field(i).Kind() {
			case reflect.String:
				v.Field(i).SetString(val.(string))
			case reflect.Int, reflect.Int32, reflect.Int64:
				intVal := val.(int64)
				v.Field(i).SetInt(intVal)
			case reflect.Float32, reflect.Float64:
				floatVal := val.(float64)
				v.Field(i).SetFloat(floatVal)
			case reflect.Struct:
				if field.Type == reflect.TypeOf(time.Time{}) {
					t, err := time.Parse("2006-01-02", val.(time.Time).String())
					if err != nil {
						mylog.Check(err)
						continue
					}
					v.Field(i).Set(reflect.ValueOf(t))
				}
			default:
				mylog.Check(fmt.Errorf("unsupported type: %s", v.Field(i).Kind()))
			}
		} else {
			// 默认格式化逻辑
			switch fieldType := field.Type; fieldType.Kind() {
			case reflect.String:
				v.Field(i).SetString(strVal)
			case reflect.Int, reflect.Int32, reflect.Int64:
				intVal := mylog.Check2(strconv.ParseInt(strVal, 10, 64))
				v.Field(i).SetInt(intVal)
			case reflect.Float32, reflect.Float64:
				floatVal := mylog.Check2(strconv.ParseFloat(strVal, 64))
				v.Field(i).SetFloat(floatVal)
			case reflect.Struct:
				if fieldType == reflect.TypeOf(time.Time{}) {
					t := mylog.Check2(time.Parse("2006-01-02", strVal))
					v.Field(i).Set(reflect.ValueOf(t))
				}
			default:
				mylog.Check(fmt.Errorf("unsupported type: %s", fieldType.Kind()))
			}
		}
	}
}

func DecodeFromForm1(fields map[string]string, dataPtr any, tagName string) {
	v := reflect.ValueOf(dataPtr)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		mylog.Check("input must be struct pointer")
	}
	v = v.Elem()

	t := v.Type()
	for i := range v.NumField() {
		field := t.Field(i)
		tagVal := field.Tag.Get(tagName)
		if tagVal == "" || tagVal == "-" {
			continue
		}

		strVal, exists := fields[tagVal]
		if !exists {
			continue
		}

		fieldType := field.Type
		fieldValue := v.Field(i)

		// 类型转换逻辑
		switch fieldType.Kind() {
		case reflect.String:
			fieldValue.SetString(strVal)
		case reflect.Int, reflect.Int32, reflect.Int64:
			intVal := mylog.Check2(strconv.ParseInt(strVal, 10, 64))
			fieldValue.SetInt(intVal)
		case reflect.Float32, reflect.Float64:
			floatVal := mylog.Check2(strconv.ParseFloat(strVal, 64))
			fieldValue.SetFloat(floatVal)
		case reflect.Struct:
			if fieldType == reflect.TypeOf(time.Time{}) {
				t := mylog.Check2(time.Parse("2006-01-02", strVal))
				fieldValue.Set(reflect.ValueOf(t))
			}
		}
	}
}
