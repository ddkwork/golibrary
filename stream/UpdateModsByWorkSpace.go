package stream

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/ddkwork/golibrary/mylog"
)

func UpdateModsByWorkSpace(isTidy, isUpdateAll bool, modWithCommitID ...string) {
	mylog.Call(func() { updateModsByWorkSpace(isTidy, isUpdateAll, modWithCommitID...) })
}

func updateModsByWorkSpace(isTidy, isUpdateAll bool, modWithCommitID ...string) {
	if !IsFilePathEx("go.work") {
		return
	}
	RunCommandArgs("go work use -r .")
	lines := NewBuffer("go.work").ToLines()
	mods := make([]string, 0)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, ".") {
			abs := mylog.Check2(filepath.Abs(line))
			mods = append(mods, abs)
		}
	}

	modChan := make(chan string, len(mods))

	var (
		wg    sync.WaitGroup
		mutex sync.Mutex
	)

	for i, mod := range mods {
		wg.Add(1)
		go func(abs string, index int) {
			defer wg.Done()

			mutex.Lock()
			defer mutex.Unlock()

			mylog.Check(os.Chdir(abs))
			for _, s := range modWithCommitID {
				RunCommand("go get -v  " + s)
			}
			if isTidy {
				RunCommand("go mod tidy -v")
			}
			if isUpdateAll {
				RunCommand("go get -v -u all")
			}

			modChan <- abs
		}(mod, i)
	}

	go func() {
		for mod := range modChan {
			mylog.Success("updated mod", strconv.Quote(mod))
		}
		close(modChan)
	}()

	wg.Wait()
	mylog.Success("all work finished")
}
