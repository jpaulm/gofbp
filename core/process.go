package core

import (
	"fmt"
	"sync/atomic"
)

type Process struct {
	name      string
	network   *Network
	inPorts   map[string]InputConn
	outPorts  map[string]OutputConn
	logFile   string
	component Component
	ownedPkts int
	status    int32
}

func (p *Process) GetName() string {
	return p.name
}

func (p *Process) OpenInPort(s string) InputConn {
	if len(p.inPorts) == 0 {
		panic(p.name + ": No input ports specified")
	}
	in := p.inPorts[s]
	if in == nil {
		panic(p.name + ": Port name not found (" + s + ")")
	}
	return in
}

func (p *Process) OpenInArrayPort(s string) InputConn {
	if len(p.inPorts) == 0 {
		panic(p.name + ": No input ports specified")
	}
	in := p.inPorts[s]
	if in == nil {
		panic(p.name + ": Port name not found (" + s + ")")
	}
	return in
}

func (p *Process) OpenOutPort(s ...string) OutputConn {
	if len(p.outPorts) == 0 {
		opt := new(NullOutPort)
		p.outPorts[s[0]] = opt
		opt.name = s[0]
	}
	out := p.outPorts[s[0]]

	if len(s) == 2 && s[1] != "opt" {
		panic(p.name + ": Invalid 2nd param (" + s[1] + ")")
	}

	return out

}

func (p *Process) OpenOutArrayPort(s ...string) OutputConn {
	if len(p.outPorts) == 0 {
		opt := new(NullOutPort)
		p.outPorts[s[0]] = opt
		opt.name = s[0]
	}
	return p.outPorts[s[0]]
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

// Create creates a new packet.
func (p *Process) Create(s string) *Packet {
	var pkt *Packet = new(Packet)
	pkt.Contents = s
	pkt.owner = p
	p.ownedPkts++
	return pkt
}

// Discard safely deletes the packet.
func (p *Process) Discard(pkt *Packet) {
	p.ownedPkts--
}

// ensureRunning starts the process if it hasn't started already.
func (p *Process) ensureRunning() {
	if !atomic.CompareAndSwapInt32(&p.status, Notstarted, Active) {
		return
	}

	go func() {
		defer p.network.wg.Done()
		p.run()
	}()
}

// run executes the process.
func (p *Process) run() {
	atomic.StoreInt32(&p.status, Dormant)
	defer atomic.StoreInt32(&p.status, Terminated)
	defer fmt.Println(p.GetName(), " terminated")
	fmt.Println(p.GetName(), " started")
	p.component.Setup(p)

	canRun := !p.allInputsClosed() || p.isSelfStarting()

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

		if p.allInputsClosed() {
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

// isMustRun checks whether component has MustRun annotation.
func isMustRun(comp Component) bool {
	_, hasMustRun := comp.(ComponentWithMustRun)
	return hasMustRun
}
