package unpackit

import (
	"github.com/c4milo/unpackit"

	"github.com/ddkwork/golibrary/mylog"
	"os"
)

// todo addprogress and pack pai
// https://github.com/a5272689/filetools
func Run(fileName, dstDir string) (ok bool) {
	file, err := os.Open(fileName)
	if !mylog.Error(err) {
		return
	}
	destPath, err := unpackit.Unpack(file, dstDir)
	if !mylog.Error(err) {
		return
	}
	mylog.Success("unpacked to", destPath)
	return mylog.Error(file.Close())
}
