package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/widget/downloader"

	"github.com/ddkwork/golibrary/stream/cmd"
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
	return ok
	fileName := filepath.Base(url)
	ext := filepath.Ext(url)
	goRoot := runtime.GOROOT()
	switch runtime.GOOS {
	case "windows":
		return true // todo
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
	//if !unpackit.Run(fileName, goRoot) {
	//	return
	//}
	if !mylog.Error(os.RemoveAll(fileName)) {
		return
	}
	mylog.Success("", "DONE")
	run := cmd.Run("go version")
	mylog.Success("finish", run)
	return true
}
