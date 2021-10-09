package core

import (
	"fmt"
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

	atomicStatus int32
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

func (p *Process) OpenOutPort(name string) OutputConn {
	out, ok := p.outPorts[name]
	if !ok {
		panic(p.name + ": Port name not found (" + name + ")")
	}
	return out.(OutputConn)
}

func (p *Process) OpenOutPortOptional(name string) OutputConn {
	out, ok := p.outPorts[name]
	if !ok {
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
	if !p.transitionFrom(NotStarted, Dormant) {
		return
	}

	go func() {
		defer p.network.wg.Done()
		defer p.transition(Terminated)
		p.run()
	}()
}

func (p *Process) run() {
	p.component.Setup(p)

	runOnce := p.isSelfStarting()
	for runOnce || !p.allInputsClosed() {
		runOnce = false

		for _, v := range p.inPorts {
			v.resetForNextExecution()
		}

		// multiple activations, if necessary!
		p.transition(Active)
		p.component.Execute(p) // single "activation"
		p.transition(Dormant)

		if p.ownedPkts > 0 {
			panic(p.name + " deactivated without disposing of all owned packets")
		}
	}

	for _, v := range p.outPorts {
		v.Close()
	}
}

// isSelfStarting returns whether the process should start at the beginning of the network.
func (p *Process) isSelfStarting() bool {
	// start anything that has a MustRun annotation
	if isMustRun(p.component) {
		return true
	}

	// start anything that doesn't have any input ports
	if len(p.inPorts) == 0 {
		return true
	}

	// start anything that has an initialization connection
	for _, in := range p.inPorts {
		if _, ok := in.(*InitializationConnection); ok {
			return true
		}
	}

	return false
}

// allInputsClosed returns whether there are any inbound connections
// that might return data.
func (p *Process) allInputsClosed() bool {
	for _, v := range p.inPorts {
		if !v.isDrained() {
			return false
		}
	}
	return true
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

// Discard safely deletes the packet.
func (p *Process) Discard(pkt *Packet) {
	p.ownedPkts--
}

// isMustRun checks whether component has MustRun annotation.
func isMustRun(comp Component) bool {
	_, hasMustRun := comp.(ComponentWithMustRun)
	return hasMustRun
}

func (p *Process) status() ProcessStatus {
	return ProcessStatus(atomic.LoadInt32(&p.atomicStatus))
}

func (p *Process) transition(next ProcessStatus) {
	previous := ProcessStatus(atomic.SwapInt32(&p.atomicStatus, int32(next)))
	fmt.Printf("%s %s -> %s\n", p.GetName(), previous, next)
	p.network.processTransitioned(previous, next)
}

func (p *Process) transitionFrom(current, next ProcessStatus) bool {
	ok := atomic.CompareAndSwapInt32(&p.atomicStatus, int32(current), int32(next))
	if ok {
		fmt.Printf("%s %s -> %s\n", p.GetName(), current, next)
		p.network.processTransitioned(current, next)
	}
	return ok
}
