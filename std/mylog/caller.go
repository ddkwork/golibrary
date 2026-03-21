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

var pcsPool = NewPool(func() []uintptr {
	return make([]uintptr, maximumCallerDepth)
})

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
		data, e := os.ReadFile("go.mod")
		if e == nil {
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
	dir, e := os.Getwd()
	if e != nil {
		return ""
	}

	for {
		_, e := os.Stat(filepath.Join(dir, "go.mod"))
		if e == nil {
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
	relPath, e := filepath.Rel(moduleRoot, fullPath)
	if e != nil {
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
	return fmt.Sprintf("%s:%d", shortenPath(frame.File), frame.Line)
}

func callerFuncName() string {
	frame := getCaller()
	CheckNil(frame)
	funcName := shortenFunction(frame.Function)
	if idx := strings.LastIndex(funcName, "."); idx != -1 {
		funcName = funcName[idx+1:]
	}
	if len(funcName) > keyLen {
		return funcName[:keyLen-3] + "..."
	}
	return funcName
}

func getCaller() *runtime.Frame {
	callerInitOnce.Do(func() {
		pcs := pcsPool.Get()
		defer pcsPool.Put(pcs)
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
	pcs := pcsPool.Get()
	defer pcsPool.Put(pcs)
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
