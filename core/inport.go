package core

import (
	//"fmt"
	"sync"
	"sync/atomic"
)

type InPort struct {
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
	//array       []*InPort
	downStrProc *Process
}

func (c *InPort) send(p *Process, pkt *Packet) bool {
	if pkt == nil {
		panic("Sending nil packet")
	}
	if pkt.owner != p {
		panic("Sending packet not owned by this process")
	}
	c.condNF.L.Lock()
	defer c.condNF.L.Unlock()
	c.network.trace(p.name, "Sending", pkt.Contents.(string))
	c.downStrProc.ensureRunning()
	c.condNE.Broadcast()
	for c.nolockIsFull() { // InPort is full
		atomic.StoreInt32(&p.status, SuspSend)
		c.condNF.Wait()
		atomic.StoreInt32(&p.status, Active)
		atomic.StoreInt32(&p.network.Active, 1)
	}
	c.network.trace(p.name, "Sent", pkt.Contents.(string))
	c.pktArray[c.is] = pkt
	c.is = (c.is + 1) % len(c.pktArray)
	//pkt.owner = nil
	p.ownedPkts--
	pkt = nil
	return true
}

func (c *InPort) receive(p *Process) *Packet {
	c.condNE.L.Lock()
	defer c.condNE.L.Unlock()

	c.network.trace(p.name, "Receiving")
	for c.nolockIsEmpty() { // InPort is empty
		if c.closed {
			c.condNF.Broadcast()
			return nil
		}
		atomic.StoreInt32(&p.status, SuspRecv)
		c.condNE.Wait()
		atomic.StoreInt32(&p.status, Active)
		atomic.StoreInt32(&p.network.Active, 1)

	}
	pkt := c.pktArray[c.ir]
	c.pktArray[c.ir] = nil
	c.network.trace(p.name, "Received", pkt.Contents.(string))
	c.ir = (c.ir + 1) % len(c.pktArray)
	pkt.owner = p
	p.ownedPkts++
	c.condNF.Broadcast()

	return pkt
}

func (c *InPort) incUpstream() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.upStrmCnt++
}

/*
func (c *InPort) decUpstream() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.upStrmCnt--
	if c.upStrmCnt == 0 {
		c.closed = true
		c.condNE.Broadcast()
		c.downStrProc.ensureRunning()

	}
}
*/
func (c *InPort) Close() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.closed = true
	c.condNE.Broadcast()
	c.downStrProc.ensureRunning()
}

func (c *InPort) isDrained() bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.nolockIsEmpty() && c.closed
}

func (c *InPort) IsEmpty() bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.nolockIsEmpty()
}

func (c *InPort) nolockIsEmpty() bool {
	return c.ir == c.is && c.pktArray[c.is] == nil
}

func (c *InPort) IsClosed() bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.closed
}

func (c *InPort) nolockIsFull() bool {
	return c.ir == c.is && c.pktArray[c.is] != nil
}

func (c *InPort) resetForNextExecution() {}
