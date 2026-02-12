package stream

import (
	"os"
)

func Fix(dir string) {
	// 保存原始值
	oldGOOS := os.Getenv("GOOS")

	// 设置 GOOS=windows 跳过 Linux 特定的 cgo 依赖
	os.Setenv("GOOS", "windows")

	RunCommandWithDir(dir, "go fix ./...")

	// 恢复原始值
	if oldGOOS == "" {
		os.Unsetenv("GOOS")
	} else {
		os.Setenv("GOOS", oldGOOS)
	}
}

func Fmt(dir string) {
	// stream.RunCommandWithDir(repoDir, "go", "run", "mvdan.cc/gofumpt@latest", "-l", "-w", ".")
	RunCommandWithDir(dir, "go", "run", "mvdan.cc/gofumpt@latest", "-w", ".")
}
