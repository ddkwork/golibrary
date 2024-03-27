package golibrary

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
	"github.com/ddkwork/golibrary/stream/cmd"
)

func UpdateModsByWorkSpace(isTidy, isUpdateAll bool, modWithCommitID ...string) {
	if !stream.FileExists("go.work") {
		mylog.Error("go.work not found")
		return
	}
	cmd.RunArgs("go work use -r .")
	lines, ok := stream.NewReadFile("go.work").ToLines()
	if !ok {
		return
	}
	mods := make([]string, 0)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, ".") {
			abs, err := filepath.Abs(line)
			if !mylog.Error(err) {
				return
			}
			mods = append(mods, abs)
		}
	}

	//modChan := make(chan string, len(mods))

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

			if !mylog.Error(os.Chdir(abs)) {
				return
			}
			for _, s := range modWithCommitID {
				cmd.Run("go get -v  " + s)
			}
			if isTidy {
				cmd.Run("go mod tidy -v")
			}
			if isUpdateAll {
				cmd.Run("go get -v -u all")
			}
			if index > 0 {
				cmd.Run("gofumpt -l -w .") // default run gofumpt,工作区目录运行这个会死循环，原因未知
			}

			mylog.Success("updated mod", strconv.Quote(abs))

			//modChan <- abs
		}(mod, i)
	}

	//go func() {
	//	for mod := range modChan { // range 可以用于接收通道 <- 操作
	//		mylog.Success("updated mod", strconv.Quote(mod))
	//	}
	//	close(modChan)
	//}()

	wg.Wait()
	mylog.Success("all work finished")
}
