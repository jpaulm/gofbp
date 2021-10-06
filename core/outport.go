package core

type OutPort struct {
	name     string
	Conn     *Connection
	optional bool
}

func (o *OutPort) send(p *Process, pkt *Packet) bool {
	return o.Conn.send(p, pkt)
}

func (o *OutPort) SetOptional(b bool) {
	o.optional = b
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
