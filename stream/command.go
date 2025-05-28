package stream

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/waitgroup"
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
		session.run(command, session.CurrentDirectory)
	})
	return
}

func RunCommandWithDir(command, dir string) (session *CommandSession) {
	mylog.Call(func() {
		session = &CommandSession{
			Output:           NewBuffer(""),
			Error:            NewBuffer(""),
			CurrentDirectory: dir,
		}
		session.run(command, session.CurrentDirectory)
	})
	return
}

func (s *CommandSession) run(command, dir string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	skipLog := false
	if strings.Contains(command, "clang") {
		skipLog = true
		//mylog.Info("", "skip log in clang")
	}

	fnInitCmd := func() *exec.Cmd {
		if runtime.GOOS == "windows" {
			return exec.Command("cmd", "/C", command)
		}
		return exec.Command("bash", "-c", command)
	}
	cmd := fnInitCmd()
	cmd.Dir = dir

	mylog.Info("command", command)

	stdoutPipe := mylog.Check2(cmd.StdoutPipe())
	stderrPipe := mylog.Check2(cmd.StderrPipe())

	mylog.Check(cmd.Start())

	 wg :=waitgroup.New()
	output := make(chan string)
	errorOutput := make(chan string)

	// 启动 goroutine 读取 stdout
	wg.Go(func() {
		mylog.Call(func() {
			scanner := bufio.NewScanner(stdoutPipe)
			for scanner.Scan() {
				output <- ConvertUtf82Gbk(scanner.Text())
			}
			mylog.Check(stdoutPipe.Close())
		})
	})

	// 启动 goroutine 读取 stderr
	wg.Go(func() {
		mylog.Call(func() {
			scanner := bufio.NewScanner(stderrPipe)
			for scanner.Scan() {
				errorOutput <- ConvertUtf82Gbk(scanner.Text())
			}
			mylog.Check(stderrPipe.Close())
		})
	})

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
				println(line)//对于json，不需要每一行都输出，而是一次性返回解码或者落地保存
			}
			s.Output.WriteStringLn(line)
		}
		done <- struct{}{}
	}()

	go func() {
		for line := range errorOutput {
			//println(line)//让clang dump ast的错误在后面写入文件，这里和后面调用s.Error.String()重复输出错误了
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
		mylog.Check(ConvertUtf82Gbk(e.Error() + "\n" + s.Error.String()))
	}
	//s.Output.NewLine()
	ss := trimTrailingEmptyLines(s.Output.String())
	s.Output.Reset()
	s.Output.WriteString(ss)
}

func trimTrailingEmptyLines(s string) string {
	// 使用正则表达式匹配末尾的所有空白行，包括空格、制表符和换行符
	re := regexp.MustCompile(`\s*\n*$`)
	return re.ReplaceAllString(s, "")
}

func ConvertUtf82Gbk(src string) string {
	if IsWindows() {
		c := mylog.Check2(simplifiedchinese.GB18030.NewDecoder().String(src)) // todo test rune
		return strings.TrimSuffix(c, "\r\n")
	}
	return src
}

// func runCmd(command string) string { // std error not support
//	fnInitCmd := func() *exec.Cmd {
//		if runtime.GOOS == "windows" {
//			return exec.Command("cmd", "/C", command)
//		}
//		return exec.Command("bash", "-c", command)
//	}
//	return ConvertUtf82Gbk(string(mylog.Check2(fnInitCmd().CombinedOutput())))
// }
