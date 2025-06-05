package udp

import "github.com/ddkwork/golibrary/std/mylog"

func (o *object) TransportUDPNoChan(DstIP string, DstPort int) {
	o.reset(DstIP, DstPort)
	o.GetSrcAddrConn()
	defer func() {
		mylog.Check(o.SrcConn == nil)
		mylog.Check(o.SrcConn.Close())
		mylog.Check(o.DstConn == nil)
		mylog.Check(o.DstConn.Close())
	}()
	for {
		o.SetDstAddrConn()
		go o.srcWriteDstBuf()
		mylog.Check2(o.DstConn.Write(o.Bytes()[:o.Len()]))
	}
}

func (o *object) srcWriteDstBuf() {
	o.Reset()
	o.BufSize = mylog.Check2(o.DstConn.Read(o.Bytes()))
	mylog.Check2(o.SrcConn.WriteToUDP(o.Bytes()[:o.BufSize], o.SrcAddr))
}
