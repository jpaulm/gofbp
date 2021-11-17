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
	condNE    *sync.Cond
	condNF    *sync.Cond
	closed    bool
	upStrmCnt int
	portName  string
	//fullName  string
	//array       []*InPort
	downStrProc *Process
}

func (c *InPort) receive(p *Process) *Packet {
	LockTr(c.condNE, "recv L", p)
	defer UnlockTr(c.condNE, "recv U", p)
	if c.isDrained() {
		return nil
	}
	trace(p, " Receiving from "+c.portName)
	for c.isEmpty() { // InPort is empty
		if c.closed {
			//c.condNF.Broadcast()
			BdcastTr(c.condNF, "bdcast in NF cl", p)
			trace(p, " Received end of stream from "+c.portName)
			return nil
		}
		atomic.StoreInt32(&p.status, SuspRecv)
		//c.condNE.Wait()
		WaitTr(c.condNE, "wait in recv", p)
		atomic.StoreInt32(&p.status, Active)
	}
	pkt := c.pktArray[c.ir]
	c.pktArray[c.ir] = nil
	if pkt.PktType != Normal {
		trace(p, " Received from "+c.portName+" <", pkt.Contents.(string),
			[...]string{"", "Open", "Close"}[pkt.PktType])
	} else {
		trace(p, " Received from "+c.portName+" <", pkt.Contents.(string))
	}
	c.ir = (c.ir + 1) % len(c.pktArray)
	pkt.owner = p
	p.ownedPkts++
	//c.condNF.Broadcast()
	BdcastTr(c.condNF, "bdcast in NF", p)

	return pkt
}

func (c *InPort) incUpstream() {
	LockTr(c.condNE, "IUS L", nil)
	defer UnlockTr(c.condNE, "IUS U", nil)

	c.upStrmCnt++
}

func (c *InPort) decUpstream() {
	//LockTr(c.condNE, "DUS L", nil) // sender is one of senders
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

func (c *InPort) IsDrained() bool {
	LockTr(c.condNE, "IDr L", c.downStrProc)
	defer UnlockTr(c.condNE, "IDr U", c.downStrProc)

	return c.isDrained()
}

func (c *InPort) isDrained() bool {
	return c.isEmpty() && c.closed
}

func (c *InPort) IsEmpty() bool {
	LockTr(c.condNE, "IE L", c.downStrProc)
	defer UnlockTr(c.condNE, "IE U", c.downStrProc)

	return c.isEmpty()
}

func (c *InPort) isEmpty() bool {
	return c.ir == c.is && c.pktArray[c.is] == nil
}

func (c *InPort) IsClosed() bool {
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
