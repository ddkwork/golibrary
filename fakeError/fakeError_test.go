package fakeError

import (
	"github.com/ddkwork/golibrary/assert"
	"github.com/ddkwork/golibrary/safemap"
	"testing"
)

func Test1(t *testing.T) {
	assert.Equal(t, m.GetMust("test1").expected, testHandle("test1", m.GetMust("test1").code))
}
func Test2(t *testing.T) {
	assert.Equal(t, m.GetMust("test2").expected, testHandle("test2", m.GetMust("test2").code))
}
func Test3(t *testing.T) {
	assert.Equal(t, m.GetMust("test3").expected, testHandle("test3", m.GetMust("test3").code))
}
func Test4(t *testing.T) {
	assert.Equal(t, m.GetMust("test4").expected, testHandle("test4", m.GetMust("test4").code))
}
func Test5(t *testing.T) {
	assert.Equal(t, m.GetMust("test5").expected, testHandle("test5", m.GetMust("test5").code))
}
func Test6(t *testing.T) {
	assert.Equal(t, m.GetMust("test6").expected, testHandle("test6", m.GetMust("test6").code))
}
func Test7(t *testing.T) {
	assert.Equal(t, m.GetMust("test7").expected, testHandle("test7", m.GetMust("test7").code))
}
func Test8(t *testing.T) {
	t.Skip()
	assert.Equal(t, m.GetMust("test8").expected, testHandle("test8", m.GetMust("test8").code))
}
func Test9(t *testing.T) {
	assert.Equal(t, m.GetMust("test9").expected, testHandle("test9", m.GetMust("test9").code))
}

type testData struct {
	code     string
	expected string
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
		expected: `package tmp

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
		expected: `package main

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
		expected: `package main

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
		expected: `package main

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
		expected: `package main

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
		expected: `package main

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
		expected: `package main

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
		expected: `package tmp

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
		expected: `package tmp

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
})
