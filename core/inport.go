package core

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type InPort struct {
	network     *Network
	pktArray    []*Packet
	is, ir      int // send index and receive index
	mtx         sync.Mutex
	condNE      *sync.Cond
	condNF      *sync.Cond
	closed      bool
	upStrmCnt   int
	portName    string
	fullName    string
	downStrProc *Process
}

func (c *InPort) Receive(p *Process) *Packet {
	return c.receive(c.downStrProc)
}

func (c *InPort) receive(p *Process) *Packet {
	LockTr(c.condNE, "recv L", p)
	defer UnlockTr(c.condNE, "recv U", p)
	trace(p, " Receiving from "+c.portName)
	for c.isEmptyNL() { // InPort is empty
		if c.closed {
			trace(p, " Received end of stream from "+c.portName)
			return nil
		}
		atomic.StoreInt32(&p.status, SuspRecv)
		WaitTr(c.condNE, "wait in recv", p)
		atomic.StoreInt32(&p.status, Active)
	}
	pkt := c.pktArray[c.ir]
	c.pktArray[c.ir] = nil
	if pkt.PktType != Normal {
		trace(p, " Received from "+c.portName+" < "+
			[...]string{"", "Open", "Close"}[pkt.PktType])
		if tracing {
			fmt.Print("  ", pkt.Contents, "\n")
		}
	} else {
		trace(p, " Received from "+c.portName+" < ")
		if tracing {
			fmt.Print("  ", pkt.Contents, "\n")
		}
	}
	c.ir = (c.ir + 1) % len(c.pktArray)
	pkt.owner = p
	p.ownedPkts++
	//c.condNF.Broadcast()
	BdcastTr(c.condNF, "bdcast recv'd", p)

	return pkt
}

func (c *InPort) incUpstream() {
	LockTr(c.condNE, "IUS L", nil)
	defer UnlockTr(c.condNE, "IUS U", nil)

	c.upStrmCnt++
}

func (c *InPort) decUpstream() {
	//LockTr(c.condNE, "DUS L", nil)
	//defer UnlockTr(c.condNE, "DUS U", nil)
	c.upStrmCnt--
}

func (c *InPort) Close() {
	LockTr(c.condNE, "ClsI L", c.downStrProc)
	defer UnlockTr(c.condNE, "ClsI U", c.downStrProc)

	c.closed = true
	//c.condNE.Broadcast()
	BdcastTr(c.condNE, "bdcast in NE", c.downStrProc)
	//c.downStrProc.activate()
}

func (c *InPort) isDrained() bool {
	LockTr(c.condNE, "IDr L", c.downStrProc)
	defer UnlockTr(c.condNE, "IDr U", c.downStrProc)

	return c.isDrainedNL()
}

func (c *InPort) isDrainedNL() bool {
	return c.isEmptyNL() && c.closed
}

func (c *InPort) isEmpty() bool {
	LockTr(c.condNE, "IE L", c.downStrProc)
	defer UnlockTr(c.condNE, "IE U", c.downStrProc)

	return c.isEmptyNL()
}

func (c *InPort) isEmptyNL() bool {
	return c.ir == c.is && c.pktArray[c.is] == nil
}

func (c *InPort) isClosed() bool {
	LockTr(c.condNE, "IC L", c.downStrProc)
	defer UnlockTr(c.condNE, "IC U", c.downStrProc)

	return c.closed
}

func (c *InPort) isFull() bool {
	return c.ir == c.is && c.pktArray[c.is] != nil
}

func (c *InPort) resetForNextExecution() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.closed = false
}

func (c *InPort) pktCount() int {
	var i int
	c.mtx.Lock()
	defer c.mtx.Unlock()
	for _, p := range c.pktArray {
		if p != nil {
			i++
		}
	}
	return i
}
