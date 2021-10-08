package core

type OutPort struct {
	name string
	conn *Connection
}

func (o *OutPort) send(p *Process, pkt *Packet) bool {
	return o.conn.send(p, pkt)
}

func (o *OutPort) IsConnected() bool { return true }

func (o *OutPort) Close() {
	o.conn.decUpstream()
}
