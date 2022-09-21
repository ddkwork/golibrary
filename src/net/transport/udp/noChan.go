package udp

import "github.com/ddkwork/golibrary/mylog"

func (o *object) TransportUDPNoChan(DstIP string, DstPort int) {
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
		if !o.SetDstAddrConn() {
			return
		}
		go o.srcWriteDstBuf()
		if !mylog.Error2(o.DstConn.Write(o.Bytes()[:o.Len()])) {
			return
		}
	}
}

func (o *object) srcWriteDstBuf() { //源连接写入目标buf
	o.Reset()
	o.BufSize, o.err = o.DstConn.Read(o.Bytes())
	if !mylog.Error(o.err) {
		return
	}
	if !mylog.Error2(o.SrcConn.WriteToUDP(o.Bytes()[:o.BufSize], o.SrcAddr)) {
		return
	}
}
