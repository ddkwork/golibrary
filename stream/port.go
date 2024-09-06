package stream

import (
	"fmt"
	"github.com/ddkwork/golibrary/mylog"
	"net"
)

// GetAvailablePort 获取可用端口
func GetAvailablePort() int {
	listener := mylog.Check2(net.ListenTCP("tcp", mylog.Check2(net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", "0.0.0.0")))))
	defer func() { mylog.Check(listener.Close()) }()
	return listener.Addr().(*net.TCPAddr).Port
}

// IsPortAvailable 判断端口是否可以（未被占用）
func IsPortAvailable(port int) bool {
	address := fmt.Sprintf("%s:%d", "0.0.0.0", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		mylog.Info("port %s is taken: %s", address, err)
		return false
	}
	defer mylog.Check(listener.Close())
	return true
}
