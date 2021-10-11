package core

type OutPort struct {
	name      string
	Conn      *InPort
	connected bool
}

func (o *OutPort) send(p *Process, pkt *Packet) bool {
	if o == nil {
		return false
	}
	return o.Conn.send(p, pkt)
}

func (o *OutPort) IsConnected() bool {
	if o == nil {
		return false
	}
	return o.connected
}

func (o *OutPort) GetArrayItem(i int) *OutPort {
	return nil
}

func (o *OutPort) SetArrayItem(op *OutPort, i int) {}

func (o *OutPort) ArrayLength() int {
	return 0
}

func (o *OutPort) Close() {
	o.Conn.decUpstream()
}
