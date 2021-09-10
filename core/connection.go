package core

import (
	"fmt"
	"sync"
)

type Connection struct {
	network   *Network
	packets   chan *Packet
	mtx       sync.Mutex
	closed    bool
	upStrmCnt int
	portName  string
	fullName  string
}

func (c *Connection) send(p *Process, pkt *Packet) bool {
	if pkt.owner != p {
		panic("Sending packet not owned by this process")
	}
	fmt.Println(p.Name, "Sending", pkt.Contents)
	pkt.owner = nil
	c.packets <- pkt
	fmt.Println(p.Name, "Sent", pkt.Contents)
	p.ownedPkts--
	return true
}

func (c *Connection) receive(p *Process) *Packet {
	fmt.Println(p.Name, "Receiving")
	pkt, ok := <-c.packets
	if !ok {
		return nil
	}
	fmt.Println(p.Name, "Received", pkt.Contents)
	pkt.owner = p
	p.ownedPkts++
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
		if !c.closed {
			c.closed = true
			close(c.packets)
		}
	}
}

func (c *Connection) Close() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if !c.closed {
		c.closed = true
		close(c.packets)
	}
}

func (c *Connection) IsEmpty() bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.isEmpty()
}

func (c *Connection) isEmpty() bool {
	return len(c.packets) == 0
}

func (c *Connection) IsClosed() bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.closed
}

func (c *Connection) IsFull() bool {
	return len(c.packets) == cap(c.packets)
}

func (c *Connection) ResetClosed() {}

func (c *Connection) GetType() string {
	return "Connection"
}
