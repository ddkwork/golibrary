package tcp

import "testing"

func TestInterfaceTCP(t *testing.T) {
	New().TransportTCP(``, 9999)
}
