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

func UpdateModsByWorkSpace(isUpdateAll bool) {
	mylog.Call(func() { updateModsByWorkSpace(isUpdateAll) })
}

var skips = []string{
	"module gioui.org",
	"module gioui.org/cmd",
	"module gioui.org/example",
	"module gioui.org/x",
}

func UpdateDependencies() {
	for s := range strings.Lines(`
     go get -x gioui.org@main
	 go get -x gioui.org/cmd@main
	 go get -x gioui.org/example@main
	 go get -x gioui.org/x@main
	 go get -x github.com/oligo/gvcode@main
	 go get -x github.com/ddkwork/golibrary@master
	 go get -x github.com/ddkwork/ux@master
	 go get -x github.com/google/go-cmp@master
	 go get -x github.com/ddkwork/app@master
	 go get -x github.com/ddkwork/toolbox@master
	 go get -x github.com/ddkwork/unison@master
	 go get -x github.com/ebitengine/purego@main
	 go get -x github.com/saferwall/pe@main
	 ::go get -u -x all
	 go mod tidy

	go install mvdan.cc/gofumpt@latest
	gofumpt -l -w .
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

`) {
		s = strings.TrimSpace(s)
		if strings.HasPrefix(s, "::") || strings.HasPrefix(s, "//") || s == "" {
			continue
		}
		RunCommand(s)
	}
}

func updateModsByWorkSpace(isUpdateAll bool) {
	if !IsFilePathEx("go.work") {
		return
	}
	RunCommandArgs("go work use -r .")
	mods := make([]string, 0)
	for line := range ReadFileToLines("go.work") {
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
			UpdateDependencies()
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

//type Cache struct {
//	store map[string]string
//	mu    sync.RWMutex
//}
//
//func TestCacheConsistency(t *testing.T) {
//	cache := &Cache{store: make(map[string]string)}
//
//	synctest.Run(t, func(sc *synctest.Scenario) {
//		// 并发写入
//		for i := 0; i < 10; i++ {
//			sc.Go(func() {
//				cache.mu.Lock()
//				defer cache.mu.Unlock()
//				cache.store["key"] = time.Now().String()
//			})
//		}
//
//		// 并发读取
//		for i := 0; i < 100; i++ {
//			sc.Go(func() {
//				cache.mu.RLock()
//				defer cache.mu.RUnlock()
//				_ = cache.store["key"]
//			})
//		}
//	}, synctest.WithOptions(
//		synctest.EnableRaceDetection(),
//		synctest.MaxGoroutines(200),
//		synctest.Timeout(10*time.Second),
//	))
//}
