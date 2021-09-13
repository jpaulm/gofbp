package core

var _ Conn = (*InArrayPort)(nil)

type InArrayPort struct {
	network *Network

	portName string
	fullName string
	array    []*Connection
	closed   bool
}

var _ Conn = (*InArrayPort)(nil)

func (c *InArrayPort) isDrained() bool {
	return false
}

func (c *InArrayPort) IsEmpty() bool {
	return false
}

func (c *InArrayPort) receive(p *Process) *Packet {
	return nil
}

func (c *InArrayPort) IsClosed() bool {
	return c.closed
}

func (c *InArrayPort) resetForNextExecution() {}

func (c *InArrayPort) GetType() string {
	return "InArrayPort"
}

func (c *InArrayPort) GetArray() []*Connection {
	return nil
}

func (c *InArrayPort) SetArray(c2 *Connection, i int) {}
