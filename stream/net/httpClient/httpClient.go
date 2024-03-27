// <nilaway struct enable>
package httpClient

import (
	"bytes"
	"context"
	"crypto/tls"
	_ "embed"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hupe1980/socks"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
	"github.com/gorilla/websocket"
)

type (
	Object struct {
		debug bool
		*http.Client
		cookiejar   *cookiejar.Jar
		form        url.Values
		requestBody []byte
		method      string
		requestUrl  string
		path        string
		head        http.Header
		stopCode    int
		responseBuf []byte
	}
)

const Localhost = "127.0.0.1"

var Default = New()

func New() *Object {
	o := &Object{
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
		requestUrl:  "",
		path:        "",
		head:        http.Header{},
		stopCode:    http.StatusOK,
		responseBuf: make([]byte, 0),
	}
	jar, err := cookiejar.New(nil)
	if !mylog.Error(err) {
		return o
	}
	o.cookiejar = jar
	return o
}

func (o *Object) SetDebug(debug bool) *Object                { o.debug = debug; return o }
func (o *Object) ResponseBufStream() *stream.Stream          { return stream.New(o.responseBuf) }
func (o *Object) Get() *Object                               { o.method = http.MethodGet; return o }
func (o *Object) Post() *Object                              { o.method = http.MethodPost; return o }
func (o *Object) Url(RequestUrl string) *Object              { o.requestUrl = RequestUrl; return o }
func (o *Object) SetPath(path string) *Object                { o.path = path; return o }
func (o *Object) StopCode(stopCode int) *Object              { o.stopCode = stopCode; return o }
func (o *Object) BaseURL() string                            { return o.requestUrl }
func (o *Object) Cookiejar() *cookiejar.Jar                  { return o.cookiejar }
func (o *Object) SetForm(form url.Values) *Object            { o.form = form; return o }
func (o *Object) Body(requestBody []byte) *Object            { o.requestBody = requestBody; return o }
func (o *Object) BodyStream(s *stream.Stream) *Object        { o.requestBody = s.Bytes(); return o }
func (o *Object) CreatNewClient(client *http.Client) *Object { o.Client = client; return o }
func (o *Object) ProxyHttp(s string) *Object                 { return o.SetProxy(HttpLayer, s) }
func (o *Object) ProxyHttps(s string) *Object                { return o.SetProxy(HttpsLayer, s) }
func (o *Object) ProxySocket5Layer(s string) *Object         { return o.SetProxy(Socket5Layer, s) }
func (o *Object) ProxySocket4Layer(s string) *Object         { return o.SetProxy(Socket4Layer, s) }
func (o *Object) ProxyWebSocketLayer(s string) *Object       { return o.SetProxy(WebSocketLayer, s) }
func (o *Object) ProxyWebsocketTlsLayer(s string) *Object    { return o.SetProxy(WebsocketTlsLayer, s) }
func (o *Object) CheckProtocol(protocol string, port string) bool {
	return false
} // todo

func (o *Object) Request() (ok bool) {
	fnReader := func() io.Reader {
		if o.requestBody != nil {
			return bytes.NewReader(o.requestBody)
		}
		return strings.NewReader(o.form.Encode())
	}
	request, err := http.NewRequest(o.method, o.requestUrl+o.path, fnReader())
	if !mylog.Error(err) {
		return // I really can't find a demo where the request is still nil after the return here, can you give me a demo to convince the request that it is nil
	}
	//if request == nil {
	//	mylog.Error("request == nil")
	//	return
	//}
	request.Close = false
	// Request.Header.Add("Connection", "close")
	request.Header = o.head

	if o.debug {
		println(mylog.DumpRequest(request, true))
	}
	response, err := o.Client.Do(request)
	if !mylog.Error(err) {
		return
	}
	//if response == nil {
	//	mylog.Error("response == nil")
	//	return
	//}
	defer func() {
		if o.debug {
			println(mylog.DumpResponse(response, true))
		}
		mylog.Error(response.Body.Close())
	}()
	switch response.StatusCode {
	case http.StatusOK, o.stopCode:
		all, err := io.ReadAll(response.Body) // todo 外部判断gzip，是否可以提进来，就不用重复劳动了
		if mylog.Error(err) {
			o.responseBuf = all
		}
		return true
	default:
		return mylog.Error(errors.New(response.Status + " != StopCode " + strconv.Itoa(o.stopCode)))
	}
}

func (o *Object) SetJsonHead(header map[string]string) *Object {
	o.SetHead(header)
	o.head.Set("content-type", "application/json")
	return o
}

func (o *Object) SetHead(header map[string]string) *Object {
	for k, v := range header {
		o.head.Set(k, v)
	}
	o.head.Set("user-agent", UserAgentRandom)
	return o
}

