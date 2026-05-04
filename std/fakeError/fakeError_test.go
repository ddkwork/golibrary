package fakeError

import (
	"go/format"
	"go/parser"
	"go/token"
	"math/rand/v2"
	"os"
	"path/filepath"
	"testing"

	"github.com/ddkwork/golibrary/std/assert"
	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/golibrary/std/safemap"
	"github.com/ddkwork/golibrary/std/stream"
	"golang.org/x/arch/x86/x86asm"
)

func TestAll(t *testing.T) {
	t.Run("多个简单if_err_替换Check", func(t *testing.T) {
		assert.Equal(t, m.GetMust("多个简单if_err_替换Check").want, get("多个简单if_err_替换Check", m.GetMust("多个简单if_err_替换Check").code))
	})
	t.Run("for循环内if_err有continue_保留", func(t *testing.T) {
		assert.Equal(t, m.GetMust("for循环内if_err有continue_保留").want, get("for循环内if_err有continue_保留", m.GetMust("for循环内if_err有continue_保留").code))
	})
	t.Run("简单return_err_替换Check2", func(t *testing.T) {
		assert.Equal(t, m.GetMust("简单return_err_替换Check2").want, get("简单return_err_替换Check2", m.GetMust("简单return_err_替换Check2").code))
	})
	t.Run("goroutine内多err_替换Check", func(t *testing.T) {
		assert.Equal(t, m.GetMust("goroutine内多err_替换Check").want, get("goroutine内多err_替换Check", m.GetMust("goroutine内多err_替换Check").code))
	})
	t.Run("for内continue加logFatal_替换", func(t *testing.T) {
		assert.Equal(t, m.GetMust("for内continue加logFatal_替换").want, get("for内continue加logFatal_替换", m.GetMust("for内continue加logFatal_替换").code))
	})
	t.Run("for内continue加logWarn_替换", func(t *testing.T) {
		assert.Equal(t, m.GetMust("for内continue加logWarn_替换").want, get("for内continue加logWarn_替换", m.GetMust("for内continue加logWarn_替换").code))
	})
	t.Run("return_nil_nil_err_替换Check2", func(t *testing.T) {
		assert.Equal(t, m.GetMust("return_nil_nil_err_替换Check2").want, get("return_nil_nil_err_替换Check2", m.GetMust("return_nil_nil_err_替换Check2").code))
	})
	t.Run("复杂条件err_替换Check", func(t *testing.T) {
		assert.Equal(t, m.GetMust("复杂条件err_替换Check").want, get("复杂条件err_替换Check", m.GetMust("复杂条件err_替换Check").code))
	})
	t.Run("多个return_nil_nil_err_替换Check2", func(t *testing.T) {
		assert.Equal(t, m.GetMust("多个return_nil_nil_err_替换Check2").want, get("多个return_nil_nil_err_替换Check2", m.GetMust("多个return_nil_nil_err_替换Check2").code))
	})
	t.Run("defer处理_各种error模式", func(t *testing.T) {
		assert.Equal(t, m.GetMust("defer处理_各种error模式").want, get("defer处理_各种error模式", m.GetMust("defer处理_各种error模式").code))
	})
	t.Run("无err变量_不替换", func(t *testing.T) {
		assert.Equal(t, m.GetMust("无err变量_不替换").want, get("无err变量_不替换", m.GetMust("无err变量_不替换").code))
	})
	t.Run("下划线err赋值_替换Check2", func(t *testing.T) {
		assert.Equal(t, m.GetMust("下划线err赋值_替换Check2").want, get("下划线err赋值_替换Check2", m.GetMust("下划线err赋值_替换Check2").code))
	})
	t.Run("已有mylog代码_保持不变", func(t *testing.T) {
		assert.Equal(t, m.GetMust("已有mylog代码_保持不变").want, get("已有mylog代码_保持不变", m.GetMust("已有mylog代码_保持不变").code))
	})
	t.Run("已有mylog代码UDP_保持不变", func(t *testing.T) {
		assert.Equal(t, m.GetMust("已有mylog代码UDP_保持不变").want, get("已有mylog代码UDP_保持不变", m.GetMust("已有mylog代码UDP_保持不变").code))
	})
	t.Run("if_else_if链_替换Check2", func(t *testing.T) {
		assert.Equal(t, m.GetMust("if_else_if链_替换Check2").want, get("if_else_if链_替换Check2", m.GetMust("if_else_if链_替换Check2").code))
	})
	t.Run("osExit_替换Check_for内补break", func(t *testing.T) {
		assert.Equal(t, m.GetMust("osExit_替换Check_for内补break").want, get("osExit_替换Check_for内补break", m.GetMust("osExit_替换Check_for内补break").code))
	})
	t.Run("简单return_替换Check2_defer闭包", func(t *testing.T) {
		assert.Equal(t, m.GetMust("简单return_替换Check2_defer闭包").want, get("简单return_替换Check2_defer闭包", m.GetMust("简单return_替换Check2_defer闭包").code))
	})
	t.Run("简单return_替换Check2_无defer", func(t *testing.T) {
		assert.Equal(t, m.GetMust("简单return_替换Check2_无defer").want, get("简单return_替换Check2_无defer", m.GetMust("简单return_替换Check2_无defer").code))
	})
	t.Run("复杂分支err判断_不替换", func(t *testing.T) {
		assert.Equal(t, m.GetMust("复杂分支err判断_不替换").want, get("复杂分支err判断_不替换", m.GetMust("复杂分支err判断_不替换").code))
	})
	t.Run("mylogWarning_return_替换Check2", func(t *testing.T) {
		assert.Equal(t, m.GetMust("mylogWarning_return_替换Check2").want, get("mylogWarning_return_替换Check2", m.GetMust("mylogWarning_return_替换Check2").code))
	})
	t.Run("for循环内if_err已有break_不替换", func(t *testing.T) {
		assert.Equal(t, m.GetMust("for循环内if_err已有break_不替换").want, get("for循环内if_err已有break_不替换", m.GetMust("for循环内if_err已有break_不替换").code))
	})
	t.Run("方法签名移除error返回", func(t *testing.T) {
		assert.Equal(t, m.GetMust("方法签名移除error返回").want, get("方法签名移除error返回", m.GetMust("方法签名移除error返回").code))
	})
	t.Run("函数签名移除error返回", func(t *testing.T) {
		assert.Equal(t, m.GetMust("函数签名移除error返回").want, get("函数签名移除error返回", m.GetMust("函数签名移除error返回").code))
	})
	t.Run("接口定义方法签名移除error", func(t *testing.T) {
		assert.Equal(t, m.GetMust("接口定义方法签名移除error").want, get("接口定义方法签名移除error", m.GetMust("接口定义方法签名移除error").code))
	})
	t.Run("接口签名多返回值移除error", func(t *testing.T) {
		assert.Equal(t, m.GetMust("接口签名多返回值移除error").want, get("接口签名多返回值移除error", m.GetMust("接口签名多返回值移除error").code))
	})
	t.Run("复杂err分支_可替换Check", func(t *testing.T) {
		assert.Equal(t, m.GetMust("复杂err分支_可替换Check").want, get("复杂err分支_可替换Check", m.GetMust("复杂err分支_可替换Check").code))
	})
	t.Run("有业务逻辑_不替换", func(t *testing.T) {
		assert.Equal(t, m.GetMust("有业务逻辑_不替换").want, get("有业务逻辑_不替换", m.GetMust("有业务逻辑_不替换").code))
	})
	t.Run("简单return_替换Check2_ioWrite", func(t *testing.T) {
		assert.Equal(t, m.GetMust("简单return_替换Check2_ioWrite").want, get("简单return_替换Check2_ioWrite", m.GetMust("简单return_替换Check2_ioWrite").code))
	})
}

