package goos

import "runtime"

const (
	android = `android`
	linux   = `linux`
	windows = `windows`
)

func IsAndroid() bool { return runtime.GOOS == android }
func IsWindows() bool { return runtime.GOOS == windows }
func IsLinux() bool   { return runtime.GOOS == linux }
