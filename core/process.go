package core

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
)

//import (
//	"github.com/gofbp/core"
//)

type Process struct {
	name    string
	network GenNet

	inPorts  map[string]inputCommon
	outPorts map[string]outputCommon

	logFile   string
	component Component
	ownedPkts int
	//done         bool
	//selfStarting bool // process has no non-IIP input ports
	//MustRun   bool
	status     int32
	mtx        sync.Mutex
	canGo      *sync.Cond
	autoInput  inputCommon
	autoOutput inputCommon
	//allDrained bool
	//hasData    bool
	//repeat     bool
}

const (
	Notstarted int32 = iota
	Active
	Dormant
	SuspSend
	SuspRecv
	Terminated
)

func (p *Process) GetName() string {
	return p.name
}

func (p *Process) OpenInPort(s string) InputConn {
	var in InputConn
	var b bool
	if len(p.inPorts) == 0 {
		panic(p.name + ": No input ports specified")
	}
	in, b = p.inPorts[s].(InputConn)
	//if in == nil {
	//	panic(p.name + ": Port name not found (" + s + ")")
	//}
	if !b {
		panic(p.name + " " + s + " InPort not connected, or found other type")
	}
	return in
}

func (p *Process) OpenInArrayPort(s string) *InArrayPort {
	var in *InArrayPort
	var b bool
	if len(p.inPorts) == 0 {
		panic(p.name + ": No input ports specified")
	}
	in, b = p.inPorts[s].(*InArrayPort)
	//if in == nil {
	//	panic(p.name + ": Port name not found (" + s + ")")
	//}
	if !b {
		panic(p.name + " " + s + " InArrayPort not connected, or found other type")
	}
	return in
}

func (p *Process) OpenOutPort(s string) OutputConn {
	var out OutputConn
	var b bool
	if len(p.outPorts) == 0 {
		panic(p.name + " " + s + " OutPort not connected")
	} else {
		out, b = p.outPorts[s].(*OutPort)
		if !b {
			panic(p.name + " " + s + " OutPort not connected, or found other type")
		}
	}

	return out

}

func (p *Process) OpenOutPortOptional(s string) OutputConn {
	var out OutputConn
	var b bool
	if len(p.outPorts) == 0 {
		out = new(NullOutPort)
		p.outPorts[s] = out
		//out.name = s
		//out.connected = false
	} else {
		out, b = p.outPorts[s].(*OutPort)
		//if !b {
		//	panic(p.name + " " + s + " OutPort not connected, or found other type")
		//}

		if !b {
			out := new(NullOutPort)
			p.outPorts[s] = out
			//out.name = s
			//out.connected = false
		}
	}

	return out
}

// not sure it makes sense to allow optional for array ports!

func (p *Process) OpenOutArrayPort(s string) *OutArrayPort {
	var out *OutArrayPort
	var b bool
	if len(p.outPorts) == 0 {
		out = new(OutArrayPort)
		p.outPorts[s] = out
		out.name = s
		out.connected = false
	} else {
		out, b = p.outPorts[s].(*OutArrayPort)
		if !b {
			panic(p.name + " " + s + " OutArrayPort not connected, or found other type")
		}
	}

	return out
}

// Send sends a packet to the output connection.
// Returns false when fails to send.
func (p *Process) Send(o OutputConn, pkt *Packet) bool {
	//o.SetSender(p)
	return o.send(p, pkt)
}

// Receive receives from the connection.
// Returns nil, when there's no more data.
func (p *Process) Receive(c InputConn) *Packet {
	return c.receive(p)
}

func (p *Process) Close(c InputConn) {
	c.Close()
}

