package stream

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/safemap"
	"golang.org/x/mod/modfile"
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

// GetLastCommitHash
// git ls-remote https://github.com/ddkwork/toolbox refs/heads/master
// git ls-remote https://github.com/gioui/gio refs/heads/main
// 但是要传递本地的仓库目录，太麻烦了
//func GetLastCommitHash(repositoryName string) string {
//	//理论上获取到hash再使用模块代理，这样才刷新的快？
//	//或者使用action得到hash先？
//
//	mylog.Check(os.Setenv("GOPROXY", "direct")) // 如果模块代理导致获取到的不是最新的提交哈希那么需要禁用模块代理，最可靠的方式是 GetLastCommitHashLocal
//	defer mylog.Check(os.Setenv("GOPROXY", "https://goproxy.cn,direct"))
//
//	//defer RunCommand("go env -w GOPROXY=https://goproxy.cn,direct")
//	s := RunCommand("git ls-remote " + repositoryName + " refs/heads/master").Output.String()
//	for hash := range strings.FieldsSeq(s) {
//		return hash
//	}
//	s = RunCommand("git ls-remote " + repositoryName + " refs/heads/main").Output.String()
//	for hash := range strings.FieldsSeq(s) {
//		return hash
//	}
//	panic("no commit hash found in master and main branch")
//	/*
//		# 获取最新提交哈希
//		$hash = (git ls-remote https://github.com/ddkwork/toolbox refs/heads/master).Split("`t")[0]
//		# 带哈希安装
//		go get -x "github.com/ddkwork/toolbox@$hash"
//
//		powershell:
//			$env:GOPROXY="direct"; go get -x github.com/ddkwork/toolbox@$(git ls-remote https://github.com/ddkwork/toolbox refs/heads/master | ForEach-Object { $_.Split()[0] })
//
//	*/
//}

func GetLastCommitHashLocal(repositoryDir string) string {
	originPath := mylog.Check2(os.Getwd())
	mylog.Check(os.Chdir(repositoryDir))
	hash := RunCommand("git rev-parse HEAD").Output.String()
	mylog.Check(os.Chdir(originPath))
	return hash
}

func ParseGoMod() *safemap.M[string, string] {
	path := "go.mod"
	f := mylog.Check2(modfile.Parse(path, mylog.Check2(os.ReadFile(path)), nil))
	return safemap.NewOrdered[string, string](func(yield func(string, string) bool) {
		for _, req := range f.Require {
			yield(req.Mod.Path, req.Mod.Version)
		}
	})
}

func GetDesktopDir() string {
	// 获取用户主目录
	homeDir := mylog.Check2(os.UserHomeDir())
	// 根据操作系统处理路径
	switch runtime.GOOS {
	case "windows", "darwin":
		// Windows和macOS直接拼接Desktop
		return filepath.Join(homeDir, "Desktop")
	case "linux":
		// Linux优先检查XDG环境变量
		if xdgDir := os.Getenv("XDG_DESKTOP_DIR"); xdgDir != "" {
			return xdgDir
		}
		// 默认使用主目录下的Desktop
		return filepath.Join(homeDir, "Desktop")
	default:
		panic("unsupported platform")
	}
}

func UpdateDependencies() { // 模块代理刷新的不及时，需要禁用代理,已经使用clone仓库远程完成更新
	var mutex sync.Mutex
	g := new(errgroup.Group)
	for s := range ReadFileToLines(filepath.Join(GetDesktopDir(), "dep.txt")) { // 因为要经常更新，我们不embed
		s = strings.TrimSpace(s)
		if strings.HasPrefix(s, "::") || strings.HasPrefix(s, "//") || s == "" {
			continue
		}
		if s == "go mod tidy" {
			continue
		}
		g.Go(func() error { // 这样之后tidy就不在最后执行了，同时升级多个依赖+读写锁定
			mutex.Lock()
			RunCommand(s)
			mutex.Unlock()
			return nil
		})
		mylog.Check(g.Wait())
	}
	RunCommand("go mod tidy") // 所有yield都执行完了，再执行tidy
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

	g := new(errgroup.Group)
	for _, path := range mods {
		g.Go(func() error { // 每个模块单独跑,这里不能加锁，否则很慢，谨慎使用读写锁
			mylog.Check(os.Chdir(path)) // 这里的必报参数不用管？
			UpdateDependencies()        // 锁应该在这里面
			if isUpdateAll {
				RunCommand("go get -u -x all")
			}
			modChan <- path
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

//func updateDependencies() { // 模块代理刷新的不及时，需要禁用代理,已经使用clone仓库远程完成更新
//	mylog.Check(os.Setenv("GOPROXY", "direct"))
//	for s := range strings.Lines(`
//     go get -x gioui.org@main
//	 go get -x gioui.org/cmd@main
//	 go get -x gioui.org/example@main
//	 go get -x gioui.org/x@main
//	 go get -x github.com/oligo/gvcode@main
//	 go get -x github.com/ddkwork/golibrary@master
//	 go get -x github.com/ddkwork/ux@master
//	 go get -x github.com/google/go-cmp@master
//	 go get -x github.com/ddkwork/app@master
//	 go get -x github.com/ddkwork/toolbox@master
//	 go get -x github.com/ddkwork/unison@master
//	 go get -x github.com/ebitengine/purego@main
//	 go get -x github.com/saferwall/pe@main
//	 ::go get -u -x all
//	 go mod tidy
//
//	go install mvdan.cc/gofumpt@latest
//	gofumpt -l -w .
//	//go install honnef.co/go/tools/cmd/staticcheck@latest
//	//staticcheck ./...
//	//go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
//
//`) {
//		s = strings.TrimSpace(s)
//		if strings.HasPrefix(s, "::") || strings.HasPrefix(s, "//") || s == "" {
//			continue
//		}
//		RunCommand(s)
//	}
//}

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
