package core

import (
	//"fmt"
	"sync"
)

type InitializationConnection struct {
	network     *Network
	portName    string
	fullName    string
	closed      bool
	value       interface{}
	mtx         sync.Mutex
	downStrProc *Process
}

func (c *InitializationConnection) isDrained() bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.closed
}

func (c *InitializationConnection) isEmpty() bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return !c.closed
}

func (c *InitializationConnection) Receive(p *Process) *Packet {
	return c.receive(c.downStrProc)
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

func (c *InitializationConnection) isClosed() bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.closed
}

func (c *InitializationConnection) resetForNextExecution() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.closed = false
}