func get(path, text string) string {
	join := filepath.Join(os.TempDir(), path+".go")
	fix := filepath.Join(os.TempDir(), path+"_fixed.go")
	mylog.Warning("original code file path", join+":1")
	stream.WriteGoFile(join, text) // 写入文件只是为了让goland检查原始代码的语法和直观的行号
	ret := ""
	mylog.Call(func() {
		fileSet := token.NewFileSet()
		file := mylog.Check2(parser.ParseFile(fileSet, fix, text, parser.ParseComments))
		ret = handle(fileSet, file, text)
		mylog.WriteGoFile(fix, ret) // 写入文件只是为了让goland检查返回代码的语法和直观的行号
	})
	return string(mylog.Check2(format.Source([]byte(ret)))) // handle中WriteGoFile已经执行格式化
}

type testData struct {
	code string
	want string
}

var m = safemap.NewOrdered[string, testData](func(yield func(string, testData) bool) {
	yield("多个简单if_err_替换Check", testData{
		code: `package tmp

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	var err error
	if err != nil {
		log.Debug(err)
	}
	if err != nil {
		panic(err)
	}
	if err := backendConn.Close(); err != nil {
		log.Debug(err)
	}
}
`,
		want: `package tmp

import (
	"github.com/ddkwork/golibrary/std/mylog"
	log "github.com/sirupsen/logrus"
)

func main() {

	mylog.Check(backendConn.Close())
}
`,
	})
	yield("for循环内if_err有continue_保留", testData{
		code: `package main

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"
)

func main() {
	config, err := initCliParas()
	if err != nil {
		log.Fatal(err)
	}

	for _, rule := range config.Rules {
		rule := rule
		go func() {
			l, err := net.Listen("tcp", rule.ListenAddr)
			if err != nil {
				log.Fatal(err)
			}
			for {
				conn, err := l.Accept()
				if err != nil {
					log.Debug(err)
					continue
				}
				go handleConnection(conn, rule.ForwardTargets)
			}
		}()
	}
	select {}
}

type config struct {
	Rules []ruleConfig
}

type ruleConfig struct {
	ListenAddr string
	// ServerName: ForwardAddress
	ForwardTargets map[string]string
}

func initCliParas() (*config, error) {
	log.SetLevel(log.DebugLevel)
	var config config
	config.Rules = []ruleConfig{
		{
			ListenAddr: ":7890",
			ForwardTargets: map[string]string{
				"www.google.com": "http://www.recaptcha.net",
			},
		},
	}
	for _, rule := range config.Rules {
		for serverName, target := range rule.ForwardTargets {
			log.Info("ADD [" + serverName + " -> " + target + "] at [" + rule.ListenAddr + "]")
		}
	}
	return &config, nil
}

func getForwardTarget(serverName string, forwardTargets map[string]string) (target string, allowed bool) {
	if _, ok := forwardTargets[serverName]; ok {
		return forwardTargets[serverName], true
	}
	for keyServerName, valueForwardTarget := range forwardTargets {
		if strings.HasPrefix(keyServerName, "*") {
			keyServerName = ".*(" + strings.ReplaceAll(keyServerName[1:], ".", "\\.") + ")$"
			matched, err := regexp.Match(keyServerName, []byte(serverName))
			if err != nil {
				log.Warn("Error when matching sni with allowed sni")
				continue
			}
			if matched {
				return valueForwardTarget, true
			}
		}
	}
	return "", false
}

func handleConnection(clientConn net.Conn, forwardTargets map[string]string) {
	defer func(clientConn net.Conn) {
		if err := clientConn.Close(); err != nil {
			log.Debug(err)
		}
	}(clientConn)

	if err := clientConn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		log.Debug(err)
		return
	}

	clientHello, clientReader, err := PeekClientHello(clientConn)
	if err != nil {
		log.Debug(err)
		return
	}

	// 设置为不会超时
	if err := clientConn.SetReadDeadline(time.Time{}); err != nil {
		log.Debug(err)
		return
	}

	forwardTarget, ok := getForwardTarget(clientHello.ServerName, forwardTargets)
	if !ok {
		log.Debug("Blocking connection to unauthorized backend.")
		log.Debug("Source Addr: " + clientConn.RemoteAddr().String())
		log.Debug("Target SNI: " + clientHello.ServerName)
		return
	}

	backendConn, err := net.DialTimeout("tcp", forwardTarget, 5*time.Second)
	if err != nil {
		log.Warn(err)
		return
	}
	defer func(backendConn net.Conn) {
		if err := backendConn.Close(); err != nil {
			log.Debug(err)
		}
	}(backendConn)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		if _, err := io.Copy(clientConn, backendConn); err != nil {
			log.Debug(err)
		}
		if err := clientConn.(*net.TCPConn).CloseWrite(); err != nil {
			log.Debug(err)
		}
		wg.Done()
	}()
	go func() {
		if _, err := io.Copy(backendConn, clientReader); err != nil {
			log.Debug(err)
		}
		if err := backendConn.(*net.TCPConn).CloseWrite(); err != nil {
			log.Debug(err)
		}
		wg.Done()
	}()

	wg.Wait()
}
`,
		want: `package main

import (
	"github.com/ddkwork/golibrary/std/mylog"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"
)

func main() {
	config := mylog.Check2(initCliParas())

	for _, rule := range config.Rules {
		rule := rule
		go func() {
			l := mylog.Check2(net.Listen("tcp", rule.ListenAddr))

			for {
				conn, err := l.Accept()
				if err != nil {
					mylog.CheckIgnore(err)
					continue
				}
				go handleConnection(conn, rule.ForwardTargets)
			}
		}()
	}
	select {}
}

type config struct {
	Rules []ruleConfig
}

type ruleConfig struct {
	ListenAddr string
	// ServerName: ForwardAddress
	ForwardTargets map[string]string
}

func initCliParas() (*config, error) {
	log.SetLevel(log.DebugLevel)
	var config config
	config.Rules = []ruleConfig{
		{
			ListenAddr: ":7890",
			ForwardTargets: map[string]string{
				"www.google.com": "http://www.recaptcha.net",
			},
		},
	}
	for _, rule := range config.Rules {
		for serverName, target := range rule.ForwardTargets {
			log.Info("ADD [" + serverName + " -> " + target + "] at [" + rule.ListenAddr + "]")
		}
	}
	return &config, nil
}

func getForwardTarget(serverName string, forwardTargets map[string]string) (target string, allowed bool) {
	if _, ok := forwardTargets[serverName]; ok {
		return forwardTargets[serverName], true
	}
	for keyServerName, valueForwardTarget := range forwardTargets {
		if strings.HasPrefix(keyServerName, "*") {
			keyServerName = ".*(" + strings.ReplaceAll(keyServerName[1:], ".", "\\.") + ")$"
			matched, err := regexp.Match(keyServerName, []byte(serverName))
			if err != nil {
				mylog.CheckIgnore(err)
				continue
			}
			if matched {
				return valueForwardTarget, true
			}
		}
	}
	return "", false
}

func handleConnection(clientConn net.Conn, forwardTargets map[string]string) {
	defer func(clientConn net.Conn) {
		mylog.Check(clientConn.Close())
	}(clientConn)

	mylog.Check(clientConn.SetReadDeadline(time.Now().Add(5 * time.Second)))

	clientHello, clientReader := mylog.Check3(PeekClientHello(clientConn))

	// 设置为不会超时
	mylog.Check(clientConn.SetReadDeadline(time.Time{}))

	forwardTarget, ok := getForwardTarget(clientHello.ServerName, forwardTargets)
	if !ok {
		log.Debug("Blocking connection to unauthorized backend.")
		log.Debug("Source Addr: " + clientConn.RemoteAddr().String())
		log.Debug("Target SNI: " + clientHello.ServerName)
		return
	}

	backendConn := mylog.Check2(net.DialTimeout("tcp", forwardTarget, 5*time.Second))

	defer func(backendConn net.Conn) {
		mylog.Check(backendConn.Close())
	}(backendConn)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		mylog.Check2(io.Copy(clientConn, backendConn))
		mylog.Check(clientConn.(*net.TCPConn).CloseWrite())
		wg.Done()
	}()
	go func() {
		mylog.Check2(io.Copy(backendConn, clientReader))
		mylog.Check(backendConn.(*net.TCPConn).CloseWrite())
		wg.Done()
	}()

	wg.Wait()
}
`,
	})
	yield("简单return_err_替换Check2", testData{
		code: `package main

import (
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

func main() {
	backendConn, err := net.DialTimeout("tcp", forwardTarget, 5*time.Second)
	if err != nil {
		log.Warn(err)
		return
	}
}`,
		want: `package main

import (
	"github.com/ddkwork/golibrary/std/mylog"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

func main() {
	backendConn := mylog.Check2(net.DialTimeout("tcp", forwardTarget, 5*time.Second))

}
`,
	})
	yield("goroutine内多err_替换Check", testData{
		code: `package main

import (
	"github.com/ddkwork/golibrary/std/mylog"
	"io"
	"net"
)

func main() {
	go func() {
		if _, err := io.Copy(clientConn, backendConn); err != nil {
			log.Debug(err)
		}
		if err := clientConn.(*net.TCPConn).CloseWrite(); err != nil {
			log.Debug(err)
		}
		wg.Done()
	}()
}
`,
		want: `package main

import (
	"github.com/ddkwork/golibrary/std/mylog"
	"io"
	"net"
)

func main() {
	go func() {
		mylog.Check2(io.Copy(clientConn, backendConn))
		mylog.Check(clientConn.(*net.TCPConn).CloseWrite())
		wg.Done()
	}()
}
`,
	})
	yield("for内continue加logFatal_替换", testData{
		code: `package main

import (
	log "github.com/sirupsen/logrus"
	"net"
)

func main() {
	config, err := initCliParas()
	if err != nil {
		log.Fatal(err)
	}

	for _, rule := range config.Rules {
		rule := rule
		go func() {
			l, err := net.Listen("tcp", rule.ListenAddr)

if err != nil {
				log.Fatal(err)
			}
			for {
				conn, err := l.Accept()
				if err != nil {
					log.Debug(err)
					continue
				}
				go handleConnection(conn, rule.ForwardTargets)
			}
		}()
	}
	select {}
}
`,
		want: `package main

import (
	"github.com/ddkwork/golibrary/std/mylog"
	log "github.com/sirupsen/logrus"
	"net"
)

func main() {
	config := mylog.Check2(initCliParas())

	for _, rule := range config.Rules {
		rule := rule
		go func() {
			l := mylog.Check2(net.Listen("tcp", rule.ListenAddr))

			for {
				conn, err := l.Accept()
				if err != nil {
					mylog.CheckIgnore(err)
					continue
				}
				go handleConnection(conn, rule.ForwardTargets)
			}
		}()
	}
	select {}
}
`,
	})
	yield("for内continue加logWarn_替换", testData{
		code: `package main

import (
	"regexp"
	"strings"
)

func main() {
	for {
		keyServerName = ".*(" + strings.ReplaceAll(keyServerName[1:], ".", "\\.") + ")$"
		matched, err := regexp.Match(keyServerName, []byte(serverName))
		if err != nil {
			log.Warn("Error when matching sni with allowed sni")
			continue
		}
	}
}
`,
		want: `package main

import (
	"github.com/ddkwork/golibrary/std/mylog"
	"regexp"
	"strings"
)

func main() {
	for {
		keyServerName = ".*(" + strings.ReplaceAll(keyServerName[1:], ".", "\\.") + ")$"
		matched, err := regexp.Match(keyServerName, []byte(serverName))
		if err != nil {
			mylog.CheckIgnore(err)
			continue
		}
	}
}
`,
	})
	yield("return_nil_nil_err_替换Check2", testData{
		code: `package main

import (
	"github.com/coreos/go-oidc"
)

func mian() {
	provider, err := oidc.NewProvider(c.Ctx, c.ProviderURL)
	if err != nil {
		return nil, nil, err
	}
}`,
		want: `package main

import (
	"github.com/coreos/go-oidc"
	"github.com/ddkwork/golibrary/std/mylog"
)

func mian() {
	provider := mylog.Check2(oidc.NewProvider(c.Ctx, c.ProviderURL))

}
`,
	})
	yield("复杂条件err_替换Check", testData{
		code: `package tmp

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	err := jsonx.Open(&token, tf)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, nil, err
	}
}`,
		want: `package tmp

import (
	"github.com/ddkwork/golibrary/std/mylog"
	log "github.com/sirupsen/logrus"
)

func main() {
	mylog.Check(jsonx.Open(&token, tf))

}
`,
	})
	yield("多个return_nil_nil_err_替换Check2", testData{
		code: `package tmp

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, nil, err
	}

	userInfo, err := provider.UserInfo(c.Ctx, tokenSource)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user info: %w", err)
	}
}
`,
		want: `package tmp

import (
	"github.com/ddkwork/golibrary/std/mylog"
	log "github.com/sirupsen/logrus"
)

func main() {
	newToken := mylog.Check2(tokenSource.Token())

	userInfo := mylog.Check2(provider.UserInfo(c.Ctx, tokenSource))

}
`,
	})
	yield("defer处理_各种error模式", testData{
		code: `package tmp

import (
	"github.com/ddkwork/golibrary/std/mylog"
	log "github.com/sirupsen/logrus"
)

func main() {
	panic(bug())
}

func bug() error {
	if _, err := w.Write([]byte("Hello, 世界")); err != nil {
		t.Errorf("could not write assets/hello_world.txt: %v", err)
	}
	var err error
	if err != nil {
		log.Error(err)
		return err
	}
	if err != nil {
		return err
	}
	if err != nil {
		return fmt.Errorf("input must be struct pointer")
	}
	if err != nil {
		return nil, fmt.Errorf("input must be struct")
	}
	if err := apkw.Close(); err != nil {
		t.Fatal(err)
	}
	defer func() { io.Close(nil) }()
	defer io.Close(nil)
	defer func() { io.Close(nil) }()
}
`,
		want: `package tmp

import (
	"github.com/ddkwork/golibrary/std/mylog"
	log "github.com/sirupsen/logrus"
)

func main() {
	panic(bug())
}

func bug() error {
	mylog.Check2(w.Write([]byte("Hello, 世界")))

	mylog.Check(apkw.Close())
	defer func() { io.Close(nil) }()
	defer func() { mylog.Check(io.Close(nil)) }()
	defer func() { io.Close(nil) }()
}
`,
	})
	yield("无err变量_不替换", testData{
		code: `package tmp

import (
	"syscall"
)

func test() {
	var GetLogicalDrives *syscall.LazyProc
	n, _ := GetLogicalDrives.Call()
	_ = n
}
`,
		want: `package tmp

import (
	"syscall"
)

func test() {
	var GetLogicalDrives *syscall.LazyProc
	n, _ := GetLogicalDrives.Call()
	_ = n
}
`,
	})
	yield("下划线err赋值_替换Check2", testData{
		code: `package tmp

import (
	"path/filepath"
)

func test() {
	files, _ := filepath.Glob("*.go")
	_ = files
}
`,
		want: `package tmp

import (
	"github.com/ddkwork/golibrary/std/mylog"
	"path/filepath"
)

func test() {
	files := mylog.Check2(filepath.Glob("*.go"))
	_ = files
}
`,
	})
	yield("已有mylog代码_保持不变", testData{
		code: `package stream

import (
	"fmt"
	"net"

	"github.com/ddkwork/golibrary/std/mylog"
)

// GetAvailablePort 获取可用端口
func GetAvailablePort() int {
	listener := mylog.Check2(net.ListenTCP("tcp", mylog.Check2(net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", "0.0.0.0")))))
	defer func() { mylog.Check(listener.Close()) }()
	return listener.Addr().(*net.TCPAddr).Port
}

// IsPortAvailable 判断端口是否可以（未被占用）
func IsPortAvailable(port int) bool {
	address := fmt.Sprintf("%s:%d", "0.0.0.0", port)
	listener := mylog.Check2(net.Listen("tcp", address))

	defer mylog.Check(listener.Close())
	return true
}
`,
		want: `package stream

import (
	"fmt"
	"net"

	"github.com/ddkwork/golibrary/std/mylog"
)

// GetAvailablePort 获取可用端口
func GetAvailablePort() int {
	listener := mylog.Check2(net.ListenTCP("tcp", mylog.Check2(net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", "0.0.0.0")))))
	defer func() { mylog.Check(listener.Close()) }()
	return listener.Addr().(*net.TCPAddr).Port
}

// IsPortAvailable 判断端口是否可以（未被占用）
func IsPortAvailable(port int) bool {
	address := fmt.Sprintf("%s:%d", "0.0.0.0", port)
	listener := mylog.Check2(net.Listen("tcp", address))

	defer mylog.Check(listener.Close())
	return true
}
`,
	})
	yield("已有mylog代码UDP_保持不变", testData{
		code: `package udp

import "github.com/ddkwork/golibrary/std/mylog"

func (o *object) TransportUDP(DstIP string, DstPort int) {
	o.reset(DstIP, DstPort)
	o.GetSrcAddrConn()
	defer func() {
		mylog.Check(o.SrcConn == nil)
		mylog.Check(o.SrcConn.Close())
		mylog.Check(o.DstConn == nil)
		mylog.Check(o.DstConn.Close())
	}()
	for {
		SrcBufChan <- o.Bytes()[:o.BufSize]
		o.SetDstAddrConn()
		go o.readDstBuf()
		mylog.Check2(o.SrcConn.WriteToUDP(<-DstBufChan, o.SrcAddr))
	}
}

func (o *object) readDstBuf() {
	select {
	case b := <-SrcBufChan:
		mylog.Check2(o.DstConn.Write(b))
		o.Reset()
		o.BufSize = mylog.Check2(o.DstConn.Read(o.Bytes()))
		DstBufChan <- o.Bytes()[:o.BufSize]
	default:
	}
}`,
		want: `package udp

import "github.com/ddkwork/golibrary/std/mylog"

func (o *object) TransportUDP(DstIP string, DstPort int) {
	o.reset(DstIP, DstPort)
	o.GetSrcAddrConn()
	defer func() {
		mylog.Check(o.SrcConn == nil)
		mylog.Check(o.SrcConn.Close())
		mylog.Check(o.DstConn == nil)
		mylog.Check(o.DstConn.Close())
	}()
	for {
		SrcBufChan <- o.Bytes()[:o.BufSize]
		o.SetDstAddrConn()
		go o.readDstBuf()
		mylog.Check2(o.SrcConn.WriteToUDP(<-DstBufChan, o.SrcAddr))
	}
}

func (o *object) readDstBuf() {
	select {
	case b := <-SrcBufChan:
		mylog.Check2(o.DstConn.Write(b))
		o.Reset()
		o.BufSize = mylog.Check2(o.DstConn.Read(o.Bytes()))
		DstBufChan <- o.Bytes()[:o.BufSize]
	default:
	}
}
`,
	})
	yield("if_else_if链_替换Check2", testData{
		code: `package tmp

func main() {
	if f, err := os.Create("ntstatus_generated.go"); err != nil {
		log.Fatal(err)
	} else if n, err := f.Write(out); err != nil {
		log.Fatal(err)
	} else if n != len(out) {
		log.Fatal("output size mismatch")
	} else {
		f.Close()
	}
}`,
		want: `package tmp

func main() {
	f := mylog.Check2(os.Create("ntstatus_generated.go"))
	n := mylog.Check2(f.Write(out))
	if n != len(out) {
		mylog.Check("output size mismatch")
	}
	mylog.Check(f.Close())

}
`,
	})
	yield("osExit_替换Check_for内补break", testData{
		code: `package tmp

func main() {
	if err := generateHelperFile(); err != nil {
		fmt.Fprintf(os.Stderr, "生成辅助函数文件失败: %v\n", err)
		os.Exit(1)
	}

	for _, config := range configs {
		if err := generateMCPServer(interfacePath, config); err != nil {
			fmt.Fprintf(os.Stderr, "生成 %s 失败: %v\n", config.Interface, err)
			os.Exit(1)
		}
	}
}`,
		want: `package tmp

func main() {
	mylog.Check(generateHelperFile())

	for _, config := range configs {
		mylog.Check(generateMCPServer(interfacePath, config))
	}
}
`,
	})
	yield("简单return_替换Check2_defer闭包", testData{
		code: `package tmp

import "os"

func main() {
	f, err := os.Open("")
	if err != nil {
		return
	}
	defer f.Close()
}`,
		want: `package tmp

import (
	"github.com/ddkwork/golibrary/std/mylog"
	"os"
)

func main() {
	f := mylog.Check2(os.Open(""))

	defer func() { mylog.Check(f.Close()) }()
}
`,
	})
	yield("简单return_替换Check2_无defer", testData{
		code: `package tmp

import "os"

func main() {
	f, err := os.Open("")
	if err != nil {
		return
	}
	_ = f
}`,
		want: `package tmp

import (
	"github.com/ddkwork/golibrary/std/mylog"
	"os"
)

func main() {
	f := mylog.Check2(os.Open(""))

	_ = f
}
`,
	})
	yield("复杂分支err判断_不替换", testData{
		code: `package tmp

func main() {
	schService, err := CreateService(sc)
	if err != nil {
		if err == ERROR_SERVICE_EXISTS {
			return false
		}
		if err == ERROR_SERVICE_MARKED_FOR_DELETE {
			return false
		}
		return false
	}
	_ = schService
}`,
		want: `package tmp

func main() {
	schService := mylog.Check2(CreateService(sc))

	_ = schService
}
`,
	})
	yield("mylogWarning_return_替换Check2", testData{
		code: `package tmp

func main() {
	h, err := CreateFile(namePtr, GENERIC_READ|GENERIC_WRITE, 0, nil, OPEN_EXISTING, 0, 0)
	if err != nil {
		mylog.Warning("CreateFile failed", "error", err)
		return false
	}
	_ = h
}`,
		want: `package tmp

func main() {
	h := mylog.Check2(CreateFile(namePtr, GENERIC_READ|GENERIC_WRITE, 0, nil, OPEN_EXISTING, 0, 0))

	_ = h
}
`,
	})
	yield("for循环内if_err已有break_不替换", testData{
		code: `package tmp

func main() {
	for off := 0; off < len(data); {
		inst, err := x86asm.Decode(data[off:], 64)
		if err != nil {
			break
		}
		off += inst.Len
	}
}`,
		want: `package tmp

func main() {
	for off := 0; off < len(data); {
		inst, err := x86asm.Decode(data[off:], 64)
		if err != nil {
			break
		}
		off += inst.Len
	}
}
`,
	})
	yield("方法签名移除error返回", testData{
		code: `package tmp

type KernelMemory struct {
	rt interface{}
}

func (k *KernelMemory) WriteUint32(addr uint64, val uint32) error {
	return k.rt.WriteUint32(addr, val)
}
`,
		want: `package tmp

type KernelMemory struct {
	rt interface{}
}

func (k *KernelMemory) WriteUint32(addr uint64, val uint32) {
	k.rt.WriteUint32(addr, val)
}
`,
	})
	yield("接口定义方法签名移除error", testData{
		code: `package tmp

type Reader interface {
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Close() error
}
`,
		want: `package tmp

type Reader interface {
	Read(p []byte) (n int)
	Write(p []byte) (n int)
	Close()
}
`,
	})
	yield("函数签名移除error返回", testData{
		code: `package tmp

func DoSomething(a int, b string) error {
	return someLib.Call(a, b)
}
`,
		want: `package tmp

func DoSomething(a int, b string) {
	someLib.Call(a, b)
}
`,
	})
	yield("接口签名多返回值移除error", testData{
		code: `package tmp

type RTCore64 struct {
	deviceHandle uintptr
}

func (r *RTCore64) readDword(addr uint64) (uint32, error) {
	var pkt struct {
		addr uint64
		size uint32
		value uint32
	}
	pkt.addr = addr
	pkt.size = 4

	mylog.Check(someFunc(r.deviceHandle, pkt))

	return pkt.value, nil
}
`,
		want: `package tmp

type RTCore64 struct {
	deviceHandle uintptr
}

func (r *RTCore64) readDword(addr uint64) uint32 {
	var pkt struct {
		addr  uint64
		size  uint32
		value uint32
	}
	pkt.addr = addr
	pkt.size = 4

	mylog.Check(someFunc(r.deviceHandle, pkt))

	return pkt.value
}
`,
	})
	yield("复杂err分支_可替换Check", testData{
		code: `package tmp

import (
	"github.com/ddkwork/golibrary/std/mylog"
)

type Driver struct {
	Name string
	Path string
}

func (d *Driver) Install() bool {
	return d.withSCManager(func(sc interface{}) bool {
		driverNamePtr := mylog.Check2(windows.UTF16PtrFromString(d.Name))
		serviceExePtr := mylog.Check2(windows.UTF16PtrFromString(d.Path))

		schService, err := windows.CreateService(
			sc,
			driverNamePtr,
			driverNamePtr,
			windows.SERVICE_ALL_ACCESS,
			windows.SERVICE_KERNEL_DRIVER,
			windows.SERVICE_DEMAND_START,
			windows.SERVICE_ERROR_NORMAL,
			serviceExePtr,
			nil, nil, nil, nil, nil,
		)
		if err != nil {
			if err == windows.ERROR_SERVICE_EXISTS {
				mylog.Warning("service already exists", "name", d.Name)
				return false
			}
			if err == windows.ERROR_SERVICE_MARKED_FOR_DELETE {
				mylog.Warning("previous instance of the service is not fully deleted. Try again...")
				return false
			}
			mylog.Warning("CreateService failed", "error", err)
			return false
		}

		if schService != 0 {
			windows.CloseServiceHandle(schService)
		}

		mylog.Success("driver installed successfully")
		return true
	})
}
`,
		want: `package tmp

import (
	"github.com/ddkwork/golibrary/std/mylog"
)

type Driver struct {
	Name string
	Path string
}

func (d *Driver) Install() bool {
	return d.withSCManager(func(sc interface{}) bool {
		driverNamePtr := mylog.Check2(windows.UTF16PtrFromString(d.Name))
		serviceExePtr := mylog.Check2(windows.UTF16PtrFromString(d.Path))

		schService := mylog.Check2(windows.CreateService(
			sc,
			driverNamePtr,
			driverNamePtr,
			windows.SERVICE_ALL_ACCESS,
			windows.SERVICE_KERNEL_DRIVER,
			windows.SERVICE_DEMAND_START,
			windows.SERVICE_ERROR_NORMAL,
			serviceExePtr,
			nil, nil, nil, nil, nil,
		))

		if schService != 0 {
			windows.CloseServiceHandle(schService)
		}

		mylog.Success("driver installed successfully")
		return true
	})
}
`,
	})
	yield("有业务逻辑_不替换", testData{
		code: `package tmp

import (
	"fmt"
	"github.com/ddkwork/golibrary/std/mylog"
)

const RTCORE_DEVICE_NAME = "RTCore64"

type RTCore struct {
	driver *Driver
	deviceHandle uintptr
}

func (r *RTCore) Open() bool {
	namePtr, _ := windows.UTF16PtrFromString(fmt.Sprintf("\\\\.\\%s", RTCORE_DEVICE_NAME))
	h, err := windows.CreateFile(namePtr,
		windows.GENERIC_READ|windows.GENERIC_WRITE,
		0, nil, windows.OPEN_EXISTING, 0, 0)
	if err != nil {
		mylog.Warning("CreateFile for RTCore64 failed", "error", err)
		r.driver.Stop()
		r.driver.Remove()
		r.driver = nil
		return false
	}

	if r.deviceHandle != 0 && r.deviceHandle != windows.InvalidHandle {
		windows.CloseHandle(r.deviceHandle)
		r.deviceHandle = 0
	}
	return true
}
`,
		want: `package tmp

import (
	"fmt"
	"github.com/ddkwork/golibrary/std/mylog"
)

const RTCORE_DEVICE_NAME = "RTCore64"

type RTCore struct {
	driver       *Driver
	deviceHandle uintptr
}

func (r *RTCore) Open() bool {
	namePtr, _ := windows.UTF16PtrFromString(fmt.Sprintf("\\\\.\\%s", RTCORE_DEVICE_NAME))
	h := mylog.Check2(windows.CreateFile(namePtr,
		windows.GENERIC_READ|windows.GENERIC_WRITE,
		0, nil, windows.OPEN_EXISTING, 0, 0))
	if err != nil {
		mylog.Warning("CreateFile for RTCore64 failed", "error", err)
		r.driver.Stop()
		r.driver.Remove()
		r.driver = nil
		return false
	}

	if r.deviceHandle != 0 && r.deviceHandle != windows.InvalidHandle {
		windows.CloseHandle(r.deviceHandle)
		r.deviceHandle = 0
	}
	return true
}
`,
	})
	yield("简单return_替换Check2_ioWrite", testData{
		code: `package tmp

import "io"

func main() {
	var w io.Writer
	n, err := w.Write([]byte("hello"))
	if err != nil {
		return
	}
	_ = n
}`,
		want: `package tmp

import (
	"github.com/ddkwork/golibrary/std/mylog"
	"io"
)

func main() {
	var w io.Writer
	n := mylog.Check2(w.Write([]byte("hello")))

	_ = n
}
`,
	})
})

