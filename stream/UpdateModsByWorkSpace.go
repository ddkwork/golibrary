package stream

import (
	"archive/zip"
	"github.com/ddkwork/golibrary/waitgroup"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/safemap"
	"golang.org/x/mod/modfile"
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
// func GetLastCommitHash(repositoryName string) string {
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
// }

func GetLastCommitHashLocal(repositoryDir string) string { // 如果失败了，发现禁用模块代理可以成功，那么需要再提交点别的，然后模块代理就会识别新的提交hash，很诡异
	originPath := mylog.Check2(os.Getwd())
	mylog.Check(os.Chdir(repositoryDir))
	hash := RunCommand("git rev-parse HEAD").Stdout.String()
	mylog.Check(os.Chdir(originPath))
	return hash
}

func ParseGoMod(file string, data []byte) *safemap.M[string, string] {
	f := mylog.Check2(modfile.Parse(file, data, nil))
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

func UpdateDependenciesFromModFile(dir string) { // 实现替换，不要网络访问了，太慢了
	if filepath.Base(dir) == "golibrary" {
		return
	}
	r := mylog.Check2(zip.OpenReader("D:\\workspace\\workspace\\mod.zip"))
	newBody := mylog.Check2(io.ReadAll(mylog.Check2(r.Open("go.mod"))))
	originMod := filepath.Join(dir, "go.mod")
	body := mylog.Check2(os.ReadFile(originMod))
	f := mylog.Check2(modfile.Parse(originMod, body, nil))
	newModMap := ParseGoMod("new.mod", newBody)
	for oldName, oldVersion := range ParseGoMod(originMod, body).Range() {
		newVersion, exist := newModMap.Get(oldName)
		if exist {
			if oldVersion != newVersion {
				for i, require := range f.Require {
					if require.Mod.Path == oldName {
						require.Mod.Version = newVersion
						// f.Require[i] = require
						setVersion(f.Require[i], newVersion)
					}
				}
			}
		}
	}
	f.Cleanup()
	f.SortBlocks()
	updateModFile := mylog.Check2(f.Format())
	// println(string(updateModFile))
	WriteTruncate(originMod, updateModFile)
	g := waitgroup.New()
	g.Go(func() {
		RunCommandWithDir("go mod tidy", dir)
		v := newModMap.GetMust("github.com/ddkwork/golibrary")
		b := NewBuffer(originMod)
		if !b.Contains("golibrary") {
			line := "require github.com/ddkwork/golibrary  " + v
			mylog.Info("add golibrary", line)
			b.WriteStringLn(line).ReWriteSelf()
		}
		// https://github.com/ddkwork/tools/blob/master/gopls/doc/analyzers.md
		RunCommandWithDir("go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -diff ./...", dir)
		RunCommandWithDir("go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -fix ./...", dir)
	})
	g.Wait()
}

func setVersion(r *modfile.Require, v string) {
	r.Mod.Version = v

	if line := r.Syntax; len(line.Token) > 0 {
		if line.InBlock {
			// If the line is preceded by an empty line, remove it; see
			// https://golang.org/issue/33779.
			if len(line.Comments.Before) == 1 && len(line.Comments.Before[0].Token) == 0 {
				line.Comments.Before = line.Comments.Before[:0]
			}
			if len(line.Token) >= 2 { // example.com v1.2.3
				line.Token[1] = v
			}
		} else {
			if len(line.Token) >= 3 { // require example.com v1.2.3
				line.Token[2] = v
			}
		}
	}
}

func UpdateDependencies(path string) { // 模块代理刷新的不及时，需要禁用代理,已经使用clone仓库远程完成更新
	var mutex sync.Mutex
	g := waitgroup.New()
	for s := range ReadFileToLines(filepath.Join(GetDesktopDir(), "dep.txt")) { // 因为要经常更新，我们不embed
		s = strings.TrimSpace(s)
		if strings.HasPrefix(s, "::") || strings.HasPrefix(s, "//") || s == "" {
			continue
		}
		if s == "go mod tidy" {
			continue
		}
		g.Go(func() { // 这样之后tidy就不在最后执行了，同时升级多个依赖+读写锁定
			RunCommandWithDir(s, path)
		})
		g.Wait()
	}
	mutex.Lock()
	RunCommandWithDir("go mod tidy", path) // 所有yield都执行完了，再执行tidy
	mutex.Unlock()
}

func updateModsByWorkSpace(isUpdateAll bool) {
	if !FileExists("go.work") {
		UpdateDependenciesFromModFile(mylog.Check2(os.Getwd()))
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
			if filepath.Base(abs) == "golibrary" {
				continue
			}
			mods = append(mods, abs)
		}
	}

	modChan := make(chan string, len(mods))

	g := waitgroup.New()
	for _, modPath := range mods {
		g.Go(func() { // 每个模块单独跑,这里不能加锁，否则很慢，谨慎使用读写锁
			UpdateDependenciesFromModFile(modPath) // 锁应该在这里面
			if isUpdateAll {
				RunCommand("go get -u -x all") // need lock,但是不使用这个，太慢了
			}
			modChan <- modPath
		})
	}
	go func() {
		for mod := range modChan {
			mylog.Success("updated mod", strconv.Quote(mod))
		}
		close(modChan)
	}()
	g.Wait()
	mylog.Success("all work finished")
}

// func updateDependencies() { // 模块代理刷新的不及时，需要禁用代理,已经使用clone仓库远程完成更新
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
// `) {
//		s = strings.TrimSpace(s)
//		if strings.HasPrefix(s, "::") || strings.HasPrefix(s, "//") || s == "" {
//			continue
//		}
//		RunCommand(s)
//	}
// }
//
// type Cache struct {
// 	store map[string]string
// 	mu    sync.RWMutex
// }
//
// func TestCacheConsistency(t *testing.T) {
// 	cache := &Cache{store: make(map[string]string)}
//
// 	synctest.Run(t, func(sc *synctest.Scenario) {
// 		// 并发写入
// 		for i := 0; i < 10; i++ {
// 			sc.Go(func() {
// 				cache.mu.Lock()
// 				defer cache.mu.Unlock()
// 				cache.store["key"] = time.Now().String()
// 			})
// 		}
//
// 		// 并发读取
// 		for i := 0; i < 100; i++ {
// 			sc.Go(func() {
// 				cache.mu.RLock()
// 				defer cache.mu.RUnlock()
// 				_ = cache.store["key"]
// 			})
// 		}
// 	}, synctest.WithOptions(
// 		synctest.EnableRaceDetection(),
// 		synctest.MaxGoroutines(200),
// 		synctest.Timeout(10*time.Second),
// 	))
// }
