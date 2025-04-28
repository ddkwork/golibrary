package fakeError

import (
	"github.com/ddkwork/golibrary/assert"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/safemap"
	"github.com/ddkwork/golibrary/stream"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"
)

func Test1(t *testing.T) {
	assert.Equal(t, m.GetMust("test1").want, get("test1", m.GetMust("test1").code))
}
func Test2(t *testing.T) {
	assert.Equal(t, m.GetMust("test2").want, get("test2", m.GetMust("test2").code))
}
func Test3(t *testing.T) {
	assert.Equal(t, m.GetMust("test3").want, get("test3", m.GetMust("test3").code))
}
func Test4(t *testing.T) {
	assert.Equal(t, m.GetMust("test4").want, get("test4", m.GetMust("test4").code))
}
func Test5(t *testing.T) {
	assert.Equal(t, m.GetMust("test5").want, get("test5", m.GetMust("test5").code))
}
func Test6(t *testing.T) {
	assert.Equal(t, m.GetMust("test6").want, get("test6", m.GetMust("test6").code))
}
func Test7(t *testing.T) {
	assert.Equal(t, m.GetMust("test7").want, get("test7", m.GetMust("test7").code))
}
func Test8(t *testing.T) {
	t.Skipf(`
不确定是否应该删除
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, nil, err
	}
`)
	assert.Equal(t, m.GetMust("test8").want, get("test8", m.GetMust("test8").code))
}
func Test9(t *testing.T) {
	assert.Equal(t, m.GetMust("test9").want, get("test9", m.GetMust("test9").code))
}
func Test10(t *testing.T) {
	t.Skip("todo bug")
	assert.Equal(t, m.GetMust("test10").want, get("test10", m.GetMust("test10").code))
}

func get(path, text string) string {
	join := filepath.Join(os.TempDir(), path+".go")
	fix := filepath.Join(os.TempDir(), path+"_fixed.go")
	mylog.Warning("original code file path", join+":1")
	stream.WriteGoFile(join, text) //写入文件只是为了让goland检查原始代码的语法和直观的行号
	ret := ""
	mylog.Call(func() {
		fileSet := token.NewFileSet()
		file := mylog.Check2(parser.ParseFile(fileSet, fix, text, parser.ParseComments))
		ret = handle(fileSet, file, text)
		mylog.WriteGoFile(fix, ret) //写入文件只是为了让goland检查返回代码的语法和直观的行号
	})
	return string(mylog.Check2(format.Source([]byte(ret)))) //handle中WriteGoFile已经执行格式化
}

type testData struct {
	code string
	want string
}

var m = safemap.NewOrdered[string, testData](func(yield func(string, testData) bool) {
	yield("test1", testData{
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
	"github.com/ddkwork/golibrary/mylog"
	log "github.com/sirupsen/logrus"
)

func main() {

	mylog.Check(backendConn.Close())
}
`,
	})
	yield("test2", testData{
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
	"github.com/ddkwork/golibrary/mylog"
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
	yield("test3", testData{
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
	"github.com/ddkwork/golibrary/mylog"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

func main() {
	backendConn := mylog.Check2(net.DialTimeout("tcp", forwardTarget, 5*time.Second))

}
`,
	})
	yield("test4", testData{
		code: `package main

import (
	"github.com/ddkwork/golibrary/mylog"
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
	"github.com/ddkwork/golibrary/mylog"
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
	yield("test5", testData{
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
	"github.com/ddkwork/golibrary/mylog"
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
	yield("test6", testData{
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
	"github.com/ddkwork/golibrary/mylog"
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
	yield("test7", testData{
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
	"github.com/ddkwork/golibrary/mylog"
)

func mian() {
	provider := mylog.Check2(oidc.NewProvider(c.Ctx, c.ProviderURL))

}
`,
	})
	yield("test8", testData{
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
	"github.com/ddkwork/golibrary/mylog"
	log "github.com/sirupsen/logrus"
)

func main() {
	mylog.Check(jsonx.Open(&token, tf))
}
`,
	})
	yield("test9", testData{
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
	"github.com/ddkwork/golibrary/mylog"
	log "github.com/sirupsen/logrus"
)

func main() {
	newToken := mylog.Check2(tokenSource.Token())

	userInfo := mylog.Check2(provider.UserInfo(c.Ctx, tokenSource))

}
`,
	})
	yield("test10", testData{
		code: `package tmp

import (
	"github.com/ddkwork/golibrary/mylog"
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
	"github.com/ddkwork/golibrary/mylog"
)

func main() {
	panic(bug())
}

func bug() error {
	mylog.Check2(w.Write([]byte("Hello, 世界")))
	mylog.Check(apkw.Close())
	defer mylog.Check(io.Close(nil))
	defer func() { mylog.Check(io.Close(nil)) }()
}
`,
	})
})
