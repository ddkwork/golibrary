package stream

import (
	"image"
	"image/color"
	"syscall"
	"unsafe"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/zzl/go-win32api/v2/win32"
	"golang.org/x/sys/windows"
)

func ExtractIcon2Image(filename string) (image image.Image, ok bool) {
	return ExtractPrivateExtractIcons(filename, 128, 128)
}

func ExtractIconToImageByExt(filename string) (img image.Image, okk bool) {
	var shFile win32.SHFILEINFO
	err := win32.SHGetFileInfoW(
		win32.StrToPwstr(filename),
		0,
		&shFile,
		uint32(unsafe.Sizeof(shFile)),
		win32.SHGFI_ICON|win32.SHGFI_USEFILEATTRIBUTES)
	if !isTrue(int32(err)) {
		return
	}
	defer win32.DestroyIcon(shFile.HIcon)
	return hICONToImage(shFile.HIcon)
}

func ExtractIconToImage(filename string) (img image.Image, okk bool) {
	large := []win32.HICON{0}
	_, err := win32.ExtractIconExW(win32.StrToPwstr(filename), 0, &large[0], nil, 1)
	if err.NilOrError() != nil {
		return
	}
	defer win32.DestroyIcon(large[0])
	return hICONToImage(large[0])
}

func isTrue[T uint32 | int32](u T) bool { return u != 0 }

// ExtractPrivateExtractIcons 提取exe高清图标
func ExtractPrivateExtractIcons(filename string, w, h int32) (img image.Image, okk bool) {
	large := []win32.HICON{0}
	var piconId uint32 = 0
	ok := win32.PrivateExtractIcons(win32.StrToPwstr(filename), 0, w, h, &large[0], &piconId, 1, 0)
	if !isTrue(ok) {
		return
	}
	if large == nil {
		mylog.Error("can't extract icon")
		return
	}
	defer win32.DestroyIcon(large[0])
	return hICONToImage(large[0])
}

func hICONToImage(hicon win32.HICON) (img image.Image, okk bool) {
	var iconInfo win32.ICONINFO
	ok, err := win32.GetIconInfo(hicon, &iconInfo)
	if err.NilOrError() != nil {
		return
	}
	w := int32(iconInfo.XHotspot * 2)
	h := int32(iconInfo.YHotspot * 2)

	dc := win32.GetDC(0)
	hdc := win32.CreateCompatibleDC(dc)
	ok = win32.ReleaseDC(0, hdc)
	if !isTrue(ok) {
		// mylog.Error("ReleaseDC Error")
		// return
	}
	defer win32.DeleteDC(hdc)
	var bits unsafe.Pointer
	info := &win32.BITMAPINFO{
		BmiHeader: win32.BITMAPINFOHEADER{
			BiSize:          0,
			BiWidth:         w,
			BiHeight:        -h,
			BiPlanes:        1,
			BiBitCount:      32,
			BiCompression:   win32.BI_RGB,
			BiSizeImage:     uint32(w * h * 4),
			BiXPelsPerMeter: 0,
			BiYPelsPerMeter: 0,
			BiClrUsed:       0,
			BiClrImportant:  0,
		},
		BmiColors: [1]win32.RGBQUAD{},
	}
	info.BmiHeader.BiSize = uint32(unsafe.Sizeof(info.BmiHeader))
	winBitmap, err := CreateDIBSection(hdc, info, win32.DIB_RGB_COLORS, &bits, 0, 0)
	defer win32.DeleteObject(winBitmap)

	pixels := (*[1 << 30]byte)(bits)[0:info.BmiHeader.BiSizeImage]
	win32.SelectObject(hdc, winBitmap)
	ok, err = win32.DrawIconEx(hdc, 0, 0, hicon, w, h, 0, 0, win32.DI_NORMAL)
	if !mylog.Error(err.NilOrError()) {
		return
	}
	if !isTrue(ok) {
		mylog.Error("DrawIcon Error")
		return
	}
	hasAlpha := false
	rgba := image.NewRGBA(image.Rectangle{
		Min: image.Point{
			X: 0,
			Y: 0,
		},
		Max: image.Point{
			X: int(w),
			Y: int(h),
		},
	})
	for y := range h {
		for x := range w {
			if pixels[((y*w+x)*4)+3] > 0 {
				hasAlpha = true
			}

			rgba.SetRGBA(int(x), int(y), color.RGBA{
				A: pixels[((y*w+x)*4)+3],
				R: pixels[((y*w+x)*4)+2],
				G: pixels[((y*w+x)*4)+1],
				B: pixels[((y*w+x)*4)+0],
			})
		}
	}
	if hasAlpha {
		return rgba, true
	}
	ok, _ = win32.DrawIconEx(hdc, 0, 0, hicon, w, h, 0, 0, win32.DI_MASK)
	if isTrue(ok) {
		for y := range h {
			for x := range w {
				tmp := rgba.RGBAAt(int(x), int(y))
				if (pixels[((y*w+x)*4)+2] | pixels[((y*w+x)*4)+1] | pixels[((y*w+x)*4)+0]) == 0 {
					tmp.A = 0xFF
					rgba.SetRGBA(int(x), int(y), tmp)
				}
			}
		}
	}
	return rgba, true
}

func CreateDIBSection(hdc win32.HDC, pbmi *win32.BITMAPINFO, usage win32.DIB_USAGE, ppvBits *unsafe.Pointer, hSection win32.HANDLE, offset uint32) (win32.HBITMAP, win32.WIN32_ERROR) {
	libGdi32 := windows.NewLazySystemDLL("gdi32.dll")
	pCreateDIBSection := uintptr(0)
	addr := win32.LazyAddr(&pCreateDIBSection, libGdi32, "CreateDIBSection")
	ret, _, err := syscall.SyscallN(addr, hdc, uintptr(unsafe.Pointer(pbmi)), uintptr(usage), uintptr(unsafe.Pointer(ppvBits)), hSection, uintptr(offset))
	return ret, win32.WIN32_ERROR(err)
}