func (p *Process) activate() {

	//
	// This function starts a goroutine if it is not started, and signal if not
	//

	//LockTr(p.canGo, "act L", p)
	//defer UnlockTr(p.canGo, "act L", p)

	if !atomic.CompareAndSwapInt32(&p.status, Notstarted, Active) {

		//LockTr(p.canGo, "act L", p)
		BdcastTr(p.canGo, "bdcast IS", p)
		//UnlockTr(p.canGo, "act U", p)

		return
	}

	netx, b := p.network.(*Network)

	var wg *sync.WaitGroup
	if b {
		wg = &netx.wg
	} else {
		nets, _ := p.network.(*Subnet)
		wg = &nets.wg
	}

	go func() { // Process goroutine
		defer wg.Done()
		fmt.Println("Starting goroutine", p.GetName())
		p.Run() //   <-------
	}()
}

func (p *Process) inputState() (bool, bool, bool) {

	//LockTr(p.canGo, "IS L", p)
	//defer UnlockTr(p.canGo, "IS U", p)

	allDrained := true
	hasData := false
	selfStarting := true

	for {
		for _, v := range p.inPorts {
			_, b := v.(*InArrayPort)
			if b {
				for _, w := range v.(*InArrayPort).array {
					allDrained = allDrained && w.IsDrained()
					hasData = hasData || !w.IsEmpty()
					selfStarting = false
				}
			} else {
				w, b := v.(*InPort)
				if b {
					allDrained = allDrained && v.IsDrained()
					hasData = hasData || !w.IsEmpty()
					selfStarting = false
				}
			}
		}

		if allDrained || hasData || selfStarting {
			return allDrained, hasData, selfStarting
		}

		//fmt.Println("waiting for more data on canGo")
		//p.canGo.Wait()
		WaitTr(p.canGo, "wait in IS", p)
	}

}

func (p *Process) Run() {

	//var autoStarting bool

	//defer UnlockTr(p.canGo, "act L", p)
	defer atomic.StoreInt32(&p.status, Terminated)
	defer trace(p.GetName(), " terminated")
	trace(p.GetName(), " started")

	fmt.Println("Goroutine", p.GetName()+":", "no.", getGID())

	p.component.Setup(p)

	//if p.selfStarting {
	//	autoStarting = true
	//}
	allDrained, hasData, selfStarting := p.inputState()

	canRun := selfStarting || hasData || !allDrained || p.autoInput != nil || p.isMustRun() && allDrained

	for canRun {

		// multiple activations, if necessary!

		//if p.repeat {
		//	LockTr(p.canGo, "act L", p)
		//}

		trace(p.GetName(), " activated")
		atomic.StoreInt32(&p.status, Active)

		p.component.Execute(p) // single "activation"

		atomic.StoreInt32(&p.status, Dormant)
		trace(p.GetName(), " deactivated")

		if p.ownedPkts > 0 {
			panic(p.name + " deactivated without disposing of all owned packets")
		}

		if selfStarting {
			break
		}

		allDrained, _, _ := p.inputState()

		if allDrained {
			if p.autoOutput != nil {
				p.autoOutput.Close()
			}
			break
		}

		for _, v := range p.inPorts {
			v.resetForNextExecution() // resets IIPs
		}
		//}

		//UnlockTr(p.canGo, "act L", p)
		//p.repeat = true
	}

	for _, v := range p.outPorts {

		if v.IsConnected() {
			v.Close()
		}
	}

}

func (p *Process) isMustRun() bool {
	_, hasMustRun := p.component.(ComponentWithMustRun)
	return hasMustRun
}

// create packet containing anything!
func (p *Process) Create(x interface{}) *Packet {
	var pkt *Packet = new(Packet)
	pkt.Contents = x
	pkt.owner = p
	p.ownedPkts++
	return pkt
}

// create bracket
func (p *Process) CreateBracket(pktType int32, s string) *Packet {
	var pkt *Packet = new(Packet)
	pkt.Contents = s
	pkt.pktType = pktType
	pkt.owner = p
	p.ownedPkts++
	return pkt
}

func (p *Process) Discard(pkt *Packet) {
	if pkt == nil {
		panic("Discarding nil packet")
	}
	p.ownedPkts--
	pkt = nil
}

//https://blog.sgmansfield.com/2015/12/goroutine-ids/

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
