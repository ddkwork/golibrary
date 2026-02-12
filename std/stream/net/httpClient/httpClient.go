package httpClient

import (
	"bytes"
	"context"
	"crypto/tls"
	_ "embed"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/golibrary/std/stream"
)

type (
	Client struct {
		debug bool
		*http.Client
		*http.Response
		cookiejar       *cookiejar.Jar
		form            url.Values
		requestBody     []byte
		method          string
		url             string
		head            http.Header
		stopCode        int
		BadRequestCount int
		*stream.Buffer
	}
)

const Localhost = "127.0.0.1"

func New() *Client {
	o := &Client{
		debug: false,
		Client: &http.Client{
			Transport:     nil,
			CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
			Jar:           nil,
			Timeout:       60 * time.Second,
		},
		cookiejar:   nil,
		form:        make(url.Values),
		requestBody: make([]byte, 0),
		method:      "",
		url:         "",
		head:        http.Header{},
		stopCode:    http.StatusOK,
	}

	o.cookiejar = mylog.Check2(cookiejar.New(nil))
	return o
}

func (c *Client) IsBadRequest() bool          { return c.BadRequestCount > 0 }
func (c *Client) SetDebug(debug bool) *Client { c.debug = debug; return c }
func (c *Client) Get(url string) *Client {
	c.url = url
	c.method = http.MethodGet
	return c
}

func (c *Client) Post(url string) *Client {
	c.url = url
	c.method = http.MethodPost
	return c
}
func (c *Client) StopCode(stopCode int) *Client           { c.stopCode = stopCode; return c }
func (c *Client) Cookiejar() *cookiejar.Jar               { return c.cookiejar }
func (c *Client) SetForm(form url.Values) *Client         { c.form = form; return c }
func (c *Client) Body(requestBody []byte) *Client         { c.requestBody = requestBody; return c }
func (c *Client) BodyStream(s *stream.Buffer) *Client     { c.requestBody = s.Bytes(); return c }
func (c *Client) ProxyHttp(s string) *Client              { return c.SetProxy(HttpType, s) }
func (c *Client) ProxyHttps(s string) *Client             { return c.SetProxy(HttpsType, s) }
func (c *Client) ProxySocket5Layer(s string) *Client      { return c.SetProxy(Socket5Type, s) }
func (c *Client) ProxySocket4Layer(s string) *Client      { return c.SetProxy(Socket4Type, s) }
func (c *Client) ProxyWebSocketLayer(s string) *Client    { return c.SetProxy(WebSocketType, s) }
func (c *Client) ProxyWebsocketTlsLayer(s string) *Client { return c.SetProxy(WebsocketTlsType, s) }
func (c *Client) CheckProtocol(protocol string, port string) bool {
	return false
}

func (c *Client) Request() *Client {
	mylog.Call(func() { c.request() })
	return c
}

func (c *Client) request() *Client {
	fnReader := func() io.Reader {
		if c.requestBody != nil {
			return bytes.NewReader(c.requestBody)
		}
		return strings.NewReader(c.form.Encode())
	}
	request := mylog.Check2(http.NewRequest(c.method, c.url, fnReader()))
	request.Close = false

	request.Header = c.head
	if c.debug {
		mylog.Request(request, true)
	}
	response := mylog.Check2(c.Client.Do(request))
	c.Response = response
	defer func() {
		if c.debug {
			mylog.Response(response, true)
		}
		mylog.Check(response.Body.Close())
	}()
	switch response.StatusCode {
	case http.StatusOK, c.stopCode:
		body, backBody := CloneBody(response.Body)
		response.Body = body
		b := mylog.Check2(io.ReadAll(backBody))
		c.Buffer = stream.NewBuffer(b)
		if request.Header.Get("Content-Encoding") == "gzip" {
			c.Buffer = stream.ReaderGzip(b)
		}
	default:
		c.BadRequestCount++
		mylog.Check(response.Status + " != StopCode " + strconv.Itoa(c.stopCode))
	}
	return c
}

func CloneBody(b io.ReadCloser) (body, backBody io.ReadCloser) {
	mylog.CheckNil(b)
	if b == http.NoBody || b == nil {
		return http.NoBody, http.NoBody
	}
	var buf bytes.Buffer
	mylog.Check2(buf.ReadFrom(b))
	mylog.Check(b.Close())
	return io.NopCloser(&buf), io.NopCloser(bytes.NewReader(buf.Bytes()))
}

func (c *Client) SetJsonHead() *Client {
	c.head.Set("content-type", "application/json")
	return c
}

func (c *Client) SetHead(header map[string]string) *Client {
	for k, v := range header {
		c.head.Set(k, v)
	}
	c.head.Set("user-agent", stream.RandomAnySlice(UserAgents))
	return c
}

