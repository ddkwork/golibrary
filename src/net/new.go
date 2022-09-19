package net

import (
	"github.com/ddkwork/golibrary/src/net/clientAndServer/client"
	"github.com/ddkwork/golibrary/src/net/clientAndServer/server"
	"github.com/ddkwork/golibrary/src/net/transport"
)

type (
	xxInterface interface {
		TransPort() transport.Interface
		TcpClient() client.Interface
		TcpServer() server.Interface
	}
	xxobject struct {
		transPort transport.Interface
		tcpClient client.Interface
		tcpServer server.Interface
	}
)

func (x xxobject) TransPort() transport.Interface {
	return x.transPort
}

func (x xxobject) TcpClient() client.Interface {
	return client.Default
}

func (x xxobject) TcpServer() server.Interface {
	return server.Default
}

func xxNew() xxInterface {
	return &xxobject{
		transPort: transport.New(),
		tcpClient: client.Default,
		tcpServer: server.Default,
	}
}
