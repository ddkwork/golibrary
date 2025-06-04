package stream

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/waitgroup"
	"golang.org/x/text/encoding/simplifiedchinese"
)

type CommandSession struct {
	cmdKey  string
	command string
	Path    string
	isClang bool
	Args    []string
	Env     []string
	Dir     string
	Stdin   *Buffer
	Stdout  *Buffer
	Stderr  *Buffer
}

func newCommandSession() *CommandSession {
	return &CommandSession{
		cmdKey:  "command", //isClang the key is clang target file path
		command: "",
		Path:    "",
		isClang: false,
		Args:    nil,
		Env:     nil,
		Dir:     mylog.Check2(os.Getwd()),
		Stdin:   NewBuffer(""),
		Stdout:  NewBuffer(""),
		Stderr:  NewBuffer(""),
	}
}

func RunCommandArgs(arg ...string) *CommandSession {
	s := newCommandSession()
	s.command = strings.Join(arg, " ")
	if arg[0] == "clang" {
		s.cmdKey = filepath.Base(arg[len(arg)-1])
		s.Path = s.cmdKey
		s.isClang = true
	}
	return s.run()
}

func RunCommand(command string) *CommandSession {
	s := newCommandSession()
	s.command = command
	return s.run()
}

func RunCommandWithDir(command, dir string) *CommandSession {
	s := newCommandSession()
	s.command = command
	s.Dir = dir
	return s.run()
}

func (s *CommandSession) run() *CommandSession {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fnInitCmd := func() *exec.Cmd {
		if runtime.GOOS == "windows" {
			return exec.Command("cmd", "/C", s.command)
		}
		return exec.Command("bash", "-c", s.command)
	}
	cmd := fnInitCmd()
	cmd.Dir = s.Dir

	mylog.Info(s.cmdKey, s.command)

	stdoutPipe := mylog.Check2(cmd.StdoutPipe())
	stderrPipe := mylog.Check2(cmd.StderrPipe())

	mylog.Check(cmd.Start())

	g := waitgroup.New()
	g.UseMutex = false
	g.SetLimit(1000)
	output := make(chan string)
	errorOutput := make(chan string)

	// 启动 goroutine 读取 stdout
	g.Go(func() {
		mylog.Call(func() {
			if s.isClang {
				s.Stdout.ReadFrom(stdoutPipe)
				return
			}
			scanner := bufio.NewScanner(stdoutPipe)
			for scanner.Scan() {
				output <- ConvertUtf82Gbk(s.isClang,scanner.Text())
			}
			mylog.Check(stdoutPipe.Close())
		})
	})

	// 启动 goroutine 读取 stderr
	g.Go(func() {
		mylog.Call(func() {
			if s.isClang {
				s.Stderr.ReadFrom(stderrPipe)
				return
			}
			scanner := bufio.NewScanner(stderrPipe)
			for scanner.Scan() {
				errorOutput <- ConvertUtf82Gbk(s.isClang,scanner.Text())
			}
			mylog.Check(stderrPipe.Close())
		})
	})

	// 启动 goroutine 统一处理输出
	go func() {
		mylog.Call(func() {
			g.Wait()
			close(output)
		})
	}()

	done := make(chan struct{})
	go func() {
		for line := range output {
			if !s.isClang {
				println(line) //对于json，不需要每一行都输出，而是一次性返回解码或者落地保存
			}
			s.Stdout.WriteStringLn(line)
		}
		done <- struct{}{}
	}()

	go func() {
		for line := range errorOutput {
			//println(line)//让clang dump ast的错误在后面写入文件，这里和后面调用s.Error.String()重复输出错误了
			s.Stderr.WriteStringLn(line)
		}
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		mylog.Check(cmd.Process.Kill())
	case <-done:
	}

	e := cmd.Wait()
	if e != nil {
		mylog.Check(ConvertUtf82Gbk(s.isClang,e.Error() + "\n" + s.Stderr.String()))
	}
	if s.isClang {
		return s
	}
	ss := trimTrailingEmptyLines(s.Stdout.String())
	s.Stdout.Reset()
	s.Stdout.WriteString(ss)
	return s
}

func trimTrailingEmptyLines(s string) string {
	// 使用正则表达式匹配末尾的所有空白行，包括空格、制表符和换行符
	re := regexp.MustCompile(`\s*\n*$`)
	return re.ReplaceAllString(s, "")
}

func ConvertUtf82Gbk(isCalng bool,src string) string {
	if isCalng {
		return src
	}
	if IsWindows() {
		c := mylog.Check2(simplifiedchinese.GB18030.NewDecoder().String(src)) // todo test rune
		return strings.TrimSuffix(c, "\r\n")
	}
	return src
}

