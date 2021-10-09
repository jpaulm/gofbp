package core

import (
	"fmt"
	"sync"
)

type Connection struct {
	network     *Network
	pktArray    []*Packet
	is, ir      int // send index and receive index
	mtx         sync.Mutex
	condNE      sync.Cond
	condNF      sync.Cond
	closed      bool
	upStrmCnt   int
	portName    string
	fullName    string
	array       []*Connection
	downStrProc *Process
}

func (c *Connection) send(p *Process, pkt *Packet) bool {
	if pkt.owner != p {
		panic("Sending packet not owned by this process")
	}
	c.mtx.Lock()
	defer c.mtx.Unlock()
	fmt.Println(p.name, "Sending", pkt.Contents)
	c.downStrProc.ensureRunning()
	for c.nolockIsFull() { // connection is full
		p.transition(SuspendedSend)
		c.condNF.Wait()
		p.transition(Active)
	}
	fmt.Println(p.name, "Sent", pkt.Contents)
	c.pktArray[c.is] = pkt
	c.is = (c.is + 1) % len(c.pktArray)
	pkt.owner = nil
	p.ownedPkts--
	c.condNE.Broadcast()
	return true
}

func (c *Connection) receive(p *Process) *Packet {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	fmt.Println(p.name, "Receiving")
	for c.nolockIsEmpty() { // connection is empty
		if c.closed {
			c.condNF.Broadcast()
			return nil
		}
		p.transition(SuspendedRecv)
		c.condNE.Wait()
		p.transition(Active)
	}
	pkt := c.pktArray[c.ir]
	c.pktArray[c.ir] = nil
	fmt.Println(p.name, "Received", pkt.Contents)
	c.ir = (c.ir + 1) % len(c.pktArray)
	pkt.owner = p
	p.ownedPkts++
	c.condNF.Broadcast()

	return pkt
}

func (c *Connection) incUpstream() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.upStrmCnt++
}

func (c *Connection) decUpstream() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.upStrmCnt--
	if c.upStrmCnt == 0 {
		c.closed = true
		c.condNE.Broadcast()
		c.downStrProc.ensureRunning()

	}
}

func (c *Connection) Close() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.closed = true
	c.condNE.Broadcast()
	c.downStrProc.ensureRunning()
}

func (c *Connection) isDrained() bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.nolockIsEmpty() && c.closed
}

func (c *Connection) IsEmpty() bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.nolockIsEmpty()
}

func (c *Connection) nolockIsEmpty() bool {
	return c.ir == c.is && c.pktArray[c.is] == nil
}

func (c *Connection) IsClosed() bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.closed
}

func (c *Connection) nolockIsFull() bool {
	return c.ir == c.is && c.pktArray[c.is] != nil
}

func (c *Connection) resetForNextExecution() {}
