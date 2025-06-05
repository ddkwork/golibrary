package i18n

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/golibrary/std/stream"
)

const (
	Extension = ".i18n"
)

var (
	Dir      string
	Language = Locale()

	Languages    = strings.Split(os.Getenv("LANGUAGE"), ":")
	altLocalizer atomic.Pointer[localizer]
	once         sync.Once
	langMap      = make(map[string]map[string]string)
	hierLock     sync.Mutex
	hierMap      = make(map[string][]string)
)

type localizer struct {
	Text func(string) string
}

func SetLocalizer(f func(string) string) {
	var trampoline *localizer
	if f != nil {
		trampoline = &localizer{Text: f}
	}
	altLocalizer.Store(trampoline)
}

func Text(text string) string {
	if f := altLocalizer.Load(); f != nil {
		return f.Text(text)
	}
	once.Do(func() {
		if Dir == "" {
			path := mylog.Check2(os.Executable())

			path = mylog.Check2(filepath.EvalSymlinks(path))

			path = stream.BaseName(path) + "_i18n"

			Dir = path
		}
		dirEntry := mylog.Check2(os.ReadDir(Dir))

		for _, one := range dirEntry {
			if !one.IsDir() {
				name := one.Name()
				if filepath.Ext(name) == Extension {
					load(name)
				}
			}
		}
	})

	var result string
	if result = lookup(text, Language); result != "" {
		return result
	}
	for _, language := range Languages {
		if result = lookup(text, language); result != "" {
			return result
		}
	}
	return text
}

func lookup(text, language string) string {
	for _, lang := range hierarchy(language) {
		if translations := langMap[lang]; translations != nil {
			if str, ok := translations[text]; ok {
				return str
			}
		}
	}
	return ""
}

func hierarchy(language string) []string {
	lang := strings.ToLower(language)
	hierLock.Lock()
	defer hierLock.Unlock()
	if s, ok := hierMap[lang]; ok {
		return s
	}
	one := strings.ReplaceAll(strings.ReplaceAll(lang, "-", "_"), ".", "_")
	var s []string
	for {
		s = append(s, one)
		if i := strings.LastIndex(one, "_"); i != -1 {
			one = one[:i]
		} else {
			break
		}
	}
	hierMap[lang] = s
	return s
}

func load(name string) {
	path := filepath.Join(Dir, name)
	f := mylog.Check2(os.Open(path))

	defer f.Close()
	lineNum := 1
	lastKeyLineStart := 1
	translations := make(map[string]string)
	var key, value string
	var hasKey, hasValue bool
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "k:") {
			if hasValue {
				if _, exists := translations[key]; !exists {
					translations[key] = value
				} else {
					slog.Warn("i18n: ignoring duplicate key", "line", lastKeyLineStart, "file", path)
				}
				hasKey = false
				hasValue = false
			}
			var buffer string
			mylog.Check2(fmt.Sscanf(line, "k:%q", &buffer))

		} else if strings.HasPrefix(line, "v:") {
			if hasKey {
				var buffer string
				mylog.Check2(fmt.Sscanf(line, "v:%q", &buffer))

			} else {
				slog.Warn("i18n: ignoring value with no previous key", "line", lineNum, "file", path)
			}
		}
		lineNum++
	}
	if hasKey {
		if hasValue {
			if _, exists := translations[key]; !exists {
				translations[key] = value
			} else {
				slog.Warn("i18n: ignoring duplicate key", "line", lastKeyLineStart, "file", path)
			}
		} else {
			slog.Warn("i18n: ignoring key with missing value", "line", lastKeyLineStart, "file", path)
		}
	}
	key = strings.ToLower(name[:len(name)-len(Extension)])
	key = strings.ReplaceAll(strings.ReplaceAll(key, "-", "_"), ".", "_")
	langMap[key] = translations
}
