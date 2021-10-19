package core

type NullOutPort struct {
	//name string
}

// NullOutPort by default discards the packet.
func (*NullOutPort) send(p *Process, pkt *Packet) bool {
	p.Discard(pkt)
	return true
}

func (*NullOutPort) IsConnected() bool {
	return false
}

func (*NullOutPort) Close() {}

func (*NullOutPort) IsClosed() bool {
	return true
}
