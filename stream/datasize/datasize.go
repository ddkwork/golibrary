package datasize

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
)

type Size uint64

const (
	B  Size = 1
	KB      = B << 10
	MB      = KB << 10
	GB      = MB << 10
	TB      = GB << 10
	PB      = TB << 10
	EB      = PB << 10

	// fnUnmarshalText string = "UnmarshalText"
	// maxUint64       uint64 = (1 << 64) - 1
	// cutoff          uint64 = maxUint64 / 10
)

var ErrBits = errors.New("unit with capital unit prefix and lower case unit (b) - bits, not bytes")

func (s Size) Bytes() uint64 {
	return uint64(s)
}

func (s Size) KBytes() float64 {
	v := s / KB
	r := s % KB
	return float64(v) + float64(r)/float64(KB)
}

func (s Size) MBytes() float64 {
	v := s / MB
	r := s % MB
	return float64(v) + float64(r)/float64(MB)
}

func (s Size) GBytes() float64 {
	v := s / GB
	r := s % GB
	return float64(v) + float64(r)/float64(GB)
}

func (s Size) TBytes() float64 {
	v := s / TB
	r := s % TB
	return float64(v) + float64(r)/float64(TB)
}

func (s Size) PBytes() float64 {
	v := s / PB
	r := s % PB
	return float64(v) + float64(r)/float64(PB)
}

func (s Size) EBytes() float64 {
	v := s / EB
	r := s % EB
	return float64(v) + float64(r)/float64(EB)
}

func (s Size) String() string {
	switch {
	case s > EB:
		return fmt.Sprintf("%.1f EB", s.EBytes())
	case s > PB:
		return fmt.Sprintf("%.1f PB", s.PBytes())
	case s > TB:
		return fmt.Sprintf("%.1f TB", s.TBytes())
	case s > GB:
		return fmt.Sprintf("%.1f GB", s.GBytes())
	case s > MB:
		return fmt.Sprintf("%.1f MB", s.MBytes())
	case s > KB:
		return fmt.Sprintf("%.1f KB", s.KBytes())
	default:
		return fmt.Sprintf("%d B", s)
	}
}

func (s Size) MachineString() string {
	switch {
	case s == 0:
		return "0B"
	case s%EB == 0:
		return fmt.Sprintf("%dEB", s/EB)
	case s%PB == 0:
		return fmt.Sprintf("%dPB", s/PB)
	case s%TB == 0:
		return fmt.Sprintf("%dTB", s/TB)
	case s%GB == 0:
		return fmt.Sprintf("%dGB", s/GB)
	case s%MB == 0:
		return fmt.Sprintf("%dMB", s/MB)
	case s%KB == 0:
		return fmt.Sprintf("%dKB", s/KB)
	default:
		return fmt.Sprintf("%dB", s)
	}
}

func (s Size) MarshalText() ([]byte, error) {
	return []byte(s.MachineString()), nil
}

func parseSizeAndUnit(input string) (sizeStr string, unitStr string) {
	input = strings.TrimSpace(input)
	for i, char := range input {
		if unicode.IsDigit(char) || char == '.' {
			sizeStr += string(char)
		} else {
			unitStr = strings.ToUpper(strings.TrimSpace(input[i:]))
			break
		}
	}
	return
}

func (s *Size) UnmarshalText(b []byte) {
	if len(b) == 0 {
		return
	}
	mylog.Check(b)

	sizeStr, unit := parseSizeAndUnit(string(b))
	var size float64

	if strings.Contains(sizeStr, ".") {
		size = stream.ParseFloat(sizeStr)
	} else {
		size = float64(stream.ParseUint(sizeStr))
	}

	switch unit {
	case "Kb", "Mb", "Gb", "Tb", "Pb", "Eb":
		mylog.Check(ErrBits)
	}
	switch unit {
	case "B":
	case "KB":
		size *= 1024
	case "MB":
		size *= 1024 * 1024
	case "GB":
		size *= 1024 * 1024 * 1024
	case "TB":
		size *= 1024 * 1024 * 1024 * 1024
	case "PB":
		size *= 1024 * 1024 * 1024 * 1024 * 1024
	case "EB":
		size *= 1024 * 1024 * 1024 * 1024 * 1024 * 1024
	}
	*s = Size(size)
}

func Parse[T string | []byte](data T) (size Size) {
	switch data := any(data).(type) {
	case string:
		size.UnmarshalText([]byte(data))
	case []byte:
		size.UnmarshalText(data)
	default:
		panic(fmt.Sprintf("unsupported type %T", data))
	}
	return
}
