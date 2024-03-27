package ico

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"testing"

	"github.com/ddkwork/golibrary/mylog"
)

func sqDiffUInt8(x, y uint8) uint64 {
	d := uint64(x) - uint64(y)
	return d * d
}

func fastCompare(img1, img2 *image.NRGBA) (int64, error) {
	if img1.Bounds() != img2.Bounds() {
		return 0, fmt.Errorf("image bounds not equal: %+v, %+v", img1.Bounds(), img2.Bounds())
	}

	accumError := int64(0)
	for i := range len(img1.Pix) {
		accumError += int64(sqDiffUInt8(img1.Pix[i], img2.Pix[i]))
	}

	return int64(math.Sqrt(float64(accumError))), nil
}

func aTestDecodeConfig(t *testing.T) {
	t.Parallel()
	file := "testdata/golang.ico"
	copyFile := "testdata/golang.png"
	reader, err := os.Open(file)
	if !mylog.Error(err) {
		return
	}
	icoImage, err := DecodeConfig(reader)
	reader.Close()
	if !mylog.Error(err) {
		return
	}
	reader, err = os.Open(copyFile)
	if !mylog.Error(err) {
		return
	}
	pngImage, err := png.DecodeConfig(reader)
	reader.Close()
	if !mylog.Error(err) {
		return
	}

	if icoImage != pngImage {
		t.Errorf("%v - %v", icoImage, pngImage)
	}
}

func TestDecode(t *testing.T) {
	t.Parallel()
	file := "testdata/golang.ico"
	copyFile := "testdata/golang.png"
	reader, err := os.Open(file)
	if !mylog.Error(err) {
		return
	}
	icoImage, err := Decode(reader)
	if !mylog.Error(err) {
		return
	}
	if icoImage == nil {
		return
	}
	reader.Close()

	reader, err = os.Open(copyFile)
	if !mylog.Error(err) {
		return
	}
	pngImage, err := png.Decode(reader)
	if !mylog.Error(err) {
		return
	}
	if pngImage == nil {
		return
	}
	reader.Close()

	if icoImage == nil || !icoImage.Bounds().Eq(pngImage.Bounds()) {
		t.Fatal("bounds differ")
	}
	inrgba, ok := icoImage.(*image.NRGBA)
	if !ok {
		t.Fatal("not nrgba")
	}
	pnrgba, ok := pngImage.(*image.NRGBA)
	if !ok {
		t.Fatal("png not nrgba")
	}

	if b, err := fastCompare(inrgba, pnrgba); err != nil || b != 670 { // todo
		t.Fatalf("pix differ %d %v\n", b, err)
	}
}
