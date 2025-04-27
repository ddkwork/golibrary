package fakeError

import (
	"github.com/ddkwork/golibrary/assert"
	"github.com/ddkwork/golibrary/safemap"
	"testing"
)

func Test1(t *testing.T) {
	assert.Equal(t, m.GetMust("test1").expected, fakeErrorTest(m.GetMust("test1").code))
}
func Test2(t *testing.T) {
	assert.Equal(t, m.GetMust("test2").expected, fakeErrorTest(m.GetMust("test2").code))
}
func Test3(t *testing.T) {
	assert.Equal(t, m.GetMust("test3").expected, fakeErrorTest(m.GetMust("test3").code))
}
func Test4(t *testing.T) {
	assert.Equal(t, m.GetMust("test4").expected, fakeErrorTest(m.GetMust("test4").code))
}
func Test5(t *testing.T) {
	assert.Equal(t, m.GetMust("test5").expected, fakeErrorTest(m.GetMust("test5").code))
}
func Test6(t *testing.T) {
	assert.Equal(t, m.GetMust("test6").expected, fakeErrorTest(m.GetMust("test6").code))
}

func Test_ApplyEdit(t *testing.T) {
	replacements := []Edit{
		{
			Start: 90,
			End:   125,
			Line:  9,
			New:   "",
		},
		{
			Start: 127,
			End:   158,
			Line:  12,
			New:   "",
		},
		{
			Start: 160,
			End:   223,
			Line:  15,
			New:   "mylog.Check(backendConn.Close())",
		},
	}
	assert.Equal(t, m.GetMust("test1").expected, Apply(m.GetMust("test1").code, replacements))
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
})
