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

	markedLive   bool
	pendingSends int64
	pendingRecvs int64
}

func (c *Connection) send(p *Process, pkt *Packet) bool {
	if pkt.owner != p {
		panic("Sending packet not owned by this process")
	}
	c.mtx.Lock()
	defer c.mtx.Unlock()
	fmt.Println(p.name, "Sending", pkt.Contents)
	c.downStrProc.ensureRunning()
	if c.nolockIsFull() {
		// Any connection that has a pair of sends and receives
		// should be considered live.
		if !c.markedLive && c.pendingRecvs > 0 {
			c.markedLive = true
			c.network.incLiveConnection()
		}
		c.pendingSends++

		for c.nolockIsFull() { // connection is full
			p.transition(SuspendedSend)
			c.condNF.Wait()
			p.transition(Active)
		}

		// We unmark liveness when there are no receivers.
		// Otherwise there's a possibility that the sender terminates
		// and the receiver is still suspended.
		c.pendingSends--
		if c.markedLive && c.pendingRecvs == 0 {
			c.markedLive = false
			c.network.decLiveConnection()
		}
	}
	fmt.Println(p.name, "Sent", pkt.Contents)
	c.pktArray[c.is] = pkt
	c.is = (c.is + 1) % len(c.pktArray)
	pkt.owner = nil
	p.ownedPkts--
	c.condNE.Signal()
	return true
}

func (c *Connection) receive(p *Process) *Packet {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	fmt.Println(p.name, "Receiving")

	if c.nolockIsEmpty() {
		// Any connection that has a pair of sends and receives
		// should be considered live. It might take some time to
		// propagate those items through.
		if !c.markedLive && c.pendingSends > 0 {
			c.markedLive = true
			c.network.incLiveConnection()
		}
		c.pendingRecvs++

		for c.nolockIsEmpty() && !c.closed {
			p.transition(SuspendedRecv)
			c.condNE.Wait()
			p.transition(Active)
		}

		// Once there are no more receivers we can unmark the connection.
		c.pendingRecvs--
		if c.markedLive && c.pendingSends == 0 && c.nolockIsEmpty() {
			c.markedLive = false
			c.network.decLiveConnection()
		}
		if c.markedLive && c.pendingRecvs == 0 {
			c.markedLive = false
			c.network.decLiveConnection()
		}

		if c.closed {
			c.condNF.Broadcast()
			return nil
		}
	}

	pkt := c.pktArray[c.ir]
	c.pktArray[c.ir] = nil
	fmt.Println(p.name, "Received", pkt.Contents)
	c.ir = (c.ir + 1) % len(c.pktArray)
	pkt.owner = p
	p.ownedPkts++
	c.condNF.Signal()

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
		c.nolockClose()
	}
}

func (c *Connection) Close() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.nolockClose()
}

func (c *Connection) nolockClose() {
	// Mark connection live when there are pending receives,
	// otherwise the receivers could be marked as deadlocked
	// although they have to receive the close signal.
	if !c.markedLive && c.pendingRecvs > 0 {
		c.markedLive = true
		c.network.incLiveConnection()
	}

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
