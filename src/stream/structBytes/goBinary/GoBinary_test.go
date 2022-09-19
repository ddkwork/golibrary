package goBinary

//
//type User struct {
//	Name string
//	Age  int
//	buf  []byte
//}
//
//type Out struct {
//	Age  int
//	Name string
//	buf  []byte
//}
//
//func TestInterfaceGob(t *testing.T) {
//	 p  := New()
//	u := New()
//	assert.True(t, p.GoBinaryEncode(u))
//	mylog.HexDump("GoBinaryEncode", p.GoBinaryBytes())
//	var out Out
//	assert.Equal(t, p.GoBinaryDecode(p.GoBinaryBytes(), &out), &out)
//	mylog.Struct(out)
//}
//
//func New() *User {
//	return &User{
//		Name: "xxxxx",
//		Age:  99,
//		buf: []byte{
//			0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
//			0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
//			0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99,
//		},
//	}
//}
