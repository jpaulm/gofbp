package core

type NullOutPort struct {
	sender *Process
}

// NullOutPort by default discards the packet.
func (o *NullOutPort) send(p *Process, pkt *Packet) bool {
	p.Discard(pkt)
	return true
}

func (o *NullOutPort) IsConnected() bool {
	return false
}

func (o *NullOutPort) Close() {}

func (o *NullOutPort) GetSender() *Process {
	return o.sender
}

func (o *NullOutPort) SetSender(p *Process) {
	o.sender = p
}
