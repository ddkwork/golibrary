package stream

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/ddkwork/golibrary/mylog"
	"golang.org/x/sync/errgroup"
)

func UpdateModsByWorkSpace(isTidy, isUpdateAll bool, modWithCommitID ...string) {
	mylog.Call(func() { updateModsByWorkSpace(isTidy, isUpdateAll, modWithCommitID...) })
}

var skips = []string{
	"module github.com/oligo/gioview",
	"module gioui.org",
	"module gioui.org/cmd",
	"module gioui.org/example",
	"module gioui.org/x",
}

func updateModsByWorkSpace(isTidy, isUpdateAll bool, modWithCommitID ...string) {
	if !IsFilePathEx("go.work") {
		return
	}
	RunCommandArgs("go work use -r .")
	lines := NewBuffer("go.work").ToLines()
	mods := make([]string, 0)
	for _, line := range lines {
		for _, skip := range skips {
			if line == skip {
				mylog.Warning("skip", skip)
				continue
			}
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, ".") {
			abs := mylog.Check2(filepath.Abs(line))
			mods = append(mods, abs)
		}
	}

	modChan := make(chan string, len(mods))

	var mutex sync.Mutex
	g := new(errgroup.Group)
	for _, mod := range mods {
		g.Go(func() error {
			mutex.Lock()
			defer mutex.Unlock()
			mylog.Check(os.Chdir(mod))
			for _, s := range modWithCommitID {
				RunCommand("go get -v  " + s)
			}
			if isTidy {
				RunCommand("go mod tidy -v")
			}
			if isUpdateAll {
				RunCommand("go get -v -u all")
			}
			modChan <- mod
			return nil
		})
	}
	go func() {
		for mod := range modChan {
			mylog.Success("updated mod", strconv.Quote(mod))
		}
		close(modChan)
	}()
	mylog.Check(g.Wait())
	mylog.Success("all work finished")
}
