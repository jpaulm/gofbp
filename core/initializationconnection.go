package core

import (
	//"fmt"
	"sync"
)

type InitializationConnection struct {
	network  *Network
	portName string
	fullName string
	closed   bool
	value    interface{}
	mtx      sync.Mutex
}

func (c *InitializationConnection) IsDrained() bool {
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
	trace(p, " Receiving IIP")
	var pkt *Packet = new(Packet)
	pkt.Contents = c.value
	pkt.owner = p
	p.ownedPkts++
	c.Close()
	trace(p, " Received IIP: ", pkt.Contents.(string))
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
