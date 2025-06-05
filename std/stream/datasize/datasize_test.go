package datasize

import (
	"testing"

	"github.com/ddkwork/golibrary/std/mylog"
)

func TestName(t *testing.T) {
	println((10480 * KB).String())
}

func TestMarshalText(t *testing.T) {
	table := []struct {
		in  Size
		out string
	}{
		{0, "0B"},
		{B, "1B"},
		{KB, "1KB"},
		{MB, "1MB"},
		{GB, "1GB"},
		{TB, "1TB"},
		{PB, "1PB"},
		{EB, "1EB"},
		{400 * TB, "400TB"},
		{2048 * MB, "2GB"},
		{B + KB, "1025B"},
		{MB + 20*KB, "1044KB"},
		{100*MB + KB, "102401KB"},
	}

	for _, tt := range table {
		b := mylog.Check2(tt.in.MarshalText())
		s := string(b)

		if s != tt.out {
			t.Errorf("MarshalText(%d) => %s, want %s", tt.in, s, tt.out)
		}
	}
}

func TestUnmarshalText(t *testing.T) {
	t.Skip("not implemented")
	table := []struct {
		in  string
		err bool
		out Size
	}{
		{"0", false, 0},
		{"0B", false, 0},
		{"0 KB", false, 0},
		{"1", false, B},
		{"1K", false, KB},
		{"2MB", false, 2 * MB},
		{"5 GB", false, 5 * GB},
		{"20480 G", false, 20 * TB},
		{"50 eB", true, (1 << 64) - 1},
		{"200000 pb", true, (1 << 64) - 1},
		{"10 Mb", true, 0},
		{"g", true, 0},
		{"10 kB ", false, 10 * KB},
		{"10 kBs ", true, 0},
	}

	for _, tt := range table {
		t.Run("UnmarshalText "+tt.in, func(t *testing.T) {
			var s Size
			s.UnmarshalText([]byte(tt.in))

			if s != tt.out {
				t.Errorf("UnmarshalText(%s) => %d bytes, want %d bytes", tt.in, s, tt.out)
			}
		})
		t.Run("Parse "+tt.in, func(t *testing.T) {
			s := Parse([]byte(tt.in))
			if s != tt.out {
				t.Errorf("Parse(%s) => %d bytes, want %d bytes", tt.in, s, tt.out)
			}
		})
		t.Run("MustParse "+tt.in, func(t *testing.T) {
			s := Parse([]byte(tt.in))
			if s != tt.out {
				t.Errorf("MustParse(%s) => %d bytes, want %d bytes", tt.in, s, tt.out)
			}
		})
		t.Run("ParseString "+tt.in, func(t *testing.T) {
			s := Parse(tt.in)

			if s != tt.out {
				t.Errorf("ParseString(%s) => %d bytes, want %d bytes", tt.in, s, tt.out)
			}
		})
		t.Run("MustParseString "+tt.in, func(t *testing.T) {
			s := Parse(tt.in)
			if s != tt.out {
				t.Errorf("MustParseString(%s) => %d bytes, want %d bytes", tt.in, s, tt.out)
			}
		})
	}
}

func Test_parseSizeAndUnit(t *testing.T) {
	mylog.Call(func() {
		parseSizeAndUnit("5.2GB")
		parseSizeAndUnit("5.2 GB")
	})
}
