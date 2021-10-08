package core

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Process struct {
	name      string
	network   *Network
	inPorts   map[string]inputCommon
	outPorts  map[string]outputCommon
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

func (p *Process) OpenInPort(name string) InputConn {
	in, ok := p.inPorts[name]
	if !ok {
		panic(p.name + ": Port name not found (" + name + ")")
	}
	return in.(InputConn)
}

func (p *Process) OpenInArrayPort(name string) InputArrayConn {
	in, ok := p.inPorts[name]
	if !ok {
		panic(p.name + ": Port name not found (" + name + ")")
	}
	return in.(InputArrayConn)
}

func (p *Process) OpenOutPort(name string, opts ...string) OutputConn {
	out, ok := p.outPorts[name]
	if !ok {
		if len(opts) == 0 || opts[0] != "opt" {
			panic(p.name + ": Port name not found (" + name + ")")
		}

		out = &NullOutPort{name: name}
		p.outPorts[name] = out
	}
	return out.(OutputConn)
}

// not sure it maes sense to allow optional for array ports!

func (p *Process) OpenOutArrayPort(name string) OutputArrayConn {
	out, ok := p.outPorts[name]
	if !ok {
		panic(p.name + ": Port name not found (" + name + ")")
	}
	return out.(OutputArrayConn)
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
			_, b := v.(*Connection)
			if b {
				if !v.(*Connection).isDrained() /*|| !v.IsClosed() */ {
					allDrained = false
				}
				hasData = hasData || !v.(*Connection).IsEmpty()
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
				v.resetForNextExecution()
			}
		}

	}

	for _, v := range p.outPorts {
		v.Close()
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
	p.ownedPkts--
}
