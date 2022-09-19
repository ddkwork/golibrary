package portscan

type (
	packet chan []byte
)

type (
	Interface interface {
		send(b []byte)
	}
	object struct {
		data packet
	}
)

func (o *object) send(b []byte) {
	o.data <- b
}

var Default = New()

func New() Interface {
	return &object{
		data: make(packet, 4),
	}
}

func main() {
	Default.send([]byte{0x11})
}
