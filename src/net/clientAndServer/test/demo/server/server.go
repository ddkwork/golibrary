package main

import (
	"encoding/json"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/src/net/clientAndServer/server"
	"net"
)

type Login struct {
	head     string
	Name     string
	Password string
}

func main() {
	s := server.New()
	if !s.ListenAndServer(net.JoinHostPort("localhost", "9999")) {
		return
	}
	go func() {
		for {
			receive := s.Receive()
			packetHeadLen := len("type1")
			head := receive.Bytes()[:packetHeadLen]
			body := receive.Bytes()[packetHeadLen:]
			mylog.Json(string(head), string(body))
			l := new(Login)
			if err := json.Unmarshal(body, l); err == nil {
				mylog.Struct(l)
				s.Replay("server replay: i am receive your message")
			}
		}
	}()
	select {}
}
