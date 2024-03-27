package stream

import (
	"image/png"
	"os"
	"testing"

	"github.com/ddkwork/golibrary/mylog"
)

func TestExtractIcon2Png(t *testing.T) {
	path := "D:\\app\\Internet Download Manager 6.42 Build 1 多语+Retail 坡姐\\Internet Download Manager 6.42 Build 1 多语+Retail 坡姐\\Patch-Ali.Dbg_v18.2\\IDM v.6.4x crack v.18.2.exe"
	path = "C:\\Windows\\notepad.exe"
	image, ok := ExtractIcon2Image(path)
	if !ok {
		return
	}
	f, err := os.Create("1.png")
	if err != nil {
		panic(err.Error())
	}
	if !mylog.Error(png.Encode(f, image)) {
		return
	}
	mylog.Error(f.Close())
	os.Remove("1.png")
}
