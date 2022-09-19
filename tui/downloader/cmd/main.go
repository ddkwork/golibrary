package main

import (
	"github.com/ddkwork/golibrary/tui/downloader"
)

//go:generate  go build .

func main() { //todo test
	if !UpdateGo("https://go.dev/dl/go1.19.linux-amd64.tar.gz") {
		return
	}
	return
	if !downloader.Run("http://download.jieiis.com/iso/windows/cn_windows_7_ultimate_with_sp1_x64_dvd_u_677408.iso") {
		return
	}
}
