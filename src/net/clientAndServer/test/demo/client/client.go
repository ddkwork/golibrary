package main

import (
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/src/net/clientAndServer/client"
	"net"
	"time"
)

type Login struct {
	Name     string
	Password string
}

func main() {
	if !client.Default.Connect(net.JoinHostPort("localhost", "9999")) {
		return
	}
	send("type1", "ddk", "12345678")
	time.Sleep(time.Second)
	send("type2", "xxoo", "999999999")
	select {}
}

func send(head, Name, Password string) {
	l := &Login{
		Name:     Name,
		Password: Password,
	}

	if !client.Default.SendJsonWithHead(head, l) {
		return
	}
	mylog.Json("", client.Default.Receive().String())
}
