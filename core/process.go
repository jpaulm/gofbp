package core

type Process struct {
	Name      string
	Network   *Network
	inPorts   map[string]InputConn
	outPorts  map[string]OutputConn
	logFile   string
	component Component
	ownedPkts int
	done      bool
}

func (p *Process) OpenInPort(s string) InputConn {
	return p.inPorts[s]
}

func (p *Process) OpenInArrayPort(s string) InputConn {
	return p.inPorts[s]
}

func (p *Process) OpenOutPort(s string) *OutPort {
	return p.outPorts[s].(*OutPort)
}

func (p *Process) OpenOutArrayPort(s string) OutputConn {
	return p.outPorts[s]
}

func (p *Process) Send(c *Connection, pkt *Packet) bool {
	return c.send(p, pkt)
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

func (p *Process) Run(net *Network) {
	p.component.OpenPorts(p)

	for !p.done {
		//if p.component.GetMustRun(p) {
		p.component.Execute(p) // activate component Execute logic
		//}

		if p.ownedPkts > 0 {
			panic(p.Name + "deactivated without disposing of all owned packets")
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
		if v.GetType() == "OutPort" {
			v.(*OutPort).Conn.decUpstream()
		} else {
			for _, w := range v.(*OutArrayPort).array {
				w.Conn.decUpstream()
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
