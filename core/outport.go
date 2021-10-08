package core

type OutPort struct {
	name     string
	conn     *Connection
	optional bool
}

func (o *OutPort) send(p *Process, pkt *Packet) bool {
	return o.conn.send(p, pkt)
}

func (o *OutPort) SetOptional(b bool) {
	o.optional = b
}

func (o *OutPort) Close() {
	o.conn.decUpstream()
}
