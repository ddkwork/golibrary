package main

import "github.com/ddkwork/golibrary/src/net/transport/tcp"

func main() {
	tcp.New().TransportTCP(``, 9999)
}
