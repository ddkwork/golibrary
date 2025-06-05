package stream

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/waitgroup"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func RunCommand(arg ...string) (stdOut *GeneratedFile) {
	return runCommand("", arg...)
}

func RunCommandWithDir(dir string, arg ...string) (stdOut *GeneratedFile) {
	return runCommand(dir, arg...)
}

func runCommand(dir string, arg ...string) (stdOut *GeneratedFile) {
	if strings.Contains(arg[0], " ") {
		panic("you shold split commands")
	}
	type setup struct {
		init       func()
		fastModel  func()
		slowModel  func()
		handleWait func() //merge exit code.and error
	}
	cmdKey := "command" //fast the key is clang target file path
	fast := false
	stdOut = NewGeneratedFile()
	stderr := NewGeneratedFile()

	var (
		cmd         *exec.Cmd
		stdoutPipe  io.ReadCloser
		stderrPipe  io.ReadCloser
		output      = make(chan string)
		errorOutput = make(chan string)
		done        = make(chan struct{})
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := setup{
		init: func() {
			switch { //todo add more need fast model
			case arg[0] == "clang":
				cmdKey = filepath.Base(arg[len(arg)-1])
				fast = true
			}
			binaryPath := arg[0]
			cmd = exec.CommandContext(ctx, binaryPath, arg[1:]...)
			cmd.Dir = dir //需要切换到对应目录，避免使用os.chdir,应用场景：批量更新工作区下的mod

			mylog.Success(cmdKey, cmd.String())

			stdoutPipe = mylog.Check2(cmd.StdoutPipe())
			stderrPipe = mylog.Check2(cmd.StderrPipe())
		},
		fastModel: func() {
			mylog.Check2(stdOut.ReadFrom(stdoutPipe))
			mylog.Check2(stderr.ReadFrom(stderrPipe))
		},
		slowModel: func() {
			g := waitgroup.New()
			g.SetLimit(1000)

			// 启动 goroutine 读取 stdout
			g.Go(func() {
				mylog.Call(func() {
					scanner := bufio.NewScanner(stdoutPipe)
					for scanner.Scan() {
						output <- ConvertUtf82Gbk(scanner.Text())
					}
					mylog.Check(stdoutPipe.Close())
				})
			})

			// 启动 goroutine 读取 stderr
			g.Go(func() {
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
					g.Wait()
					close(output)
				})
			}()

			go func() {
				for line := range output {
					println(line)
					stdOut.P(line)
				}
				done <- struct{}{}
			}()

			go func() {
				for line := range errorOutput {
					stderr.P(line)
				}
				done <- struct{}{}
			}()

			select { // 等待 goroutine 完成,不写在这里会导致全部携程死锁，原因不明
			case <-ctx.Done():
				mylog.Check(cmd.Process.Kill())
			case <-done:
			}
		},
		handleWait: func() {
			e := cmd.Wait()
			if e != nil {
				bug := stderr.String() + "\n" + e.Error()
				mylog.Check(bug)
			}
		},
	}
	s.init()
	mylog.Check(cmd.Start())
	if fast {
		s.fastModel()
		s.handleWait()
		return
	}
	s.slowModel()
	s.handleWait()
	ss := trimTrailingEmptyLines(stdOut.String())
	stdOut.Reset()
	stdOut.WriteString(ss)
	return
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
