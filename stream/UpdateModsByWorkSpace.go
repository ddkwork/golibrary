package stream

import (
	"github.com/ddkwork/golibrary/safemap"
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

// GetLastCommitHash
// git ls-remote https://github.com/ddkwork/toolbox refs/heads/master
// git ls-remote https://github.com/gioui/gio refs/heads/main
// 但是要传递本地的仓库目录，太麻烦了
func GetLastCommitHash(repositoryName string) string {
	//理论上获取到hash再使用模块代理，这样才刷新的快？
	//或者使用action得到hash先？

	mylog.Check(os.Setenv("GOPROXY", "direct")) // 如果模块代理导致获取到的不是最新的提交哈希那么需要禁用模块代理，最可靠的方式是 GetLastCommitHashLocal
	defer mylog.Check(os.Setenv("GOPROXY", "https://goproxy.cn,direct"))

	//defer RunCommand("go env -w GOPROXY=https://goproxy.cn,direct")
	s := RunCommand("git ls-remote " + repositoryName + " refs/heads/master").Output.String()
	for hash := range strings.FieldsSeq(s) {
		return hash
	}
	s = RunCommand("git ls-remote " + repositoryName + " refs/heads/main").Output.String()
	for hash := range strings.FieldsSeq(s) {
		return hash
	}
	panic("no commit hash found in master and main branch")
	/*
		# 获取最新提交哈希
		$hash = (git ls-remote https://github.com/ddkwork/toolbox refs/heads/master).Split("`t")[0]
		# 带哈希安装
		go get -x "github.com/ddkwork/toolbox@$hash"

		powershell:
			$env:GOPROXY="direct"; go get -x github.com/ddkwork/toolbox@$(git ls-remote https://github.com/ddkwork/toolbox refs/heads/master | ForEach-Object { $_.Split()[0] })

	*/
}

type ModInfo struct {
	ModName       string //gioui.org
	RepositoryUrl string //https://github.com/gioui/gio
	repositoryDir string //GetLastCommitHashLocal
	Hash          string //GetLastCommitHash or GetLastCommitHashLocal
	UpdateCommand string //go get -x gioui.org@hash
}

var m = safemap.NewOrdered[string, ModInfo](func(yield func(string, ModInfo) bool) {
	yield("https://github.com/gioui/gio", ModInfo{
		ModName:       "gioui.org",
		RepositoryUrl: "https://github.com/gioui/gio",
		repositoryDir: "",
		Hash:          "",
		UpdateCommand: "",
	})
	yield("https://github.com/gioui/gio/cmd", ModInfo{
		ModName:       "gioui.org/cmd",
		RepositoryUrl: "https://github.com/gioui/gio/cmd",
		repositoryDir: "",
		Hash:          "",
		UpdateCommand: "",
	})
	yield("https://github.com/gioui/gio/example", ModInfo{
		ModName:       "gioui.org/example",
		RepositoryUrl: "https://github.com/gioui/gio/example",
		repositoryDir: "",
		Hash:          "",
		UpdateCommand: "",
	})
	yield("https://github.com/gioui/gio/x", ModInfo{
		ModName:       "gioui.org/x",
		RepositoryUrl: "https://github.com/gioui/gio/x",
		repositoryDir: "",
		Hash:          "",
		UpdateCommand: "",
	})
	yield("https://github.com/oligo/gvcode", ModInfo{
		ModName:       "github.com/oligo/gvcode",
		RepositoryUrl: "https://github.com/oligo/gvcode",
		repositoryDir: "",
		Hash:          "",
		UpdateCommand: "",
	})
	yield("https://github.com/ddkwork/golibrary", ModInfo{
		ModName:       "github.com/ddkwork/golibrary",
		RepositoryUrl: "https://github.com/ddkwork/golibrary",
		repositoryDir: "",
		Hash:          "",
		UpdateCommand: "",
	})
	yield("https://github.com/ddkwork/ux", ModInfo{
		ModName:       "github.com/ddkwork/ux",
		RepositoryUrl: "https://github.com/ddkwork/ux",
		repositoryDir: "",
		Hash:          "",
		UpdateCommand: "",
	})
	yield("https://github.com/google/go-cmp", ModInfo{
		ModName:       "github.com/google/go-cmp",
		RepositoryUrl: "https://github.com/google/go-cmp",
		repositoryDir: "",
		Hash:          "",
		UpdateCommand: "",
	})
	yield("https://github.com/ddkwork/app", ModInfo{
		ModName:       "github.com/ddkwork/app",
		RepositoryUrl: "https://github.com/ddkwork/app",
		repositoryDir: "",
		Hash:          "",
		UpdateCommand: "",
	})
	yield("https://github.com/ddkwork/toolbox", ModInfo{
		ModName:       "github.com/ddkwork/toolbox",
		RepositoryUrl: "https://github.com/ddkwork/toolbox",
		repositoryDir: "",
		Hash:          "",
		UpdateCommand: "",
	})
	yield("https://github.com/ddkwork/unison", ModInfo{
		ModName:       "github.com/ddkwork/unison",
		RepositoryUrl: "https://github.com/ddkwork/unison",
		repositoryDir: "",
		Hash:          "",
		UpdateCommand: "",
	})
	yield("https://github.com/ebitengine/purego", ModInfo{
		ModName:       "github.com/ebitengine/purego",
		RepositoryUrl: "https://github.com/ebitengine/purego",
		repositoryDir: "",
		Hash:          "",
		UpdateCommand: "",
	})
	yield("https://github.com/saferwall/pe", ModInfo{
		ModName:       "github.com/saferwall/pe",
		RepositoryUrl: "https://github.com/saferwall/pe",
		repositoryDir: "",
		Hash:          "",
		UpdateCommand: "",
	})

})

func GetLastCommitHashLocal(repositoryName, repositoryDir string) string {
	originPath := mylog.Check2(os.Getwd())
	mylog.Check(os.Chdir(repositoryDir))
	hash := RunCommand("git rev-parse HEAD").Output.String()
	mylog.Check(os.Chdir(originPath))
	return hash
}

func UpdateDependencies() { //模块代理刷新的不及时，需要禁用代理
	mylog.Check(os.Setenv("GOPROXY", "direct"))
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
	//go install honnef.co/go/tools/cmd/staticcheck@latest
	//staticcheck ./...
	//go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

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
				RunCommand("go get -u -x all")
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
