package core

import (
	"fmt"
	"sync/atomic"
)

//import (
//	"github.com/gofbp/core"
//)

type Process struct {
	name      string
	network   *Network
	inPorts   map[string]InputConn
	outPorts  map[string]OutputConn
	logFile   string
	component Component
	ownedPkts int
	done      bool
	starting  bool
	//MustRun   bool
	status int32
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

func (p *Process) Send(o *OutPort, pkt *Packet) bool {
	return o.Conn.send(p, pkt)
}

func (p *Process) Receive(c InputConn) *Packet {
	return c.receive(p)
}

// allDrained returns whether any input port might return new data.
func (p *Process) allDrained() bool {
	for _, v := range p.inPorts {
		if !v.isDrained() {
			return false
		}
	}
	return true
}

func (p *Process) ensureRunning() {
	status := atomic.LoadInt32(&p.status)
	fmt.Println(p.GetName(), []string{"notStarted",
		"dormant",
		"suspSend",
		"suspRecv",
		"active",
		"terminated"}[status])
	if !atomic.CompareAndSwapInt32(&p.status, Notstarted, Active) {
		return
	}

	p.network.wg.Add(1)
	go func() { // Process goroutine
		defer p.network.wg.Done()
		p.Run()
	}()
}

func (p *Process) Run() {
	atomic.StoreInt32(&p.status, Dormant)
	defer atomic.StoreInt32(&p.status, Terminated)

	p.component.Setup(p)

	for {
		//if p.MustRun {
		atomic.StoreInt32(&p.status, Active)
		p.component.Execute(p) // single "activation"
		atomic.StoreInt32(&p.status, Dormant)
		//}

		if p.ownedPkts > 0 {
			panic(p.name + " deactivated without disposing of all owned packets")
		}

		p.done = p.allDrained()
		if p.done {
			break
		}

		for _, v := range p.inPorts {
			v.resetForNextExecution()
		}
	}

	for _, v := range p.outPorts {
		if v.GetType() != "NullOutPort" {
			if v.GetType() == "OutPort" {
				v.(*OutPort).Conn.decUpstream()
			} else {
				for _, w := range v.(*OutArrayPort).array {
					w.Conn.decUpstream()
				}
			}
		}
	}
}

func (p *Process) Create(s string) *Packet {
	var pkt *Packet = new(Packet)
	pkt.Contents = s
	pkt.owner = p
	p.ownedPkts++
	return pkt
}

func (p *Process) Discard(pkt *Packet) {
	p.ownedPkts--
}
