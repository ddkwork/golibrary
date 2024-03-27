// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package golibrary

import "runtime"

func IsAix() bool       { return runtime.GOOS == "aix" }
func IsAndroid() bool   { return runtime.GOOS == "android" }
func IsDarwin() bool    { return runtime.GOOS == "darwin" }
func IsDragonfly() bool { return runtime.GOOS == "dragonfly" }
func IsFreebsd() bool   { return runtime.GOOS == "freebsd" }
func IsHurd() bool      { return runtime.GOOS == "hurd" }
func IsIllumos() bool   { return runtime.GOOS == "illumos" }
func IsIos() bool       { return runtime.GOOS == "ios" }
func IsJs() bool        { return runtime.GOOS == "js" }
func IsLinux() bool     { return runtime.GOOS == "linux" }
func IsNacl() bool      { return runtime.GOOS == "nacl" }
func IsNetbsd() bool    { return runtime.GOOS == "netbsd" }
func IsOpenbsd() bool   { return runtime.GOOS == "openbsd" }
func IsPlan9() bool     { return runtime.GOOS == "plan9" }
func IsSolaris() bool   { return runtime.GOOS == "solaris" }
func IsWasip1() bool    { return runtime.GOOS == "wasip1" }
func IsWindows() bool   { return runtime.GOOS == "windows" }
func IsZos() bool       { return runtime.GOOS == "zos" }

// Note that this file is read by internal/goarch/gengoarch.go and by
// internal/goos/gengoos.go. If you change this file, look at those
// files as well.

// knownOS is the list of past, present, and future known GOOS values.
// Do not remove from this list, as it is used for filename matching.
// If you add an entry to this list, look at unixOS, below.
var knownOS = map[string]bool{
	"aix":       true,
	"android":   true,
	"darwin":    true,
	"dragonfly": true,
	"freebsd":   true,
	"hurd":      true,
	"illumos":   true,
	"ios":       true,
	"js":        true,
	"linux":     true,
	"nacl":      true,
	"netbsd":    true,
	"openbsd":   true,
	"plan9":     true,
	"solaris":   true,
	"wasip1":    true,
	"windows":   true,
	"zos":       true,
}

// unixOS is the set of GOOS values matched by the "unix" build tag.
// This is not used for filename matching.
// This list also appears in cmd/dist/build.go and
// cmd/go/internal/imports/build.go.
var unixOS = map[string]bool{
	"aix":       true,
	"android":   true,
	"darwin":    true,
	"dragonfly": true,
	"freebsd":   true,
	"hurd":      true,
	"illumos":   true,
	"ios":       true,
	"linux":     true,
	"netbsd":    true,
	"openbsd":   true,
	"solaris":   true,
}

// knownArch is the list of past, present, and future known GOARCH values.
// Do not remove from this list, as it is used for filename matching.
var knownArch = map[string]bool{
	"386":         true,
	"amd64":       true,
	"amd64p32":    true,
	"arm":         true,
	"armbe":       true,
	"arm64":       true,
	"arm64be":     true,
	"loong64":     true,
	"mips":        true,
	"mipsle":      true,
	"mips64":      true,
	"mips64le":    true,
	"mips64p32":   true,
	"mips64p32le": true,
	"ppc":         true,
	"ppc64":       true,
	"ppc64le":     true,
	"riscv":       true,
	"riscv64":     true,
	"s390":        true,
	"s390x":       true,
	"sparc":       true,
	"sparc64":     true,
	"wasm":        true,
}
