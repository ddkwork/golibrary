package environment

import (
	"github.com/ddkwork/golibrary/mylog"
	"golang.org/x/sys/windows/registry"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

// 功能就是刷新windows所有版本的系统path环境变量
// [HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment]
type (
	Interface interface {
		//canvasobjectapi.Interface
		WalkDirs(roots ...string) (ok bool) //遍历需要加入到path的文件夹集生成目录切片,map+goruntine
		Orig() (ok bool)                    //读取系统Environment并用map去重+delete bad path
		Update() (ok bool)                  //sort by strings
	}
	object struct {
		paths    []string
		pathsMap map[string]string
		key      registry.Key
	}
)

// ComSpec=%SystemRoot%\system32\cmd.exe
var (
	name    = "Path"
	EnvPath = `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`
	skip    = []string{
		"genx",
		"steam",
		"todo",
		"vmware",
		"ndk",
		"媒体",
		"apk",
		"vs2022",
		"clone",
	}
	defalutRoots = []string{
		"D:\\bin",
		//"D:\\codespace",
	}
)

func (o *object) Orig() (ok bool) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, EnvPath, registry.ALL_ACCESS)
	if !mylog.Error(err) {
		return
	}
	o.key = key
	value, _, err := key.GetStringValue(name)
	if !mylog.Error(err) {
		return
	}
	split := strings.Split(value, ";")
	o.paths = append(o.paths, split...)
	return true
}

func (o *object) WalkDirs(roots ...string) (ok bool) {
	if len(roots) == 0 {
		roots = defalutRoots
	}
	for _, root := range roots {
		filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
			for _, s := range skip {
				if strings.Contains(path, s) {
					return nil
				}
			}
			ext := filepath.Ext(path)
			switch ext {
			case ".exe", ".cmd", ".bat":
				dir := filepath.Dir(path)
				o.pathsMap[dir] = dir
			}
			return err
		})
	}
	for _, s2 := range o.pathsMap {
		o.paths = append(o.paths, s2)
	}
	return true
}

func (o *object) Update() (ok bool) {
	sort.Strings(o.paths)
	for _, path := range o.paths {
		mylog.Info("update", path)
	}
	return mylog.Error(o.key.SetStringValue(name, strings.Join(o.paths, ";")))
}

func New() Interface {
	return &object{
		paths:    make([]string, 0),
		pathsMap: make(map[string]string),
	}
}
