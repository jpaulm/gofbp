package core

import (
	"fmt"
	"sync"
	"sync/atomic"
)

//import (
//	"github.com/gofbp/core"
//)

type Process struct {
	name    string
	network *Network
	//inPorts   map[string]InputConn
	inPorts map[string]interface{}
	//outPorts  map[string]OutputConn
	outPorts  map[string]interface{}
	logFile   string
	component Component
	ownedPkts int
	//done         bool
	selfStarting bool // process has no non-IIP input ports
	//MustRun   bool
	status int32
	mtx    sync.Mutex
}

func (p *Process) GetName() string {
	return p.name
}

func (p *Process) OpenInPort(s string) *InPort {
	var in *InPort
	var b bool
	if len(p.inPorts) == 0 {
		panic(p.name + ": No input ports specified")
	}
	in, b = p.inPorts[s].(*InPort)
	//if in == nil {
	//	panic(p.name + ": Port name not found (" + s + ")")
	//}
	if !b {
		panic(p.name + " " + s + " InPort not connected, or found other type")
	}
	return in
}

func (p *Process) OpenInitializationPort(s string) *InitializationPort {
	var in *InitializationPort
	var b bool
	if len(p.inPorts) == 0 {
		panic(p.name + ": No input ports specified")
	}
	in, b = p.inPorts[s].(*InitializationPort)
	//if in == nil {
	//	panic(p.name + ": Port name not found (" + s + ")")
	//}
	if !b {
		panic(p.name + " " + s + " InitializationPort not connected, or found other type")
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

func (p *Process) OpenOutPort(s ...string) *OutPort {
	var out *OutPort
	var b bool
	if len(p.outPorts) == 0 {
		out = new(OutPort)
		p.outPorts[s[0]] = out
		out.name = s[0]
		out.connected = false
	} else {
		out, b = p.outPorts[s[0]].(*OutPort)
		if !b {
			//panic(p.name + " " + s[0] + " OutPort wrong type")
			return nil // fix later...
		}
	}

	if len(s) == 2 {

		if s[1] != "opt" {
			panic(p.name + ": Invalid 2nd param (" + s[1] + ")")
		}

		if out == nil {
			out := new(OutPort)
			p.outPorts[s[0]] = out
			out.name = s[0]
			out.connected = false

		}
	}

	return out

}

// not sure it makes sense to allow optional for array ports!

func (p *Process) OpenOutArrayPort(s ...string) *OutArrayPort {
	var out *OutArrayPort
	var b bool
	if len(p.outPorts) == 0 {
		out = new(OutArrayPort)
		p.outPorts[s[0]] = out
		out.name = s[0]
		out.connected = false
	} else {
		out, b = p.outPorts[s[0]].(*OutArrayPort)
		if !b {
			panic(p.name + " " + s[0] + " OutArrayPort not connected, or found other type")
		}
	}

	if len(s) == 2 {

		if s[1] != "opt" {
			panic(p.name + ": Invalid 2nd param (" + s[1] + ")")
		}

		if out == nil {
			out := new(OutArrayPort)
			p.outPorts[s[0]] = out
			out.name = s[0]
			out.connected = false
		}
	}

	return out
}

// Send sends a packet to the output connection.
// Returns false when fails to send.
func (p *Process) Send(o OutputConn, pkt *Packet) bool {
	return o.send(p, pkt)
}

// Receive receives from the connection.
// Returns nil, when there's no more data.
func (p *Process) Receive(c InputConn) *Packet {
	return c.receive(p)
}

func (p *Process) ensureRunning() {

	if !atomic.CompareAndSwapInt32(&p.status, Notstarted, Active) {
		return
	}

	//p.network.wg.Add(1)
	go func() { // Process goroutine
		defer p.network.wg.Done()
		p.Run()
	}()
}

func (p *Process) inputState() (bool, bool) {
	allDrained := true
	hasData := false
	for _, v := range p.inPorts {
		//if v.GetType() == "InArrayPort" {
		_, b := v.(*InArrayPort)
		if b {
			//allClosed = true
			for _, w := range v.(*InArrayPort).array {
				if !w.isDrained() /* || !w.IsClosed() */ {
					allDrained = false
				}
				hasData = hasData || !w.IsEmpty()
			}
		} else {
			_, b := v.(*InPort)
			if b {
				if !v.(*InPort).isDrained() /*|| !v.IsClosed() */ {
					allDrained = false
				}
				hasData = hasData || !v.(*InPort).IsEmpty()
			}
		}
	}
	return allDrained, hasData
}

func (p *Process) Run() {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	atomic.StoreInt32(&p.status, Dormant)
	defer atomic.StoreInt32(&p.status, Terminated)
	defer fmt.Println(p.GetName(), " terminated")
	fmt.Println(p.GetName(), " started")
	p.component.Setup(p)

	//var allDrained bool
	//var hasData bool

	allDrained, hasData := p.inputState()

	canRun := p.selfStarting || hasData || !allDrained || p.isMustRun()

	for canRun {
		// multiple activations, if necessary!
		fmt.Println(p.GetName(), " activated")
		atomic.StoreInt32(&p.status, Active)
		p.component.Execute(p) // single "activation"
		atomic.StoreInt32(&p.status, Dormant)
		fmt.Println(p.GetName(), " deactivated")

		if p.ownedPkts > 0 {
			panic(p.name + " deactivated without disposing of all owned packets")
		}

		allDrained, _ := p.inputState()

		if allDrained {
			canRun = false
		} else {
			for _, v := range p.inPorts {
				_, b := v.(InitializationPort)
				if b {
					v.(*InitializationPort).resetForNextExecution()
				}
			}
		}

	}

	for _, v := range p.outPorts {
		var port interface{}
		var b bool

		port, b = v.(*OutPort)
		//fmt.Println(p.name, v, port, " close after run")
		if b {
			if !port.(*OutPort).IsConnected() {
				continue
			}
			port.(*OutPort).Close()
			continue
		}
		port, b = v.(*OutArrayPort)
		if b {
			port.(*OutArrayPort).Close()
		}
		// if anything else, just continue
	}
}

func (p *Process) isMustRun() bool {
	_, hasMustRun := p.component.(ComponentWithMustRun)
	return hasMustRun
}

/*
// create packet containing string
func (p *Process) Create(s string) *Packet {
	var pkt *Packet = new(Packet)
	pkt.Contents = s
	pkt.owner = p
	p.ownedPkts++
	return pkt
}
*/

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
