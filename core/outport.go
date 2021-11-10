package core

import (
	"sync/atomic"
)

type OutPort struct {
	name      string
	Conn      *InPort
	connected bool
	sender    *Process
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
	trace(p.name, "Sending to "+o.name+":", pkt.Contents.(string))

	o.Conn.downStrProc.activate()
	//o.Conn.downStrProc.canGo.Broadcast()

	for o.Conn.isFull() { // InPort is full
		atomic.StoreInt32(&p.status, SuspSend)
		//o.Conn.condNF.Wait()
		WaitTr(o.Conn.condNF, "wait in send", p)
		atomic.StoreInt32(&p.status, Active)
	}

	BdcastTr(o.Conn.condNE, "bdcast out", p)

	trace(p.name, "Sent  to "+o.name)
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

//func (o *OutPort) Close() {
//	o.decUpstream()
//}

func (o *OutPort) Close() {
	LockTr(o.Conn.condNF, "close L", o.sender)
	defer UnlockTr(o.Conn.condNF, "close U", o.sender)

	o.Conn.upStrmCnt--
	if o.Conn.upStrmCnt == 0 {
		o.Conn.closed = true
		//o.Conn.condNE.Broadcast()
		BdcastTr(o.Conn.condNE, "bdcast out", o.sender)
		o.Conn.downStrProc.activate()
		//o.Conn.downStrProc.canGo.Signal()

	}
}

func (o *OutPort) GetSender() *Process {
	return o.sender
}

func (o *OutPort) SetSender(p *Process) {
	o.sender = p
}
