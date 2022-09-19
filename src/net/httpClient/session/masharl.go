package session

import (
	_ "embed"
	"github.com/ddkwork/golibrary/mylog"

	"google.golang.org/protobuf/proto"
)

//go:embed session.bin
var sessionBuf []byte

func GooglePb() (b []byte) {
	obj := &Session{
		Group1: &Session_Group1{
			Binary1: proto.String("game/system/session/info"),
			Packed2: &Packed2Obj{
				Group2: []*Packed2Obj_Group2{
					{
						Varint1: proto.Uint64(0),
						Binary2: proto.String("d3048a459417e6c0b7d39c971b99a58029f2720f7b2a70c992c826ce48184069"),
						Binary3: proto.String("6593D03B-92BC-4BC6-BF54-D16BB0271AF9"),
						Binary4: proto.String("Apple-iPhone10,3"),
					},
					{
						Varint1: proto.Uint64(1),
						Binary2: proto.String(""),
						Binary3: proto.String("EAAIe5YPC68wBAPGE5l7JEtIE9BfPQmbcpQ92b0c9fD29vKn5ZCwHkutEjpX2PEcvyBLDqo15gNi1x0VN7U6d26QDaABEDaVzu3vuZBYKtvHH130O9Kna4742s6B8dtr1aKJUw7HuuyNWObpWYZCqBDGvypB7Js93oBISkwWrjbuYZAooY5vHLHFPyuIcLAcV8ZAX9sRFn3un18KZA0MEuDtgekU4ZBTJ1doEQmQ4DCLXwZDZD"),
						//Binary4: myreflect.String(""), //todo
					},
				},
				Binary3: proto.String("6593D03B-92BC-4BC6-BF54-D16BB0271AF9"),
			},
			Varint3: proto.Uint64(0),
			Binary4: proto.String(""),
			Varint5: proto.Uint64(0),
		},
	}
	marshal, err := proto.Marshal(obj)
	if !mylog.Error(err) {
		return
	}
	return marshal
}
