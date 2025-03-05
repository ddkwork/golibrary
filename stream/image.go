package stream

import (
	"bytes"
	"embed"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"path/filepath"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream/ico"
	"golang.org/x/image/bmp"
)

// 实现任意图片格式作为按钮图标

// paint.NewImageOp(LoadImage(src))
//func (i *Image) Layout(gtx layout.Context) layout.Dimensions {
//	return widget.Image{
//		Src:      i.imageOp,
//		Fit:      widget.Unscaled,
//		Position: layout.Center,
//		Scale:    1.0, //todo 测试按钮图标和层级图标并支持到button风格包
//	}.Layout(gtx)
//}

// LoadImage 一般而言，我们使用embed,对于进程图标，制作任务管理器和音速启动，直接得到image.Image
// 所以我们不用懒解码探测图片格式
func LoadImage(fileName string, fs embed.FS) image.Image {
	b := mylog.Check2(fs.ReadFile(fileName))
	var img image.Image
	switch filepath.Ext(fileName) {
	case ".png":
		img = mylog.Check2(png.Decode(bytes.NewReader(b)))
	case ".jpg", ".jpeg":
		img = mylog.Check2(jpeg.Decode(bytes.NewReader(b)))
	case ".gif":
		img = mylog.Check2(gif.Decode(bytes.NewReader(b)))
	case ".ico":
		img = mylog.Check2(ico.Decode(bytes.NewReader(b)))
	case ".bmp":
		img = mylog.Check2(bmp.Decode(bytes.NewReader(b)))
	// case ".svg": //svg的话是giosvg直接解码元数据实现layout方法渲染
	default:
		mylog.Check("unsupported image format")
	}
	return img
}
