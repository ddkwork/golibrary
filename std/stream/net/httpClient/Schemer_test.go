package httpClient

import (
	"go/format"
	"testing"

	"github.com/ddkwork/golibrary/std/safemap"

	"github.com/ddkwork/golibrary/std/assert"
	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/golibrary/std/stream"
)

func TestLayer_AssertKind(t *testing.T) {
	assert.Equal(t, TcpType, HttpType.AssertBy("tcp"))
	assert.Equal(t, HttpsType, HttpType.AssertBy("https"))
}

func TestGeneratedFile_Iota(t *testing.T) {
	m := safemap.NewOrdered[string, string](func(yield func(string, string) bool) {
	})
	m.Set("Http", "Http")
	m.Set("Https", "Https")
	m.Set("Socket4", "Socket4")
	m.Set("Socket5", "Socket5")
	m.Set("WebSocket", "WebSocket")
	m.Set("WebsocketTls", "WebsocketTls")
	m.Set("Tcp", "Tcp")
	m.Set("TcpTls", "TcpTls")
	m.Set("Udp", "Udp")
	m.Set("Kcp", "Kcp")
	m.Set("Pipe", "Pipe")
	m.Set("Quic", "Quic")
	m.Set("Rpc", "Rpc")
	m.Set("Ssh", "Ssh")
	stream.NewGeneratedFile().EnumTypes("Schemer", m)
	b := stream.NewBuffer("Schemer_types_gen.go")
	b.WriteStringLn(expansions)
	f := mylog.Check2(format.Source(b.Bytes()))
	stream.WriteTruncate("Schemer_types_gen.go", f)
}

var expansions = `

func (l SchemerType) IsStream() bool {
	switch l {
	case HttpType, HttpsType:
		return false
	default:
		return true

	}
}

func (l SchemerType) IsContainer() bool {
	return l.IsStream()
}

func (l SchemerType) Containers() []SchemerType {
	return []SchemerType{
		WebSocketType,
		KcpType,
		PipeType,
		QuicType,
		RpcType,
		Socket4Type,
		Socket5Type,
		SshType,
		TcpType,
		TcpTlsType,
		UdpType,
		WebsocketTlsType,
	}
}
`
