package httpClient

import (
	"bytes"
	"crypto/tls"
	"errors"

	"github.com/ddkwork/golibrary/mylog"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type (
	Object struct {
		client      *http.Client
		cookiejar   *cookiejar.Jar
		form        url.Values
		requestBody []byte
		method      string
		requestUrl  string
		path        string
		head        map[string]string
		stopCode    int
		responseBuf []byte
		error       error
	}
)

var Default = New()

func New() *Object {
	jar, err := cookiejar.New(nil)
	if !mylog.Error(err) {
		return nil
	}
	return &Object{
		client: &http.Client{
			Transport:     nil,
			CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
			Jar:           nil,
			Timeout:       60 * time.Second,
		},
		cookiejar:  jar,
		form:       nil,
		method:     "",
		requestUrl: "",
		path:       "",
		head:       nil,
		stopCode:   http.StatusOK,
	}
}

func (o *Object) Request() (ok bool) {
	var (
		request  = new(http.Request)
		response = new(http.Response)
	)

	var body io.Reader
	if o.requestBody != nil {
		body = bytes.NewReader(o.requestBody)
	} else {
		body = strings.NewReader(o.form.Encode())
	}

	request, o.error = http.NewRequest(o.method, o.requestUrl+o.path, body)
	if !mylog.Error(o.error) {
		return
	}
	request.Close = true //强制短链接
	//Request.Header.Add("Connection", "close")
	for k, v := range o.head {
		request.Header.Set(k, v)
	}
	response, o.error = o.client.Do(request)
	if !mylog.Error(o.error) {
		return
	}
	defer func() {
		if response == nil {
			o.error = errors.New("response == nil")
			mylog.Error(o.error)
			return
		}
		mylog.Error(response.Body.Close())
	}()
	switch response.StatusCode {
	case http.StatusOK, o.stopCode:
		o.responseBuf, o.error = io.ReadAll(response.Body) //todo 外部判断gzip，是否可以提进来，就不用重复劳动了
		return mylog.Error(o.error)
	default:
		o.error = errors.New(response.Status + " != StopCode " + strconv.Itoa(o.stopCode))
		return mylog.Error(o.error)
	}
}
func (o *Object) ResponseBuf() []byte { return o.responseBuf }
func (o *Object) Error() error        { return o.error }

func (o *Object) Get() *Object {
	o.method = http.MethodGet
	return o
}
func (o *Object) Post() *Object {
	o.method = http.MethodPost
	return o
}
func (o *Object) Url(RequestUrl string) *Object {
	o.requestUrl = RequestUrl
	return o
}
func (o *Object) SetPath(path string) *Object {
	o.path = path
	return o
}
func (o *Object) SetHead(head map[string]string) *Object {
	o.head = head
	return o
}
func (o *Object) StopCode(stopCode int) *Object {
	o.stopCode = stopCode
	return o
}
func (o *Object) BaseURL() string           { return o.requestUrl }
func (o *Object) Cookiejar() *cookiejar.Jar { return o.cookiejar }
func (o *Object) SetForm(form url.Values) *Object {
	o.form = form
	return o
}
func (o *Object) Body(requestBody []byte) *Object {
	o.requestBody = requestBody
	return o
}
func (o *Object) Client() *http.Client { return o.client }
func (o *Object) CreatNewClient(client *http.Client) *Object {
	o.client = client
	return o
}
func (o *Object) hasCookieInJar(jar *cookiejar.Jar, cookieName, Host string) (ok bool) {
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
func (o *Object) ProxyHttp(hostPort string) *Object   { return o.setProxy("http", hostPort) }
func (o *Object) ProxyHttps(hostPort string) *Object  { return o.setProxy("https", hostPort) }
func (o *Object) ProxySocks4(hostPort string) *Object { return o.setProxy("socks4", hostPort) }
func (o *Object) ProxySocks5(hostPort string) *Object { return o.setProxy("socks5", hostPort) }
func (o *Object) setProxy(protocol, hostPort string) *Object {
	host, port, err := net.SplitHostPort(hostPort)
	if err != nil {
		return nil
	}
	var ObjTransport struct {
		Transport    *http.Transport
		dialFunc     func(network, addr string) (net.Conn, error)
		proxyURLFunc func(*http.Request) (*url.URL, error)
	}
	switch protocol {
	case ProtoName.Socks4(), ProtoName.Socks5():
		ObjTransport.dialFunc = SDial(protocol + "://" + hostPort + "?timeout=20s")
	case ProtoName.Http(), ProtoName.Https():
		URL, err := url.Parse(ProtoName.Http() + "://" + host + ":" + port)
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
	}
	o.client.Transport = &http.Transport{
		Proxy:                  ObjTransport.proxyURLFunc,
		DialContext:            nil,
		Dial:                   ObjTransport.dialFunc,
		DialTLSContext:         nil,
		DialTLS:                nil,
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
		WriteBufferSize:        0,
		ReadBufferSize:         0,
		ForceAttemptHTTP2:      false,
	}
	return o
}
