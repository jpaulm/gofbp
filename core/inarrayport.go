package core

//var _ Conn = (*InArrayPort)(nil)

type InArrayPort struct {
	network *Network

	portName string
	fullName string
	array    []*Connection
	closed   bool
}

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

func (c *InArrayPort) GetArrayItem(i int) *Connection {
	return c.array[i]
}

func (c *InArrayPort) SetArrayItem(c2 *Connection, i int) {
	c.array[i] = c2
}

func (c *InArrayPort) ArrayLength() int {
	return len(c.array)
}
