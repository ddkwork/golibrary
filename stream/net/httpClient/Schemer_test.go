package httpClient

import (
	"testing"

	"github.com/ddkwork/golibrary/mylog"
	"mvdan.cc/gofumpt/format"

	"github.com/stretchr/testify/assert"

	"github.com/ddkwork/golibrary/stream"
)

func TestLayer_AssertKind(t *testing.T) {
	assert.Equal(t, TcpKind, InvalidKind.AssertKind("tcp"))
	assert.Equal(t, HttpsKind, InvalidKind.AssertKind("https"))
}

func TestGeneratedFile_Iota(t *testing.T) {
	stream.NewGeneratedFile().Enum("Schemer", []string{
		"Invalid",
		"Http",
		"Https",
		"Socket4",
		"Socket5",
		"WebSocket",
		"WebsocketTls",
		"Tcp",
		"TcpTls",
		"Udp",
		"Kcp",
		"Pipe",
		"Quic",
		"Rpc",
		"Ssh",
	}, nil)
	b := stream.NewBuffer("Schemer_enum_gen.go")
	b.WriteStringLn(expansions)
	f := mylog.Check2(format.Source(b.Bytes(), format.Options{}))
	stream.WriteTruncate("Schemer_enum_gen.go", f)
}

var expansions = `

func (l SchemerKind) IsStream() bool {
	switch l {
	case InvalidKind, HttpKind, HttpsKind:
		return false
	default:
		return true

	}
}

func (l SchemerKind) IsContainer() bool {
	return l.IsStream()
}

func (l SchemerKind) Containers() []SchemerKind {
	return []SchemerKind{
		WebSocketKind,
		KcpKind,
		PipeKind,
		QuicKind,
		RpcKind,
		Socket4Kind,
		Socket5Kind,
		SshKind,
		TcpKind,
		TcpTlsKind,
		UdpKind,
		WebsocketTlsKind,
	}
}
`
