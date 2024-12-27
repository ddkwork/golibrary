package txt

import "io"

type RuneReader struct {
	Src []rune
	Pos int
}

func (rr *RuneReader) ReadRune() (r rune, size int, err error) {
	if rr.Pos >= len(rr.Src) {
		return -1, 0, io.EOF
	}
	nextRune := rr.Src[rr.Pos]
	rr.Pos++
	return nextRune, 1, nil
}
