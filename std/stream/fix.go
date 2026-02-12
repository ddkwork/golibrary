package stream

import "os"

func Fix(dir string) {
	// 禁用 cgo 执行避免缺少 Linux 系统依赖
	os.Setenv("CGO_ENABLED", "0")
	// 忽略无法解析的包错误
	RunCommandWithDir(dir, "bash", "-c", "go fix ./... 2>/dev/null || true")
}

func Fmt(dir string) {
	// stream.RunCommandWithDir(repoDir, "go", "run", "mvdan.cc/gofumpt@latest", "-l", "-w", ".")
	RunCommandWithDir(dir, "go", "run", "mvdan.cc/gofumpt@latest", "-w", ".")
}
