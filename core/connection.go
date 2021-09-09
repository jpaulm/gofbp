package core

import (
	"fmt"
	"sync"
)

type Connection struct {
	network   *Network
	pktArray  []*Packet
	is, ir    int // send index and receive index
	mtx       sync.Mutex
	condNE    sync.Cond
	condNF    sync.Cond
	closed    bool
	upStrmCnt int
	portName  string
	fullName  string
}

func (c *Connection) send(p *Process, pkt *Packet) bool {
	if pkt.owner != p {
		panic("Sending packet not owned by this process")
	}
	c.condNF.L.Lock()
	fmt.Println(p.Name, "Sending", pkt.Contents)
	for c.IsFull() { // connection is full
		c.condNF.Wait()
	}
	fmt.Println(p.Name, "Sent", pkt.Contents)
	c.pktArray[c.is] = pkt
	c.is = (c.is + 1) % len(c.pktArray)
	pkt.owner = nil
	p.ownedPkts--
	c.condNE.Broadcast()
	c.condNF.L.Unlock()
	return true
}

func (c *Connection) receive(p *Process) *Packet {
	c.condNE.L.Lock()
	fmt.Println(p.Name, "Receiving")
	if c.isEmpty() { // connection is empty
		if c.closed {
			c.condNF.Broadcast()
			c.condNE.L.Unlock()
			return nil
		}
		c.condNE.Wait()
	}
	pkt := c.pktArray[c.ir]
	c.pktArray[c.ir] = nil
	fmt.Println(p.Name, "Received", pkt.Contents)
	c.ir = (c.ir + 1) % len(c.pktArray)
	pkt.owner = p
	p.ownedPkts++
	c.condNF.Broadcast()
	c.condNE.L.Unlock()
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
	}
}

func (c *Connection) Close() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.closed = true
}

func (c *Connection) IsEmpty() bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.isEmpty()
}

func (c *Connection) isEmpty() bool {
	return c.ir == c.is && c.pktArray[c.is] == nil
}

func (c *Connection) IsClosed() bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.closed
}

func (c *Connection) IsFull() bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.ir == c.is && c.pktArray[c.is] != nil
}
