package udp

import "testing"

func TestInterfaceUDP(t *testing.T) {
	t.Skip()
	p := New()
	p.TransportUDPNoChan(``, 6001)
	p.TransportUDP(``, 6001)
}
