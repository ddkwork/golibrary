package netpacket

import (
	"github.com/ddkwork/golibrary/src/mydoc"
	"testing"
)

func TestDoc(t *testing.T) {
	d := mydoc.New()
	d.Append(mydoc.Row{
		Api:      "RegisterPaApicketHandles(handles []Handle)",
		Function: "Register Packet Handles as work pool",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	d.Append(mydoc.Row{
		Api:      "SetHandles(handles []Handle)",
		Function: "RegisterPacketHandles means",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	d.Append(mydoc.Row{
		Api:      "Handles() []Handle",
		Function: "work pool",
		Note:     "work center",
		Chinese:  "",
		Todo:     "",
	})
	d.Append(mydoc.Row{
		Api:      "HandlePacket() (ok bool)",
		Function: "rang work pool packet and post them",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	d.Append(mydoc.Row{
		Api:      "SetEvent(event any)",
		Function: "set a event for any work",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	d.Append(mydoc.Row{
		Api:      "SetEventsCap(cap int)",
		Function: "Set Events Cap",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	d.Append(mydoc.Row{
		Api:      "Events() <-chan any",
		Function: "pop events",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	d.Append(mydoc.Row{
		Api:      "HandleEvent()",
		Function: "rang event and handle them",
		Note:     "need handle by your self",
		Chinese:  "",
		Todo:     "",
	})
	d.Append(mydoc.Row{
		Api:      "HttpClient() httpClient.Interface",
		Function: "http client",
		Note:     "",
		Chinese:  "",
		Todo:     "add udp wss etc",
	})
	d.Append(mydoc.Row{
		Api:      "",
		Function: "",
		Note:     "",
		Chinese:  "",
		Todo:     "add Worn echo and saveDataBase api",
	})
	body := d.Gen()
	println(body)
}
