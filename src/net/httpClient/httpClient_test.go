package httpClient_test

import (
	"github.com/ddkwork/golibrary/src/net/httpClient"
	"github.com/ddkwork/golibrary/src/net/httpClient/session"
	"net"
	"testing"
)

func TestPb2(t *testing.T) {
	c := httpClient.New()
	//https://cloud.tencent.com/developer/article/1624700
	m := map[string]string{
		"Content-Type": "application/x-google-protobuf",
		//"Content-Type": "application/x-protobuf",
		//"Content-Type": "application/x-protobuffer",
	}
	c.Url("https://www.baidu.com").ProxyHttp(":8080").Body(session.GooglePb()).Post().SetHead(m).Request()
}

func TestName(t *testing.T) {
	host, port, err := net.SplitHostPort(":8080")
	if err != nil {
		println(err.Error())
		return
	}
	println(host)
	println(port)
}
