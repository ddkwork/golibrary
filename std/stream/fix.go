package stream

import (
	"os"
	"path/filepath"
)

func fix(dir string, args ...string) {
	oldGOOS := os.Getenv("GOOS")
	os.Setenv("GOOS", "windows")
	RunCommandWithDir(dir, args...)
	if oldGOOS == "" {
		os.Unsetenv("GOOS")
	} else {
		os.Setenv("GOOS", oldGOOS)
	}
}

func FixDir(dir string) { fix(dir, "go", "fix", "./...") }

func FixFile(file string) { fix(filepath.Dir(file), "go", "fix", file) }

func FmtDir(dir string) {
	RunCommandWithDir(dir, "go", "run", "mvdan.cc/gofumpt@latest", "-l", "-w", ".")
}

func FmtFile(file string) {
	RunCommandWithDir(filepath.Dir(file), "go", "run", "mvdan.cc/gofumpt@latest", "-l", "-w", file)
}

func Fix(path string) {
	if IsDir(path) {
		FixDir(path)
	} else {
		FixFile(path)
	}
}

func Fmt(path string) {
	if IsDir(path) {
		FmtDir(path)
	} else {
		FmtFile(path)
	}
}
