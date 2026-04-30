package cmake

import (
	"path/filepath"

	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/golibrary/std/stream"
)

const BinDir = "C:/Program Files/CMake/bin"

var (
	EwdkCmakeFile = filepath.Join(BinDir, "ewdk.cmake")
	EwdkEnvFile   = filepath.Join(BinDir, "ewdk.env.json")
)

func Module() string {
	matches, e := filepath.Glob(filepath.Join(filepath.Join(filepath.Dir(BinDir), "share"), "cmake-*"))
	if e != nil || len(matches) == 0 {
		mylog.Warning("未找到 cmake 模块目录")
		mylog.Check(e)
	}
	globalModuleDir := filepath.Join(matches[0], "Modules")
	if !stream.IsDir(globalModuleDir) {
		panic("模块路径不存在: " + globalModuleDir)
	}
	return globalModuleDir
}
