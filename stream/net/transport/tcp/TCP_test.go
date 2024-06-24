package tcp

import (
	"testing"
)

func TestInterfaceTCP(t *testing.T) {
	t.Skip()
	New().TransportTCP(``, 9999)
}