func (c *Client) HasCookieInJar(jar *cookiejar.Jar, cookieName, Host string) (ok bool) {
	URL := mylog.Check2(url.Parse(Host))
	for _, v := range jar.Cookies(URL) {
		if v.Name == cookieName {
			mylog.Success(" find cookie by name", v)
			return true
		}
	}
	return
}

func (c *Client) SetProxyEx(layer SchemerType, hostPort string) *Client {
	c.SetProxy(layer, hostPort)
	return c
}

func (c *Client) SetProxy(layer SchemerType, hostPort string) *Client {
	host, port := mylog.Check3(net.SplitHostPort(hostPort))
	var Transport struct {
		Transport    *http.Transport
		dialContext  func(ctx context.Context, network, addr string) (net.Conn, error)
		proxyURLFunc func(*http.Request) (*url.URL, error)
	}
	switch layer {
	case Socket4Type:
		// c.Client.Transport = &http.Transport{
		//	DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
		//		d := socks.NewSocks4Dialer("tcp", net.JoinHostPort(host, port))
		//		return d.DialContext(ctx, network, addr)
		//	},
		// }
	case Socket5Type:
		c.Client.Transport = &http.Transport{
			Proxy: func(request *http.Request) (*url.URL, error) {
				return url.Parse("socks5://" + net.JoinHostPort(host, port))
			},
		}
	case HttpType, HttpsType:
		URL := mylog.Check2(url.Parse(HttpType.String() + "://" + net.JoinHostPort(host, port)))
		Transport.proxyURLFunc = http.ProxyURL(URL)
		Transport.dialContext = (&net.Dialer{
			Timeout:   6 * time.Second,
			Deadline:  time.Now().Add(20 * time.Second),
			LocalAddr: nil,

			FallbackDelay: 0,
			KeepAlive:     0,
			Resolver:      nil,

			Control:        nil,
			ControlContext: nil,
		}).DialContext
		c.Client.Transport = &http.Transport{
			Proxy:       Transport.proxyURLFunc,
			DialContext: Transport.dialContext,

			DialTLSContext: nil,

			TLSClientConfig:        &tls.Config{InsecureSkipVerify: true},
			TLSHandshakeTimeout:    0,
			DisableKeepAlives:      true,
			DisableCompression:     false,
			MaxIdleConns:           10,
			MaxIdleConnsPerHost:    10,
			MaxConnsPerHost:        0,
			IdleConnTimeout:        0,
			ResponseHeaderTimeout:  0,
			ExpectContinueTimeout:  0,
			TLSNextProto:           nil,
			ProxyConnectHeader:     nil,
			GetProxyConnectHeader:  nil,
			MaxResponseHeaderBytes: 0,
			WriteBufferSize:        bufferSize,
			ReadBufferSize:         bufferSize,
			ForceAttemptHTTP2:      false,
		}
	case WebsocketTlsType:
		// endpointURL := "wss://localhost:12345"
		//
		// proxyURL := "http://" + net.JoinHostPort(host, port)
		//
		// surl := mylog.Check2(url.Parse(proxyURL))
		// dialer := websocket.Dialer{
		//	NetDial:           nil,
		//	NetDialContext:    nil,
		//	NetDialTLSContext: nil,
		//	Proxy:             http.ProxyURL(surl), // todo need tls can be work
		//
		//	HandshakeTimeout:  0,
		//	ReadBufferSize:    bufferSize,
		//	WriteBufferSize:   bufferSize,
		//	WriteBufferPool:   nil,
		//	Subprotocols:      []string{"p1"},
		//	EnableCompression: false,
		//	Jar:               nil,
		// }
		// mylog.Check3(dialer.Dial(endpointURL, nil))
	default:
		mylog.Check("unhandled default case")
	}
	return c
}

const bufferSize = 1024 * 8

//go:embed logevent.bin
var LogeventBuf []byte

func MockProtoBufPacket(proxyPort string) {
	c := New()
	c.SetDebug(true)
	c.SetProxy(HttpType, net.JoinHostPort(Localhost, proxyPort))

	header := map[string]string{
		"Content-Type": "application/x-protobuf",
	}

	c.Post("https://www.baidu.com").BodyStream(stream.NewBuffer(LogeventBuf)).SetHead(header).Request()
}

var (
	UserAgentRandom = stream.RandomAnySlice(UserAgents)
	UserAgents      = []string{
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.78 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.78 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.78 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.78 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.78 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.78 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.78 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.78 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.78 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.78 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.79 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.79 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.79 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.79 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.79 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.79 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.79 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.79 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.79 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.79 Safari/537.36",
	}
)
