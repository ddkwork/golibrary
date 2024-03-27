package cmd

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
	"golang.org/x/text/encoding/simplifiedchinese"
)

type Session struct {
	pid       chan int
	logWriter io.Writer
	*exec.Cmd
	Result           string
	CurrentDirectory string
}

func MakeArg(arg ...string) string {
	args := make([]string, 0)
	args = append(args, arg...)
	return strings.Join(args, " ")
}
func RunArgs(arg ...string) *Session { return Run(MakeArg(arg...)) }
func Run(command string) *Session {
	session := NewSession()
	session.Result = session.run(context.Background(), command)
	return session
}

func NewSession() *Session {
	return &Session{
		pid:              make(chan int, 1),
		logWriter:        nil,
		Cmd:              nil,
		Result:           "",
		CurrentDirectory: "",
	}
}

func (s *Session) SetDir(dir string)   { s.Dir = strings.TrimSpace(dir) }
func (s *Session) SetLog(wr io.Writer) { s.logWriter = wr }
func (s *Session) GetPid() <-chan int  { return s.pid }
func (s *Session) run(ctx context.Context, command string) (ss string) {
	fnInitCmd := func() *exec.Cmd {
		if runtime.GOOS == "windows" {
			return exec.Command("cmd", "/C", command)
		}
		return exec.Command("bash", "-c", command) //"linux", "darwin", "freebsd":
	}
	cmd := fnInitCmd()
	s.Cmd = cmd
	dir, err := os.Getwd()
	if !mylog.Error(err) {
		return
	}
	s.CurrentDirectory = dir

	mylog.Warning("go-command", cmd.String())
	outputErr := &bytes.Buffer{}
	outputOut := &bytes.Buffer{}

	cmd.Stderr = io.MultiWriter(outputErr, os.Stderr)
	pipe, err := cmd.StdoutPipe()
	if !mylog.Error(err) {
		return err.Error()
	}
	go func() {
		reader := bufio.NewReader(pipe)
		for {
			line, _, err := reader.ReadLine()
			if err != nil || errors.Is(err, io.EOF) {
				break
			}
			ss = ConvertUtf82Gbk(line)
			mylog.Warning("line", ss)
			outputOut.WriteString(ss)
		}
	}()
	if !mylog.Error(cmd.Start()) {
		return
	}
	if cmd.Process == nil {
		return
	}
	s.pid <- cmd.Process.Pid

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	done := make(chan struct{}, 1)
	go func() {
		select {
		case <-ctx.Done():
			mylog.Error(ConvertUtf82Gbk([]byte((ctx.Err().Error()))))
			mylog.Error(cmd.Process.Kill())
		case <-done:
		}
	}()

	err = cmd.Wait()
	done <- struct{}{}
	if err != nil {
		ss = ConvertUtf82Gbk((outputErr).Bytes())
		mylog.Error(ss)
		return
	}
	return ConvertUtf82Gbk((outputOut).Bytes())
}

func ConvertUtf82Gbk(src []byte) (ss string) {
	defer func() {
		if ss != "" {
			ss = strings.TrimRight(ss, "\r\n")
			ss = stream.Utf82Gbk(ss)
		}
	}()
	if stream.IsWindows() { // windows控制台和vs2022都是gbk编码的，包括他编写的代码
		return stream.Utf82Gbk(string(src))
		s, err := simplifiedchinese.GB18030.NewDecoder().Bytes(src)
		if !mylog.Error(err) {
			return err.Error()
		}
		return string(s)
	}
	return string(src)
}
