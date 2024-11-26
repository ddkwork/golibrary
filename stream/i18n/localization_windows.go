package i18n

import (
	"syscall"
	"unsafe"
)

func Locale() string {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("GetUserDefaultLocaleName")
	buffer := make([]uint16, 128)
	if ret, _, _ := proc.Call(uintptr(unsafe.Pointer(&buffer[0])), uintptr(len(buffer))); ret == 0 {
		proc = kernel32.NewProc("GetSystemDefaultLocaleName")
		if ret, _, _ = proc.Call(uintptr(unsafe.Pointer(&buffer[0])), uintptr(len(buffer))); ret == 0 {
			return "en_US.UTF-8"
		}
	}
	return syscall.UTF16ToString(buffer)
}
