package orderedmap_test

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream/orderedmap"
)

func TestNew(t *testing.T) {
	type Data struct {
		A int
		B int
	}
	o := orderedmap.New[string, string]()

	o.Set("c", "d")
	o.Set("a", "b")
	vv := map[string]string{
		"c": "d",
		"a": "b",
	}
	bb := mylog.Check2(json.Marshal(vv))
	log.Println(string(bb))

	for _, v := range o.List() {
		log.Println(v.Key, v.Value)
	}
	start := time.Now()
	b := mylog.Check2(json.Marshal(o))
	log.Println(string(b), time.Since(start))
	newO := orderedmap.New[string, Data]()
	b = []byte(`{"c": {"A":2333}}`)
	log.Println(json.Unmarshal(b, newO))
	b = mylog.Check2(json.Marshal(newO))
	log.Println(string(b), time.Since(start))
}

func TestOrderedMap_UnmarshalJSON(t *testing.T) {
	mylog.Call(func() {
		o := orderedmap.New[string, any]()
		assert.NoError(t, json.Unmarshal([]byte(data), o))
		mylog.Struct(o.Keys())
		for _, p := range o.List() {
			mylog.Info(p.Key, p.Value)
		}
	})
}

var data = `
{
  "Binary1": "game/system/session/info",
  "Message2": {
    "Packed2": [
      {
        "Varint1": 0,
        "Binary2": "d3048a459417e6c0b7d39c971b99a58029f2720f7b2a70c992c826ce48184069",
        "Binary3": "6593D03B-92BC-4BC6-BF54-D16BB0271AF9",
        "Binary4": "Apple-iPhone10,3"
      },
      {
        "Varint1": 1,
        "Binary2": "",
        "Binary3": "EAAIe5YPC68wBAPGE5l7JEtIE9BfPQmbcpQ92b0c9fD29vKn5ZCwHkutEjpX2PEcvyBLDqo15gNi1x0VN7U6d26QDaABEDaVzu3vuZBYKtvHH130O9Kna4742s6B8dtr1aKJUw7HuuyNWObpWYZCqBDGvypB7Js93oBISkwWrjbuYZAooY5vHLHFPyuIcLAcV8ZAX9sRFn3un18KZA0MEuDtgekU4ZBTJ1doEQmQ4DCLXwZDZD",
        "Binary4": "_"
      }
    ],
    "Binary3": "6593D03B-92BC-4BC6-BF54-D16BB0271AF9"
  },
  "Varint3": 0,
  "Binary4": "",
  "Varint5": 0
}
`