// FuzzX86asmBreakVsCheck2 用真正的 Go fuzz 验证随机数据下 break 与 Check2 的行为差异
func FuzzX86asmBreakVsCheck2(f *testing.F) {
	seed1 := make([]byte, 4096)
	rand.New(rand.NewChaCha8([32]byte(seed1)))

	seed2 := make([]byte, 64)
	for i := range seed2 {
		seed2[i] = byte(0xFF) // 全0xFF垃圾
	}

	seed3 := make([]byte, 128)
	for i := range seed3 {
		seed3[i] = byte(i*7 + 3) // 模式化垃圾
	}

	seed4 := []byte{0x48, 0x89, 0xD8, 0xC3} // 正常指令

	f.Add(seed1)
	f.Add(seed2)
	f.Add(seed3)
	f.Add(seed4)
	f.Add([]byte{})
	f.Add([]byte{0x00})
	f.Add([]byte{0xFF, 0xFF, 0xFF})

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) == 0 {
			return
		}
		off := 0
		for off < len(data) {
			inst, err := x86asm.Decode(data[off:], 64)
			if err != nil { //ErrUnrecognized = errors.New("unrecognized instruction")
				//mylog.Struct(inst) //失败的时候它居然不是nil
				//mylog.CheckIgnore(err)//不会panic
				break //现在我明白我们的预期了，break需要check插入，前提是判断err := 在 for 循环内
				//panic(err)
			}
			off += inst.Len
		}
	})
}
