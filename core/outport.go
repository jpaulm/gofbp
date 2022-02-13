package core

import (
	"fmt"
	"sync/atomic"
)

type OutPort struct {
	portName  string
	fullName  string
	conn      *InPort
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

	LockTr(o.conn.condNE, "send L", p)
	defer UnlockTr(o.conn.condNE, "send U", p)

	x := fmt.Sprint(pkt.Contents)
	if pkt.PktType == OpenBracket || pkt.PktType == CloseBracket {
		trace(p, " Sending to "+o.portName+" > "+
			[...]string{"", "Open", "Close"}[pkt.PktType]+" Bracket "+x)
		//if p.network.tracing {
		//	fmt.Print("  contents: ", pkt.Contents, "\n")
		//}
	} else {
		if pkt.PktType == Signal {
			//x, _ := pkt.Contents.(string)
			trace(p, " Sending to "+o.portName+" > "+
				"Signal: "+x)
		} else {
			trace(p, " Sending to "+o.portName+" > "+x)

			//if p.network.tracing {
			//	fmt.Print("  ", pkt.Contents, "\n")
			//}
		}
	}

	if o.conn.dropOldest {
		if o.conn.isFull() { // if connection is full
			old_pkt := o.conn.pktArray[o.conn.ir]
			p.DiscardOldest(old_pkt)
			o.conn.pktArray[o.conn.ir] = pkt
			o.conn.ir = (o.conn.ir + 1) % len(o.conn.pktArray)
			o.conn.is = o.conn.ir
		} else {
			o.conn.pktArray[o.conn.is] = pkt
			o.conn.is = (o.conn.is + 1) % len(o.conn.pktArray)
		}
	} else {
		for o.conn.isFull() { // while connection is full
			atomic.StoreInt32(&p.status, SuspSend)
			WaitTr(o.conn.condNF, "wait in send", p)
			atomic.StoreInt32(&p.status, Active)
			trace(p, " Sent to "+o.portName)
		}
		o.conn.pktArray[o.conn.is] = pkt
		o.conn.is = (o.conn.is + 1) % len(o.conn.pktArray)
	}
	pkt.owner = nil
	p.ownedPkts--
	//pkt = nil
	BdcastTr(o.conn.condNE, "bdcast sent", p)

	//trace(o.conn.downStrProc, "activated from send")
	o.conn.downStrProc.activate()

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
	LockTr(o.conn.condNF, "close L", o.sender)
	defer UnlockTr(o.conn.condNF, "close U", o.sender)
	trace(o.sender, " Close "+o.portName)

	o.conn.decUpstream()
	if o.conn.upStrmCnt == 0 {
		o.conn.closed = true
		BdcastTr(o.conn.condNE, "bdcast out", o.sender)
		trace(o.conn.downStrProc, "activated from close")
		o.conn.downStrProc.activate()
	}
}
