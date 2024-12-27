package udp

import "github.com/ddkwork/golibrary/mylog"

func (o *object) TransportUDP(DstIP string, DstPort int) {
	o.reset(DstIP, DstPort)
	o.GetSrcAddrConn()
	defer func() {
		mylog.Check(o.SrcConn == nil)
		mylog.Check(o.SrcConn.Close())
		mylog.Check(o.DstConn == nil)
		mylog.Check(o.DstConn.Close())
	}()
	for {
		SrcBufChan <- o.Bytes()[:o.BufSize]
		o.SetDstAddrConn()
		go o.readDstBuf()
		mylog.Check2(o.SrcConn.WriteToUDP(<-DstBufChan, o.SrcAddr))
	}
}

func (o *object) readDstBuf() {
	select {
	case b := <-SrcBufChan:
		mylog.Check2(o.DstConn.Write(b))
		o.Reset()
		o.BufSize = mylog.Check2(o.DstConn.Read(o.Bytes()))
		DstBufChan <- o.Bytes()[:o.BufSize]
	}
}
