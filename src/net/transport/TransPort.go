package transport

import (
	"github.com/ddkwork/golibrary/src/net/transport/internal/tcp"
	"github.com/ddkwork/golibrary/src/net/transport/internal/udp"
)

type (
	Interface interface {
		Tcp() tcp.Interface
		Udp() udp.Interface
	}
	object struct {
		tcp tcp.Interface
		udp udp.Interface
	}
)

func (o *object) Tcp() tcp.Interface {
	return o.tcp
}

func (o *object) Udp() udp.Interface {
	return o.udp
}

func New() Interface {
	return &object{
		tcp: tcp.New(),
		udp: udp.New(),
	}
}
