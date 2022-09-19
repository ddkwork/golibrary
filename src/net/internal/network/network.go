package network

type (
	Object struct {
		Tcp string
		Udp string
	}
)

var Default = New()

func New() Object {
	return Object{
		Tcp: "tcp",
		Udp: "udp",
	}
}