func (o *Object) HasCookieInJar(jar *cookiejar.Jar, cookieName, Host string) (ok bool) {
	URL, err := url.Parse(Host)
	if !mylog.Error(err) {
		return
	}
	for _, v := range jar.Cookies(URL) {
		if v.Name == cookieName {
			mylog.Success(" find cookie by name", v)
			return true
		}
	}
	return
}

// SetProxyEx todo add auth and cert
func (o *Object) SetProxyEx(layer Layer, hostPort string) *Object {
	o.SetProxy(layer, hostPort)
	return o
}

func (o *Object) SetProxy(layer Layer, hostPort string) *Object {
	host, port, err := net.SplitHostPort(hostPort)
	if !mylog.Error(err) {
		return nil
	}
	var ObjTransport struct {
		Transport    *http.Transport
		dialFunc     func(network, addr string) (net.Conn, error)
		proxyURLFunc func(*http.Request) (*url.URL, error)
	}
	switch layer {
	case Socket4Layer:
		o.Client.Transport = &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				d := socks.NewSocks4Dialer("tcp", net.JoinHostPort(host, port)) // todo test
				return d.DialContext(ctx, network, addr)
			},
		}
	case Socket5Layer:
		o.Client.Transport = &http.Transport{
			Proxy: func(request *http.Request) (*url.URL, error) {
				return url.Parse("socks5://" + net.JoinHostPort(host, port))
			},
		}
	case HttpLayer, HttpsLayer:
		URL, err := url.Parse(HttpLayer.String() + "://" + net.JoinHostPort(host, port))
		if !mylog.Error(err) {
			return nil
		}
		ObjTransport.proxyURLFunc = http.ProxyURL(URL)
		ObjTransport.dialFunc = (&net.Dialer{
			Timeout:       6 * time.Second,
			Deadline:      time.Now().Add(20 * time.Second),
			LocalAddr:     nil,
			DualStack:     false,
			FallbackDelay: 0,
			KeepAlive:     0,
			Resolver:      nil,
			Cancel:        nil,
			Control:       nil,
		}).Dial
		o.Client.Transport = &http.Transport{
			Proxy:                  ObjTransport.proxyURLFunc,
			DialContext:            nil,
			Dial:                   ObjTransport.dialFunc,
			DialTLSContext:         nil,
			DialTLS:                nil,
			TLSClientConfig:        &tls.Config{InsecureSkipVerify: true}, // todo load cert
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
	case WebsocketTlsLayer:
		endpointURL := "wss://localhost:12345" // todo
		// proxyURL := "http://localhost:6666"
		proxyURL := "http://" + net.JoinHostPort(host, port)

		surl, err := url.Parse(proxyURL)
		if !mylog.Error(err) {
			return nil
		}
		dialer := websocket.Dialer{
			NetDial:           nil,
			NetDialContext:    nil,
			NetDialTLSContext: nil,
			Proxy:             http.ProxyURL(surl),
			// TLSClientConfig:   ca.MitmCfg.NewTlsConfigForHost(Localhost),
			// TLSClientConfig:   ca.MitmCfg.NewTlsConfigForHost("localhost"),
			HandshakeTimeout:  0,
			ReadBufferSize:    bufferSize,
			WriteBufferSize:   bufferSize,
			WriteBufferPool:   nil,
			Subprotocols:      []string{"p1"},
			EnableCompression: false,
			Jar:               nil,
		}
		dialer.Dial(endpointURL, nil)

	}
	return o
}

const bufferSize = 1024 * 8

//go:embed logevent.bin
var LogeventBuf []byte

func MockProtoBufPacket(proxyPort string) {
	c := New()
	c.SetDebug(true)
	c.SetProxy(HttpLayer, net.JoinHostPort(Localhost, proxyPort))
	// https://cloud.tencent.com/developer/article/1624700
	header := map[string]string{
		//"Content-Type": "application/x-google-protobuf",
		"Content-Type": "application/x-protobuf",
		//"Content-Type": "application/x-protobuffer",
	}
	// c.Url("https://www.baidu.com").ProxyHttp(":8080").BodyStream(stream.NewHexDump(TestPbBuf)).Post().SetHead(header).Request()
	c.Url("https://www.baidu.com").BodyStream(stream.New(LogeventBuf)).Post().SetHead(header).Request()
	// c.Url("https://www.baidu.com").ProxyHttp(":8080").BodyStream(stream.New(logeventBuf)).Post().SetHead(header).Request()
	// c.Url("https://www.baidu.com").ProxyHttp(":8080").Body(session.GooglePb()).Post().SetHead(m).Request()
}
