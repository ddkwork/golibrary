package stream

import (
	"bufio"
	"context"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/ddkwork/golibrary/mylog"
	"golang.org/x/text/encoding/simplifiedchinese"
)

type CommandSession struct {
	Output           *Buffer
	Error            *Buffer
	CurrentDirectory string
}

func RunCommandArgs(arg ...string) *CommandSession { return RunCommand(strings.Join(arg, " ")) }
func RunCommand(command string) (session *CommandSession) {
	mylog.Call(func() {
		session = &CommandSession{
			Output:           NewBuffer(""),
			Error:            NewBuffer(""),
			CurrentDirectory: mylog.Check2(os.Getwd()),
		}
		session.run(command)
	})
	return
}

func (s *CommandSession) run(command string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	skipLog := false
	if strings.Contains(command, "clang") {
		skipLog = true
		mylog.Info("", "skip log in clang")
	}

	fnInitCmd := func() *exec.Cmd {
		if runtime.GOOS == "windows" {
			return exec.Command("cmd", "/C", command)
		}
		return exec.Command("bash", "-c", command)
	}
	cmd := fnInitCmd()

	log.Println(command)

	stdoutPipe := mylog.Check2(cmd.StdoutPipe())
	stderrPipe := mylog.Check2(cmd.StderrPipe())

	mylog.Check(cmd.Start())

	var wg sync.WaitGroup
	output := make(chan string)
	errorOutput := make(chan string)

	// 启动 goroutine 读取 stdout
	wg.Add(1)
	go func() {
		mylog.Call(func() {
			defer wg.Done()
			scanner := bufio.NewScanner(stdoutPipe)
			for scanner.Scan() {
				output <- ConvertUtf82Gbk(scanner.Text())
			}
			mylog.Check(stdoutPipe.Close())
		})
	}()

	// 启动 goroutine 读取 stderr
	wg.Add(1)
	go func() {
		mylog.Call(func() {
			defer wg.Done()
			scanner := bufio.NewScanner(stderrPipe)
			for scanner.Scan() {
				errorOutput <- ConvertUtf82Gbk(scanner.Text())
			}
			mylog.Check(stderrPipe.Close())
		})
	}()

	// 启动 goroutine 统一处理输出
	go func() {
		mylog.Call(func() {
			wg.Wait()
			close(output)
		})
	}()

	done := make(chan struct{})
	go func() {
		for line := range output {
			if !skipLog {
				println(line)
			}
			s.Output.WriteStringLn(line)
		}
		done <- struct{}{}
	}()

	go func() {
		for line := range errorOutput {
			println(line)
			s.Error.WriteStringLn(line)
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
		mylog.Check(ConvertUtf82Gbk(e.Error()))
	}
	s.Output.NewLine()
	s.Output.WriteStringLn(s.Error.String())
	ss := trimTrailingEmptyLines(s.Output.String())
	s.Output.Reset()
	s.Output.WriteString(ss)
}

func ConvertUtf82Gbk(src string) string {
	if IsWindows() {
		c := mylog.Check2(simplifiedchinese.GB18030.NewDecoder().String(src)) // todo test rune
		return strings.TrimSuffix(c, "\r\n")
	}
	return src
}

//func runCmd(command string) string { // std error not support
//	fnInitCmd := func() *exec.Cmd {
//		if runtime.GOOS == "windows" {
//			return exec.Command("cmd", "/C", command)
//		}
//		return exec.Command("bash", "-c", command)
//	}
//	return ConvertUtf82Gbk(string(mylog.Check2(fnInitCmd().CombinedOutput())))
//}
