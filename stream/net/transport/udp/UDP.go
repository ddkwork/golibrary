package udp

import "github.com/ddkwork/golibrary/mylog"

func (o *object) TransportUDP(DstIP string, DstPort int) {
	o.reset(DstIP, DstPort)
	if !o.GetSrcAddrConn() {
		return
	}
	defer func() {
		if o.SrcConn == nil {
			mylog.Error("SrcConn == nil ")
			return
		}
		mylog.Error(o.SrcConn.Close())

		if o.DstConn == nil {
			mylog.Error("DstConn == nil ")
			return
		}
		mylog.Error(o.DstConn.Close())
	}()
	for {
		SrcBufChan <- o.Bytes()[:o.BufSize]
		if !o.SetDstAddrConn() {
			return
		}
		go o.readDstBuf()
		if !mylog.Error2(o.SrcConn.WriteToUDP(<-DstBufChan, o.SrcAddr)) { // 这句提到协程内即可不用信道
			return
		}
	}
}

func (o *object) readDstBuf() { // 读目标buf
	select {
	case b := <-SrcBufChan:
		if !mylog.Error2(o.DstConn.Write(b)) {
			return
		}
		o.Reset()
		o.BufSize, o.err = o.DstConn.Read(o.Bytes())
		if !mylog.Error(o.err) {
			return
		}
		DstBufChan <- o.Bytes()[:o.BufSize]
	}
}
