package httpClient

import (
	"strings"

	"golang.org/x/exp/constraints"
)

// Code generated by GeneratedFile enum - DO NOT EDIT.

type SchemerKind byte

const (
	InvalidKind SchemerKind = iota
	HttpKind
	HttpsKind
	Socket4Kind
	Socket5Kind
	WebSocketKind
	WebsocketTlsKind
	TcpKind
	TcpTlsKind
	UdpKind
	KcpKind
	PipeKind
	QuicKind
	RpcKind
	SshKind
	InvalidSchemerKind
)

func ConvertInteger2SchemerKind[T constraints.Integer](v T) SchemerKind {
	return SchemerKind(v)
}

func (k SchemerKind) AssertKind(kinds string) SchemerKind {
	for _, kind := range k.Kinds() {
		if strings.ToLower(kinds) == strings.ToLower(kind.String()) {
			return kind
		}
	}
	return InvalidSchemerKind
}

func (k SchemerKind) String() string {
	switch k {
	case InvalidKind:
		return "Invalid"
	case HttpKind:
		return "Http"
	case HttpsKind:
		return "Https"
	case Socket4Kind:
		return "Socket4"
	case Socket5Kind:
		return "Socket5"
	case WebSocketKind:
		return "WebSocket"
	case WebsocketTlsKind:
		return "WebsocketTls"
	case TcpKind:
		return "Tcp"
	case TcpTlsKind:
		return "TcpTls"
	case UdpKind:
		return "Udp"
	case KcpKind:
		return "Kcp"
	case PipeKind:
		return "Pipe"
	case QuicKind:
		return "Quic"
	case RpcKind:
		return "Rpc"
	case SshKind:
		return "Ssh"
	default:
		return "InvalidSchemerKind"
	}
}

func (k SchemerKind) Keys() []string {
	return []string{
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
		"InvalidSchemerKind",
	}
}

func (k SchemerKind) Kinds() []SchemerKind {
	return []SchemerKind{
		InvalidKind,
		HttpKind,
		HttpsKind,
		Socket4Kind,
		Socket5Kind,
		WebSocketKind,
		WebsocketTlsKind,
		TcpKind,
		TcpTlsKind,
		UdpKind,
		KcpKind,
		PipeKind,
		QuicKind,
		RpcKind,
		SshKind,
		InvalidSchemerKind,
	}
}

func (k SchemerKind) SvgFileName() string {
	switch k {
	case InvalidKind:
		return "Invalid"
	case HttpKind:
		return "Http"
	case HttpsKind:
		return "Https"
	case Socket4Kind:
		return "Socket4"
	case Socket5Kind:
		return "Socket5"
	case WebSocketKind:
		return "WebSocket"
	case WebsocketTlsKind:
		return "WebsocketTls"
	case TcpKind:
		return "Tcp"
	case TcpTlsKind:
		return "TcpTls"
	case UdpKind:
		return "Udp"
	case KcpKind:
		return "Kcp"
	case PipeKind:
		return "Pipe"
	case QuicKind:
		return "Quic"
	case RpcKind:
		return "Rpc"
	case SshKind:
		return "Ssh"
	default:
		return "InvalidSchemerKind"
	}
}

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
