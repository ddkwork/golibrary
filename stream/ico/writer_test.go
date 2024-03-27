package ico

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/ddkwork/golibrary/mylog"
)

func TestEncode(t *testing.T) {
	t.Parallel()
	origfile := "testdata/golang.ico"
	file := "testdata/golang_test.ico"

	f, err := os.Open("testdata/golang.png")
	img, err := png.Decode(f)
	if !mylog.Error(err) {
		return
	}
	f.Close()

	var newFile *os.File
	if newFile, err = os.Create(filepath.Join(file)); err != nil {
		t.Error(err)
	}
	if !Encode(newFile, img) {
		return
	}
	newFile.Close()

	f, err = os.Open(origfile)
	if !mylog.Error(err) {
		return
	}
	origICO, err := Decode(f)
	if !mylog.Error(err) {
		return
	}
	f.Close()

	newFile, err = os.Open(file)
	if !mylog.Error(err) {
		return
	}
	newICO, err := Decode(newFile)
	if !mylog.Error(err) {
		return
	}
	newFile.Close()

	inrgba, ok := origICO.(*image.NRGBA)
	if !ok {
		t.Fatal("not nrgba")
	}
	pnrgba, ok := newICO.(*image.NRGBA)
	if !ok {
		t.Fatal("new not nrgba")
	}
	if b, err := fastCompare(inrgba, pnrgba); err != nil || b != 0 {
		t.Fatalf("pix differ %d %v\n", b, err)
	}
}
