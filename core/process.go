package core

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
)

//Process type defines a gofbp process.
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

//Current process status
const (
	Notstarted int32 = iota
	Active
	Dormant
	SuspSend
	SuspRecv
	Terminated
)

//OpenInPort function opens and returns InputConn
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

//OpenInArrayPort method opens and returns InArrayPort
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

//OpenOutPort method opens and returns OutputConn
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
		//p.network.conns[out.(*OutPort).fullName] = out.(*OutPort).Conn
	}

	return out

}

//OpenOutPortOptional function opens and returns OutputConn
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
			//p.network.conns[out.(*OutPort).fullName] = out.(*OutPort).Conn
		} else {
			out := new(NullOutPort)
			p.outPorts[s] = out
		}
	}

	return out
}

// not sure it makes sense to allow optional for array ports!

//OpenOutArrayPort method opens and returns OutArrayPort
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

//Send method emits Packet,  returning false when fails to send.
func (p *Process) Send(o OutputConn, pkt *Packet) bool {
	//o.SetSender(p)
	return o.send(p, pkt)
}

//Receive method accepts InputConn and returns Packet
// Receive receives from the connection.
// Returns nil, when there's no more data.
func (p *Process) Receive(c InputConn) *Packet {
	return c.receive(p)
}

//Close method closes InputConn
func (p *Process) Close(c InputConn) {
	c.Close()
}

func (p *Process) activate() {

	//
	// This function starts a goroutine if it is not started, and signal if it has been
	//
	///LockTr(p.canGo, "act L", p)
	//defer UnlockTr(p.canGo, "act U", p)

	st := atomic.LoadInt32(&p.status)

	if !atomic.CompareAndSwapInt32(&p.status, Notstarted, Active) {
		if atomic.CompareAndSwapInt32(&p.status, Dormant, Active) {
			BdcastTr(p.canGo, "bdcast act", p)
			trace(p, "Activating: from status "+[...]string{"Not Started", "Active", "Dormant",
				"SuspSend", "SuspRecv", "Terminated"}[st])
		}
		///UnlockTr(p.canGo, "act U", p)
		return
	} else {
		trace(p, "Activating: from status "+[...]string{"Not Started", "Active", "Dormant",
			"SuspSend", "SuspRecv", "Terminated"}[st])
	}
	// if status was NotStarted...
	///UnlockTr(p.canGo, "act U", p)
	go func() { // Process goroutine
		defer p.network.wg.Done()
		trace(p, "Starting goroutine "+strconv.FormatUint(getGID(), 10))
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
					allDrained = allDrained && w.isDrained()
					hasData = hasData || !w.isEmpty()
					selfStarting = false
				}
			} else {
				w, b := v.(*InPort)
				if b {
					allDrained = allDrained && v.isDrained()
					hasData = hasData || !w.isEmpty()
					selfStarting = false
				}
			}
		}

		if allDrained || hasData || selfStarting {
			return allDrained, hasData, selfStarting
		}

		atomic.StoreInt32(&p.status, Dormant)
		WaitTr(p.canGo, "wait in IS", p)
		//checkPending()

	}
}

