package core

import (
	"sync/atomic"
)

type OutPort struct {
	portName string
	//fullName  string
	Conn      *InPort
	connected bool
	sender    *Process
	network   *Network
}

func (o *OutPort) send(p *Process, pkt *Packet) bool {
	if o == nil {
		return false
	}

	if pkt == nil {
		panic("Sending nil packet")
	}
	if pkt.owner != p {
		panic("Sending packet not owned by this process")
	}

	LockTr(o.Conn.condNF, "send L", p)
	defer UnlockTr(o.Conn.condNF, "send U", p)

	if pkt.PktType != Normal {
		trace(p, " Sending to "+o.portName+" >", pkt.Contents.(string),
			[...]string{"", "Open", "Close"}[pkt.PktType])
	} else {
		trace(p, " Sending to "+o.portName+" >", pkt.Contents.(string))
	}

	o.Conn.downStrProc.activate()
	//o.Conn.downStrProc.canGo.Broadcast()

	for o.Conn.isFull() { // InPort is full
		atomic.StoreInt32(&p.status, SuspSend)
		//o.Conn.condNF.Wait()
		WaitTr(o.Conn.condNF, "wait in send", p)
		atomic.StoreInt32(&p.status, Active)
	}

	BdcastTr(o.Conn.condNE, "bdcast out", p)

	trace(p, " Sent to "+o.portName)

	o.Conn.pktArray[o.Conn.is] = pkt
	o.Conn.is = (o.Conn.is + 1) % len(o.Conn.pktArray)
	//pkt.owner = nil
	p.ownedPkts--
	pkt = nil
	return true
}

func (o *OutPort) IsConnected() bool {
	if o == nil {
		return false
	}
	return o.connected
}

func (o *OutPort) GetArrayItem(i int) *OutPort {
	return nil
}

func (o *OutPort) SetArrayItem(op *OutPort, i int) {}

func (o *OutPort) ArrayLength() int {
	return 0
}

func (o *OutPort) Close() {
	LockTr(o.Conn.condNF, "close L", o.sender)
	defer UnlockTr(o.Conn.condNF, "close U", o.sender)
	trace(o.sender, " Close "+o.portName)
	//o.Conn.upStrmCnt--
	o.Conn.decUpstream()
	if o.Conn.upStrmCnt == 0 {
		o.Conn.closed = true
		//o.Conn.condNE.Broadcast()
		BdcastTr(o.Conn.condNE, "bdcast out", o.sender)
		o.Conn.downStrProc.activate()
		//o.Conn.downStrProc.canGo.Signal()

	}
}
