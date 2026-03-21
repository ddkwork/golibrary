package mylog

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	logrusPackage      string
	minimumCallerDepth = 1
	callerInitOnce     sync.Once
	modulePath         string
	moduleRoot         string
	moduleInitOnce     sync.Once
)

const (
	maximumCallerDepth int = 25
	knownLogrusFrames  int = 4
)

func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}
	return f
}

func initModuleInfo() {
	moduleInitOnce.Do(func() {
		if data, err := os.ReadFile("go.mod"); err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "module ") {
					modulePath = strings.TrimSpace(strings.TrimPrefix(line, "module "))
					break
				}
			}
		}
		if modulePath != "" {
			moduleRoot = findModuleRoot()
		}
	})
}

func findModuleRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

func shortenPath(fullPath string) string {
	initModuleInfo()
	if moduleRoot == "" {
		return filepath.Base(fullPath)
	}
	relPath, err := filepath.Rel(moduleRoot, fullPath)
	if err != nil {
		return filepath.Base(fullPath)
	}
	return filepath.ToSlash(relPath)
}

func shortenFunction(fullFunc string) string {
	initModuleInfo()
	if modulePath == "" {
		return filepath.Base(fullFunc)
	}
	if idx := strings.Index(fullFunc, modulePath); idx != -1 {
		return fullFunc[idx+len(modulePath)+1:]
	}
	return fullFunc
}

func caller() string {
	frame := getCaller()
	CheckNil(frame)
	return fmt.Sprintf("%s %s:%d", shortenFunction(frame.Function), shortenPath(frame.File), frame.Line)
}

func getCaller() *runtime.Frame {
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, maximumCallerDepth)
		runtime.Callers(0, pcs)
		for i := range maximumCallerDepth {
			funcName := runtime.FuncForPC(pcs[i]).Name()
			if strings.Contains(funcName, "getCaller") {
				logrusPackage = getPackageName(funcName)
				break
			}
		}
		minimumCallerDepth = knownLogrusFrames
	})
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)
		if pkg != logrusPackage {
			return &f
		}
	}
	return nil
}
