package core

import "fmt"

type ConnArray struct {
	network  *Network
	portName string
	fullName string
	closed   bool
	value    []*Connection
}

func (c *ConnArray) isDrained() bool {
	// initialization connection can be considered always as drained,
	// since it won't produce new values.
	return true
}

func (c *ConnArray) IsEmpty() bool {
	return !c.closed
}

func (c *ConnArray) receive(p *Process) *Packet {

	if c.closed {
		return nil
	}
	fmt.Println(p.Name, "Receiving IIP")
	var pkt *Packet = new(Packet)
	pkt.Contents = c.value
	pkt.owner = p
	p.ownedPkts++
	c.closed = true
	fmt.Println(p.Name, "Received IIP: ", pkt.Contents)
	return pkt
}

func (c *ConnArray) IsClosed() bool {
	return c.closed
}

func (c *ConnArray) resetForNextExecution() {}

func (c *ConnArray) GetType() string {
	return "InitializationConnection"
}

func (c *ConnArray) GetArray() []*Connection {
	return nil
}

func (c *ConnArray) SetArray(c2 *Connection, i int) {}
