package stream

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/ddkwork/golibrary/mylog"
)

func CurrentDirName(path string) (currentDirName string) {
	if path == "" {
		getwd, err := os.Getwd()
		if !mylog.Error(err) {
			return
		}
		path = getwd
	}
	split := strings.Split(path, "\\")
	if split == nil {
		return BaseName(filepath.Dir(path)) // todo test
	}
	return split[len(split)-1]
}

func CopyDir(source, destination string) (ok bool) {
	base := filepath.Base(source)
	return mylog.Error(filepath.WalkDir(source, func(p string, info fs.DirEntry, err error) error {
		split := strings.Split(p, base)
		dst := filepath.Join(destination, base, split[1])
		switch {
		case info.IsDir():
			if !CreatDirectory(dst) {
				return err
			}
		default:
			if !ReadFileAndWriteTruncate(p, dst) {
				return err
			}
		}
		return err
	}))
}

func ReadFileAndWriteTruncate(path, dstPath string) (ok bool) {
	return WriteTruncate(dstPath, NewReadFile(path).Bytes())
}

func FileExists(path string) bool {
	fi, err := os.Stat(path)
	if !mylog.Error(err) {
		return false
	}
	if fi == nil {
		return false
	}
	mode := fi.Mode()
	return !mode.IsDir() && mode.IsRegular()
}

func FixUncPath(path string) string {
	all := strings.ReplaceAll(path, "\\", "/")
	return strings.ReplaceAll(all, "//", "/")
}

// IsDirDeep1 是否是一级目录，包含目录分隔符的深度一定大于1
func IsDirDeep1(path string) bool { return !strings.Contains(FixUncPath(path), "/") }

func IsDir(path string) bool {
	if strings.HasPrefix(path, ".") && IsDirDeep1(path) { //.git .github .vs
		return true // 一般这种目录太深，遍历浪费时间
	}
	fi, err := os.Stat(path)
	if fi == nil {
		return false
	}
	return err == nil && fi.IsDir()
}

func BaseName(path string) string {
	abs, err := filepath.Abs(path)
	if !mylog.Error(err) {
		return ""
	}
	return TrimExtension(filepath.Base(abs))
}
func TrimExtension(path string) string { return path[:len(path)-len(filepath.Ext(path))] }

func JoinHomeDir(path string) (join string, ok bool)  { return joinHome(path, true) }
func JoinHomeFile(path string) (join string, ok bool) { return joinHome(path, false) }
func joinHome(path string, isDir bool) (join string, ok bool) {
	join = filepath.Join(HomeDir(), path)
	if !FileExists(join) {
		switch isDir {
		case true:
			if !mylog.Error(os.MkdirAll(join, 0o777)) {
				return
			}
		default:
			f, err := os.Create(join)
			if !mylog.Error(err) {
				return
			}
			if !mylog.Error(f.Close()) {
				return
			}
		}
	}
	ok = true
	return
}

func HomeDir() string {
	if u, err := user.Current(); err == nil {
		return u.HomeDir
	}
	if dir, err := os.UserHomeDir(); err == nil {
		return dir
	}
	return "."
}

func MoveFile(src, dst string) (err error) {
	var srcInfo, dstInfo os.FileInfo
	srcInfo, err = os.Stat(src)
	if err != nil {
		return errors.Unwrap(err)
	}
	if !srcInfo.Mode().IsRegular() {
		return errors.New(fmt.Sprintf("%s is not a regular file", src))
	}
	dstInfo, err = os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return errors.Unwrap(err)
		}
	} else {
		if !dstInfo.Mode().IsRegular() {
			return errors.New(fmt.Sprintf("%s is not a regular file", dst))
		}
		if os.SameFile(srcInfo, dstInfo) {
			return nil
		}
	}
	if os.Rename(src, dst) == nil {
		return nil
	}
	var in, out *os.File
	out, err = os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return errors.Unwrap(err)
	}
	defer func() {
		if closeErr := out.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()
	if in, err = os.Open(src); err != nil {
		err = errors.Unwrap(err)
		return
	}
	_, err = io.Copy(out, in)
	if !mylog.Error(in.Close()) {
		return
	}
	if err = os.Remove(src); err != nil {
		err = errors.Unwrap(err)
	}
	return
}

func RunDirAbs() string {
	abs, err := filepath.Abs(RunDir())
	if !mylog.Error(err) {
		return ""
	}
	return abs
}

func RunDir() string {
	dir, err := os.Getwd()
	if !mylog.Error(err) {
		return ""
	}
	return filepath.Base(dir)
}