//Run method initializes and Executes Process
func (p *Process) Run() {

	defer atomic.StoreInt32(&p.status, Terminated)
	defer trace(p, " terminated")
	trace(p, " started")

	if p.network.generateGids {
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

//Create method  creates and interface and returns a Packet
// create packet containing anything!
func (p *Process) Create(x interface{}) *Packet {
	var pkt *Packet = new(Packet)
	pkt.Contents = x
	pkt.owner = p
	p.ownedPkts++
	//y := fmt.Sprint(x)
	if p.network.tracepkts {
		fmt.Print(p.Name, "  Packet created < ", pkt.Contents, "\n")
		//fmt.Println("  ", pkt.Contents)
	}
	return pkt
}

//CreateBracket method builds a new Bracket and returns it
func (p *Process) CreateBracket(pktType int32, s string) *Packet {
	var pkt *Packet = new(Packet)
	pkt.Contents = s
	pkt.PktType = pktType
	pkt.owner = p
	p.ownedPkts++
	if p.network.tracepkts {
		fmt.Print(p.Name, "  Bracket created: ", [...]string{"", "Open", "Close"}[pkt.PktType]+" Bracket ",
			pkt.Contents, "\n")
	}
	return pkt
}

//CreateSignal method builds a new Signal IP and returns it
func (p *Process) CreateSignal(s string) *Packet {
	var pkt *Packet = new(Packet)
	pkt.Contents = s
	pkt.PktType = Signal
	pkt.owner = p
	p.ownedPkts++
	if p.network.tracepkts {
		fmt.Print(p.Name, "  Signal created: ", pkt.Contents, "\n")
	}
	return pkt
}

//Discard method sets Packet to nil
func (p *Process) Discard(pkt *Packet) {
	if pkt == nil {
		panic("Discarding nil packet")
	}
	p.ownedPkts--
	if p.network.tracepkts {
		//x := fmt.Sprint(pkt.Contents)
		if pkt.PktType == OpenBracket || pkt.PktType == CloseBracket {
			fmt.Print(p.Name, "  Bracket discarded: ", [...]string{"", "Open", "Close"}[pkt.PktType]+" Bracket ", pkt.Contents, "\n")

		} else {
			if pkt.PktType == Signal {
				fmt.Print(p.Name, "  Signal discarded: ", pkt.Contents, "\n")
			} else {
				fmt.Print(p.Name, "  Packet discarded > ", pkt.Contents, "\n")
			}
		}
	}
	//}
	pkt = nil
}

//DiscardOldest method sets Packet to nil
func (p *Process) discardOldest(pkt *Packet) {
	if pkt == nil {
		panic("Discarding nil packet (DO)")
	}
	p.ownedPkts--
	if p.network.tracepkts {
		//x := fmt.Sprint(pkt.Contents)
		if pkt.PktType == OpenBracket || pkt.PktType == CloseBracket {
			fmt.Print(p.Name, "  Bracket discarded (DO): ", [...]string{"", "Open", "Close"}[pkt.PktType]+" Bracket ", pkt.Contents, "\n")
		} else {
			if pkt.PktType == Signal {
				fmt.Print(p.Name, "  Signal discarded (DO): ", pkt.Contents, "\n")
			} else {
				fmt.Print(p.Name, "  Packet discarded (DO) > ", pkt.Contents, "\n")
			}
		}
	}
	//}
	pkt = nil
}

func (p *Process) NewChain(pkt *Packet, name string) (*Chain, bool) {
	if pkt == nil {
		panic("Creating chain onto nil packet")
	}
	if pkt.chains == nil {
		pkt.chains = make(map[string]*Chain)
	}
	x := pkt.chains[name]
	if x != nil {
		return x, false
	}
	x = &Chain{}
	x.name = name
	pkt.chains[name] = x
	x.owner = pkt
	return x, true
}

func (p *Process) GetChain(pkt *Packet, name string) (*Chain, bool) {
	if pkt == nil {
		panic("Getting chain from nil packet")
	}
	if pkt.chains == nil {
		panic("No chains attached to packet")
	}
	x := pkt.chains[name]
	if x == nil {
		return x, false
	}
	return x, true
}

// Attach `subpkt` to `pkt` via chain named `name`
func (p *Process) Attach(c *Chain, subpkt *Packet) {

	if subpkt == nil {
		panic("Attaching nil packet")
	}
	if subpkt.owner != p {
		panic("Attaching packet not owned by this process")
	}
	if c.First == nil {
		c.First = subpkt
	}
	if c.Last != nil {
		c.Last.Next = subpkt
	}
	c.Last = subpkt
	subpkt.owner = c
	p.ownedPkts--
}

// Detach `subpkt` from `pkt` via chain named `name`
func (p *Process) Detach(chn *Chain, subpkt *Packet) {

	if subpkt == nil {
		panic("Detaching nil packet")
	}
	if subpkt.owner != chn {
		panic("Detaching packet not owned by this chain")
	}

	pt := chn.First
	if pt == subpkt {
		chn.First = pt.Next
		subpkt.owner = p
		p.ownedPkts++
		return
	}
	pt2 := pt
	for {
		if pt == nil {
			panic("Packet being detached not found")
		}
		if pt == subpkt {
			pt2.Next = pt.Next
			if chn.Last == pt {
				chn.Last = pt2
			}
			subpkt.owner = p
			p.ownedPkts++
			break
		}
		pt2 = pt
		pt = pt.Next
	}
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
