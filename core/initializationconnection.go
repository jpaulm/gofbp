package core

import (
	"fmt"
	"sync"
)

type InitializationConnection struct {
	network  *Network
	portName string
	fullName string
	closed   bool
	value    string
	mtx      sync.Mutex
}

func (c *InitializationConnection) isDrained() bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.closed
}

func (c *InitializationConnection) IsEmpty() bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return !c.closed
}

func (c *InitializationConnection) receive(p *Process) *Packet {

	if c.closed {
		return nil
	}
	fmt.Println(p.name, "Receiving IIP")
	var pkt *Packet = new(Packet)
	pkt.Contents = c.value
	pkt.owner = p
	p.ownedPkts++
	c.Close()
	fmt.Println(p.name, "Received IIP: ", pkt.Contents)
	return pkt
}

func (c *InitializationConnection) Close() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.closed = true

}

func (c *InitializationConnection) IsClosed() bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.closed
}

func (c *InitializationConnection) resetForNextExecution() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.closed = false
}

func (c *InitializationConnection) GetType() string {
	return "InitializationConnection"
}

func (c *InitializationConnection) GetArrayItem(i int) *Connection {
	return nil
}

func (c *InitializationConnection) SetArrayItem(c2 *Connection, i int) {}

func (c *InitializationConnection) ArrayLength() int {
	return 0
}
