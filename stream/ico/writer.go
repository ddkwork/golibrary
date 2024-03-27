package ico

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"image"
	"image/draw"
	"image/png"
	"io"

	"github.com/ddkwork/golibrary/mylog"
)

func Encode(w io.Writer, im image.Image) (ok bool) {
	if w == nil {
		return
	}
	b := im.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, im, b.Min, draw.Src)

	header := head{
		0,
		1,
		1,
	}
	entry := direntry{
		Plane:  1,
		Bits:   32,
		Offset: 22,
	}

	pngbuffer := new(bytes.Buffer)
	pngwriter := bufio.NewWriter(pngbuffer)
	if !mylog.Error(png.Encode(pngwriter, m)) {
		return
	}
	if !mylog.Error(pngwriter.Flush()) {
		return
	}
	entry.Size = uint32(len(pngbuffer.Bytes()))

	bounds := m.Bounds()
	entry.Width = uint8(bounds.Dx())
	entry.Height = uint8(bounds.Dy())
	bb := new(bytes.Buffer)

	if !mylog.Error(binary.Write(bb, binary.LittleEndian, header)) {
		return
	}

	if !mylog.Error(binary.Write(bb, binary.LittleEndian, entry)) {
		return
	}

	if !mylog.Error2(w.Write(bb.Bytes())) {
		return
	}

	return mylog.Error2(w.Write(pngbuffer.Bytes()))
}
