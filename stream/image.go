package stream

import (
	"bytes"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream/ico"
	"golang.org/x/image/bmp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

//实现任意图片格式作为按钮图标

// paint.NewImageOp(LoadImage(src))
//func (i *Image) Layout(gtx layout.Context) layout.Dimensions {
//	return widget.Image{
//		Src:      i.imageOp,
//		Fit:      widget.Unscaled,
//		Position: layout.Center,
//		Scale:    1.0, //todo 测试按钮图标和层级图标并支持到button风格包
//	}.Layout(gtx)
//}

func LoadImage(fileName string) image.Image {
	file := mylog.Check2(os.ReadFile(fileName))
	var img image.Image
	switch filepath.Ext(fileName) {
	case ".png":
		img = mylog.Check2(png.Decode(bytes.NewReader(file)))
	case ".jpg", ".jpeg":
		img = mylog.Check2(jpeg.Decode(bytes.NewReader(file)))
	case ".gif":
		img = mylog.Check2(gif.Decode(bytes.NewReader(file)))
	case ".ico":
		img = mylog.Check2(ico.Decode(bytes.NewReader(file)))
	case ".bmp":
		img = mylog.Check2(bmp.Decode(bytes.NewReader(file)))
	//case ".svg": //svg的话是giosvg直接解码元数据实现layout方法渲染
	default:
		mylog.Check("unsupported image format")
	}
	return img
}
