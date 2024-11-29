package pretty

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	r "reflect"
	"strconv"
	"strings"
)

const (
	DEFAULT_INDENT = "  "
	DEFAULT_NIL    = "nil"
)

type Pretty struct {
	// indent string
	Indent string
	// output recipient
	Out io.Writer
	// string for nil
	NilString string
	// compact empty array and struct
	Compact bool
	// Maximum nesting level
	MaxLevel int
}

func Print(i interface{}) {
	PrintTo(os.Stdout, i, true)
}

func Format(i interface{}) string {
	var out bytes.Buffer
	PrintTo(&out, i, false)
	return out.String()
}

func PrintTo(out io.Writer, i interface{}, nl bool) {
	p := &Pretty{Indent: DEFAULT_INDENT, Out: out, NilString: DEFAULT_NIL}
	if nl {
		p.Println(i)
	} else {
		p.Print(i)
	}
}

func (p *Pretty) Print(i interface{}) {
	p.PrintValue(r.ValueOf(i), 0)
}

func (p *Pretty) Println(i interface{}) {
	p.PrintValue(r.ValueOf(i), 0)
	io.WriteString(p.Out, "\n")
}

func (p *Pretty) PrintValue(val r.Value, level int) {
	if !val.IsValid() {
		io.WriteString(p.Out, p.NilString)
		return
	}

	cur := strings.Repeat(p.Indent, level)
	next := strings.Repeat(p.Indent, level+1)

	nl := "\n"
	if len(p.Indent) == 0 {
		nl = " "
	}

	if p.MaxLevel > 0 && level >= p.MaxLevel {
		io.WriteString(p.Out, val.String())
		return
	}

	switch val.Kind() {
	case r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
		io.WriteString(p.Out, "0x")
		switch val.Kind() {
		case r.Int, r.Int64:
			io.WriteString(p.Out, FormatInteger(val.Int()))
		case r.Int8:
			io.WriteString(p.Out, FormatInteger(int8(val.Int())))
		case r.Int16:
			io.WriteString(p.Out, FormatInteger(int16(val.Int())))
		case r.Int32:
			io.WriteString(p.Out, FormatInteger(int32(val.Int())))
		}
		io.WriteString(p.Out, "│")
		io.WriteString(p.Out, strconv.FormatInt(val.Int(), 10))
		m, ok := val.Interface().(fmt.Stringer)
		if ok {
			io.WriteString(p.Out, " "+m.String())
		}
	case r.Uint, r.Uint8, r.Uint16, r.Uint32, r.Uint64, r.Uintptr:
		io.WriteString(p.Out, "0x")
		switch val.Kind() {
		case r.Uint, r.Uint64:
			io.WriteString(p.Out, FormatInteger(val.Uint()))
		case r.Uint8:
			io.WriteString(p.Out, FormatInteger(byte(val.Uint())))
		case r.Uint16:
			io.WriteString(p.Out, FormatInteger(uint16(val.Uint())))
		case r.Uint32:
			io.WriteString(p.Out, FormatInteger(uint32(val.Uint())))
		}
		io.WriteString(p.Out, "│")
		io.WriteString(p.Out, strconv.FormatUint(val.Uint(), 10))
		m, ok := val.Interface().(fmt.Stringer)
		if ok {
			io.WriteString(p.Out, " "+m.String())
		}
	case r.Float32, r.Float64:
		io.WriteString(p.Out, strconv.FormatFloat(val.Float(), 'f', -1, 64))

	case r.String:
		io.WriteString(p.Out, strconv.Quote(val.String()))

	case r.Bool:
		io.WriteString(p.Out, strconv.FormatBool(val.Bool()))

	case r.Map:
		l := val.Len()

		io.WriteString(p.Out, "{"+nl)
		for i, k := range val.MapKeys() {
			io.WriteString(p.Out, next)
			io.WriteString(p.Out, strconv.Quote(k.String()))
			io.WriteString(p.Out, ": ")
			p.PrintValue(val.MapIndex(k), level+1)
			if i < l-1 {
				io.WriteString(p.Out, ","+nl)
			} else {
				io.WriteString(p.Out, nl)
			}
		}
		io.WriteString(p.Out, cur)
		io.WriteString(p.Out, "}")

	case r.Array, r.Slice:
		l := val.Len()
		if ValueIsBytesType(val) {
			if l == 0 {
				io.WriteString(p.Out, "[]bytes{}")
				return
			}
			dump := hex.Dump(val.Bytes())
			dump = strings.TrimPrefix(dump, "\n")
			dump = strings.TrimSuffix(dump, "\n")
			mixHex := "00000000  3c 4d 55 fe 16 14 00 00  00 d8 70 c5 70 e8 ab 30  |<MU.......p.p..0|"
			if len(dump) > len(mixHex) {
				dump = "\n" + dump // todo indent left it

				for i, s := range strings.Split(dump, "\n") {
					if i > 0 {
						io.WriteString(p.Out, "\n")
					}
					io.WriteString(p.Out, cur) // todo 取字段名称长度?
					io.WriteString(p.Out, s)
				}
				return
			}
			io.WriteString(p.Out, fmt.Sprintf("%#v", val.Bytes()))
			io.WriteString(p.Out, " //"+hex.EncodeToString(val.Bytes()))
			// todo get tag and print string
			if i, ok := val.Interface().(fmt.Stringer); ok {
				io.WriteString(p.Out, " //"+strconv.Quote(i.String()))
			}
			return
		}

		if p.Compact && l == 0 {
			io.WriteString(p.Out, "[]")
		} else {
			io.WriteString(p.Out, "["+nl)
			for i := 0; i < l; i++ {

				io.WriteString(p.Out, next)

				io.WriteString(p.Out, "<")
				io.WriteString(p.Out, strconv.Itoa(i))
				io.WriteString(p.Out, "> ")

				p.PrintValue(val.Index(i), level+1)
				if i < l-1 {
					io.WriteString(p.Out, ","+nl)
				} else {
					io.WriteString(p.Out, nl)
				}
			}
			io.WriteString(p.Out, cur)
			io.WriteString(p.Out, "]")
		}

	case r.Interface, r.Ptr:
		p.PrintValue(val.Elem(), level)

	case r.Struct:
		if val.CanInterface() {
			i := val.Interface()
			if i, ok := i.(fmt.Stringer); ok {
				io.WriteString(p.Out, i.String())
			} else {
				l := val.NumField()
				sOpen := val.Type().Name() + " {"
				if p.Compact {
					sOpen = "{"
				}
				if p.Compact && l == 0 {
					io.WriteString(p.Out, "{}")
				} else {
					io.WriteString(p.Out, sOpen+nl)
					// 计算最大字段名长度，以便对齐
					maxKeyLen := 0
					for i := 0; i < l; i++ {
						keyLen := len(val.Type().Field(i).Name)
						if keyLen > maxKeyLen {
							maxKeyLen = keyLen
						}
					}
					for i := 0; i < l; i++ {
						io.WriteString(p.Out, next)
						fieldName := val.Type().Field(i).Name

						// 计算当前键名后需要的空格数
						spaces := strings.Repeat(" ", maxKeyLen-len(fieldName))
						io.WriteString(p.Out, spaces) // 填充空格实现对齐
						io.WriteString(p.Out, fieldName)
						io.WriteString(p.Out, ": ")
						p.PrintValue(val.Field(i), level+1)
						if i < l-1 {
							io.WriteString(p.Out, ","+nl)
						} else {
							io.WriteString(p.Out, nl)
						}
					}
					io.WriteString(p.Out, cur)
					io.WriteString(p.Out, "}")
				}
			}
		} else {
			io.WriteString(p.Out, "protected")
		}
	default:
		io.WriteString(p.Out, "unsupported:")
		io.WriteString(p.Out, val.String())
	}
}
