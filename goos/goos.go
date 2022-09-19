package goos

import (
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/src/stream/tool"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

const (
	android = `android`
	linux   = `linux`
	windows = `windows`
)

func IsAndroid() bool { return runtime.GOOS == android }
func IsWindows() bool { return runtime.GOOS == windows }
func IsLinux() bool   { return runtime.GOOS == linux }

func AppInstallPaths() (installPaths []string) {
	installPaths = make([]string, 0)
	switch {
	case IsWindows():
		installPaths = []string{"C:/Program Files/JetBrains"}
	case IsLinux():
		installPaths = []string{"/opt", "/usr/share"}
	case IsAndroid():
		installPaths = []string{"/data/data/com.termux/files/usr/opt"}
	}
	return
}
func UserConfigDir(username string) (configDir string, ok bool) {
	if username == "" {
		username = "ddk"
	}
	configDir, err := os.UserConfigDir()
	if !mylog.Error(err) {
		return
	}
	switch {
	case IsWindows():
	case IsLinux():
		if runtime.GOOS == "linux" {
			lookup, err := user.Lookup(username)
			if !mylog.Error(err) {
				return
			}
			configDir = lookup.HomeDir
			configDir += "/.config"
		}
	case IsAndroid():
		configDir = HomeDir() + `/.config`
	}
	ok = true
	return
}

func InitUsername() (username string, ok bool) {
	usernamePath := "username.txt"
	_, err := os.Stat(usernamePath)
	if os.IsNotExist(err) {
		tool.File().WriteTruncate(usernamePath, "ddk")
	}
	lines, o := tool.File().ReadToLines("username.txt")
	if !o {
		return
	}
	return lines[0], true
}

func HomeDir() string {
	switch {
	case IsAndroid():
		return `/data/data/com.termux/files/home`
	case IsLinux():
		return `/home`
	case IsWindows():
		return `C:\Users\Admin`
	}
	mylog.Error("HomeDir not find")
	return ""
}

func UserHomeDir() (userHomeDir string, ok bool) {
	username, ok := InitUsername()
	if !ok {
		return
	}
	if IsWindows() {
		mylog.Info("userHomeDir", HomeDir())
		return HomeDir(), true
	}
	userHomeDir = filepath.Join(HomeDir(), username)
	mylog.Info("userHomeDir", userHomeDir)
	ok = true
	return
}

func GoBin() (GoBin string, ok bool) {
	dir, ok := UserHomeDir()
	if !ok {
		return
	}
	GoBin = filepath.Join(dir, "go/bin/")
	mylog.Info("GoBin", GoBin)
	ok = true
	return
}
