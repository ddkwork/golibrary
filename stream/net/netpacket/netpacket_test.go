package netpacket

import (
	"testing"

	"github.com/ddkwork/golibrary/widget/table"
)

type (
	Doc struct {
		Api      string
		Function string
		Note     string
		Todo     string
		Chinese  string
	}
	doc struct{ infos []Doc }
)

func TestDoc(t *testing.T) {
	root := table.New("netPacket", Doc{
		Api:      "",
		Function: "",
		Note:     "",
		Todo:     "",
		Chinese:  "",
	})
	root.SetRowCallback(func(node *table.Node[Doc]) (cell []string) {
		if node.Container() {
			node.MetaData.Api = node.Type
		}
		return []string{
			node.MetaData.Api,
			node.MetaData.Function,
			node.MetaData.Note,
			node.MetaData.Todo,
			node.MetaData.Chinese,
		}
	})
	root.AddChildByData(Doc{
		Api:      "RegisterPacketHandles(handles []Handle)",
		Function: "Register Packet Handles as work pool",
		Note:     "",
		Todo:     "",
		Chinese:  "方是的是的是否电风扇所发生的",
	})
	root.AddChildByData(Doc{
		Api:      "SetHandles(handles []Handle)",
		Function: "RegisterPacketHandles means",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	root.AddChildByData(Doc{
		Api:      "Handles() []Handle",
		Function: "work pool",
		Note:     "work center",
		Chinese:  "",
		Todo:     "",
	})
	root.AddChildByData(Doc{
		Api:      "HandlePacket() (ok bool)",
		Function: "rang work pool packet and post them",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	root.AddChildByData(Doc{
		Api:      "SetEvent(event any)",
		Function: "set a event for any work",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	root.AddChildByData(Doc{
		Api:      "SetEventsCap(cap int)",
		Function: "Set Events Cap",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	root.AddChildByData(Doc{
		Api:      "Events() <-chan any",
		Function: "pop events",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	root.AddChildByData(Doc{
		Api:      "HandleEvent()",
		Function: "rang event and handle them",
		Note:     "need handle by your self",
		Chinese:  "",
		Todo:     "",
	})
	root.AddChildByData(Doc{
		Api:      "HttpClient() httpClient.Interface",
		Function: "http client",
		Note:     "",
		Chinese:  "",
		Todo:     "add udp wss etc",
	})
	root.AddChildByData(Doc{
		Api:      "",
		Function: "",
		Note:     "",
		Todo:     "",
		Chinese:  "",
	})
	root.AddChildByData(Doc{
		Api:      "",
		Function: "",
		Note:     "",
		Chinese:  "",
		Todo:     "add Worn echo and saveDataBase api",
	})
	println(root.Document())
}
