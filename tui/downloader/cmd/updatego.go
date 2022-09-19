package main

import (
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/src/stream/tool/cmd"
	"github.com/ddkwork/golibrary/src/unpackit"
	"github.com/ddkwork/golibrary/tui/downloader"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func runCommand(command string) error {
	parts := strings.Fields(command)
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stderr = os.Stdin
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func UpdateGo(url string) (ok bool) {
	fileName := filepath.Base(url)
	ext := filepath.Ext(url)
	goRoot := runtime.GOROOT()
	switch runtime.GOOS {
	case "windows":
		return true //todo
	case "linux":
		if ext != ".gz" {
			if !downloader.Run(url) {
				return
			}
		}
	}
	if !mylog.Error(os.RemoveAll(goRoot)) {
		return
	}
	//cmd.Run()
	//if !mylog.Error(runCommand(fmt.Sprintf("sudo tar -C %s -xzf %s", name, filepath.Base(url)))) {
	//	return
	//}
	if !unpackit.Run(fileName, goRoot) {
		return
	}
	if !mylog.Error(os.RemoveAll(fileName)) {
		return
	}
	mylog.Success("", "DONE")
	run, err := cmd.Run("go version")
	if !mylog.Error(err) {
		return
	}
	mylog.Success("finish", run)
	return true
}
