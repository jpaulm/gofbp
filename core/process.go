package core

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
)

type Process struct {
	gid     uint64
	Name    string
	network *Network

	inPorts  map[string]inputCommon
	outPorts map[string]outputCommon

	logFile    string
	component  Component
	ownedPkts  int
	status     int32
	mtx        sync.Mutex
	canGo      *sync.Cond
	autoInput  inputCommon
	autoOutput inputCommon
}

const (
	Notstarted int32 = iota
	Active
	Dormant
	SuspSend
	SuspRecv
	Terminated
)

func (p *Process) OpenInPort(s string) InputConn {
	var in InputConn
	var b bool
	if len(p.inPorts) == 0 {
		panic(p.Name + ": No input ports specified")
	}
	in, b = p.inPorts[s].(InputConn)

	if !b {
		panic(p.Name + " " + s + " InPort not connected, or found other type")
	}
	return in
}

func (p *Process) OpenInArrayPort(s string) *InArrayPort {
	var in *InArrayPort
	var b bool
	if len(p.inPorts) == 0 {
		panic(p.Name + ": No input ports specified")
	}
	in, b = p.inPorts[s].(*InArrayPort)
	//if in == nil {
	//	panic(p.Name + ": Port Name not found (" + s + ")")
	//}
	if !b {
		panic(p.Name + " " + s + " InArrayPort not connected, or found other type")
	}
	return in
}

func (p *Process) OpenOutPort(s string) OutputConn {
	var out OutputConn
	var b bool
	if len(p.outPorts) == 0 {
		panic(p.Name + " " + s + " OutPort not connected")
	} else {
		out, b = p.outPorts[s].(*OutPort)
		if !b {
			panic(p.Name + " " + s + " OutPort not connected, or found other type")
		}
		out.(*OutPort).portName = s
		out.(*OutPort).fullName = p.Name + "." + s
		p.network.conns[out.(*OutPort).fullName] = out.(*OutPort).Conn
	}

	return out

}

func (p *Process) OpenOutPortOptional(s string) OutputConn {
	var out OutputConn
	var b bool
	if len(p.outPorts) == 0 {
		out = new(NullOutPort)
		p.outPorts[s] = out
	} else {
		out, b = p.outPorts[s].(*OutPort)
		if b {
			out.(*OutPort).portName = s
			out.(*OutPort).fullName = p.Name + "." + s
			p.network.conns[out.(*OutPort).fullName] = out.(*OutPort).Conn
		} else {
			out := new(NullOutPort)
			p.outPorts[s] = out
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
		out.portName = s
		out.fullName = p.Name + "." + s
		out.connected = false
	} else {
		out, b = p.outPorts[s].(*OutArrayPort)
		if !b {
			panic(p.Name + " " + s + " OutArrayPort not connected, or found other type")
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

	if !atomic.CompareAndSwapInt32(&p.status, Notstarted, Active) {
		if atomic.CompareAndSwapInt32(&p.status, Dormant, Active) {

			BdcastTr(p.canGo, "bdcast act", p)
		}
		return
	}

	go func() { // Process goroutine
		defer p.network.wg.Done()
		if tracing {
			fmt.Println("Starting goroutine", p.Name)
		}
		p.Run() //   <-------
	}()
}

func (p *Process) inputState() (bool, bool, bool) {

	LockTr(p.canGo, "IS L", p)
	defer UnlockTr(p.canGo, "IS U", p)

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

		atomic.StoreInt32(&p.status, Dormant)
		WaitTr(p.canGo, "wait in IS", p)

	}
}

func (p *Process) Run() {

	defer atomic.StoreInt32(&p.status, Terminated)
	defer trace(p, " terminated")
	trace(p, " started")

	if generate_gids {
		fmt.Println("Goroutine", p.Name+":", "no.", getGID())
	}

	p.component.Setup(p)

	//if p.selfStarting {
	//	autoStarting = true
	//}
	allDrained, hasData, selfStarting := p.inputState()

	canRun := selfStarting || hasData || !allDrained || p.autoInput != nil || p.isMustRun() && allDrained

	for canRun {

		// multiple activations, if necessary!

		trace(p, " activated")
		atomic.StoreInt32(&p.status, Active)

		p.component.Execute(p) // single "activation"

		atomic.StoreInt32(&p.status, Dormant)
		trace(p, " deactivated")

		if p.ownedPkts > 0 {
			panic(p.Name + " deactivated without disposing of all owned packets")
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

			_, b := v.(*InitializationConnection)
			if b {
				v.resetForNextExecution()
			}
		}
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
	pkt.PktType = pktType
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
