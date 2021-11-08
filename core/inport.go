package core

import (
	//"fmt"
	"sync"
	"sync/atomic"
)

type InPort struct {
	network   GenNet
	pktArray  []*Packet
	is, ir    int // send index and receive index
	mtx       sync.Mutex
	condNE    *sync.Cond
	condNF    *sync.Cond
	closed    bool
	upStrmCnt int
	//portName  string
	name string
	//array       []*InPort
	downStrProc *Process
}

func (c *InPort) receive(p *Process) *Packet {
	LockTr(c.condNE, "recv L", p)
	defer UnlockTr(c.condNE, "recv U", p)

	trace(p.name, "Receiving from "+c.name+":")
	for c.isEmpty() { // InPort is empty
		if c.closed {
			//c.condNF.Broadcast()
			BdcastTr(c.condNF, "bdcast in NF", p)
			return nil
		}
		atomic.StoreInt32(&p.status, SuspRecv)
		//c.condNE.Wait()
		WaitTr(c.condNE, "wait in recv", p)
		atomic.StoreInt32(&p.status, Active)

	}
	pkt := c.pktArray[c.ir]
	c.pktArray[c.ir] = nil
	trace(p.name, "Received from "+c.name+":", pkt.Contents.(string))
	c.ir = (c.ir + 1) % len(c.pktArray)
	pkt.owner = p
	p.ownedPkts++
	//c.condNF.Broadcast()
	BdcastTr(c.condNF, "bdcast in NF", p)

	return pkt
}

func (c *InPort) incUpstream() {
	LockTr(c.condNE, "IUS L", c.downStrProc)
	defer UnlockTr(c.condNE, "IUS U", c.downStrProc)

	c.upStrmCnt++
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

func (c *InPort) resetForNextExecution() {}
