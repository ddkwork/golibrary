package main

import (
	"github.com/ddkwork/golibrary/stream/net/transport/tcp"
)

func main() {
	tcp.New().TransportTCP(``, 9999)
}
