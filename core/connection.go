package core

import (
	"fmt"
	"sync"
)

// Based on https://stackoverflow.com/questions/36857167/how-to-correctly-use-sync-cond

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
	array     []*Connection
}

func (c *Connection) send(p *Process, pkt *Packet) bool {
	if pkt.owner != p {
		panic("Sending packet not owned by this process")
	}
	c.condNF.L.Lock()
	fmt.Println(p.Name, "Sending", pkt.Contents)
	for c.isFull() { // connection is full
		p.status = suspSend
		c.condNF.Wait()
		p.status = active
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
		p.status = suspRecv
		c.condNE.Wait()
		p.status = active
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

func (c *Connection) isDrained() bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.isEmpty() && c.closed
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

func (c *Connection) isFull() bool {
	return c.ir == c.is && c.pktArray[c.is] != nil
}

func (c *Connection) resetForNextExecution() {}

func (c *Connection) GetType() string {
	return "Connection"
}

func (c *Connection) GetArrayItem(i int) *Connection {
	return nil
}

func (c *Connection) SetArrayItem(c2 *Connection, i int) {}

func (c *Connection) ArrayLength() int {
	return 0
}
