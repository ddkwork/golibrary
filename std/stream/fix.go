package stream

import (
	"os"
	"os/exec"
)

func Fix(dir string) {
	// 设置 GOOS=windows 跳过 Linux 特定的 cgo 依赖
	os.Setenv("GOOS", "windows")
	os.Setenv("CGO_ENABLED", "0")
	// 执行 go fix，忽略错误（某些包可能因缺少 cgo 依赖无法解析）
	cmd := exec.Command("go", "fix", "./...")
	cmd.Dir = dir
	cmd.Run() // 忽略错误
}

func Fmt(dir string) {
	// stream.RunCommandWithDir(repoDir, "go", "run", "mvdan.cc/gofumpt@latest", "-l", "-w", ".")
	RunCommandWithDir(dir, "go", "run", "mvdan.cc/gofumpt@latest", "-w", ".")
}
