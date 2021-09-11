package core

type Process struct {
	Name         string
	Network      *Network
	inPorts      map[string]Conn
	outPorts     map[string]*OutPort
	logFile      string
	component    Component
	ownedPkts    int
	done         bool
	allDrained   bool
	hasData      bool
	selfStarting bool
}

//func (p *Process) OpenInPort(s string) *InPort {
func (p *Process) OpenInPort(s string) Conn {
	return p.inPorts[s]
}

func (p *Process) OpenOutPort(s string) *OutPort {
	return p.outPorts[s]
}

func (p *Process) Send(c *Connection, pkt *Packet) bool {
	return c.send(p, pkt)
}

func (p *Process) Receive(c Conn) *Packet {
	return c.receive(p)
}

func (p *Process) Run(net *Network) {
	p.component.OpenPorts(p)

	p.selfStarting = true
	for _, v := range p.inPorts {
		if v.GetType() == "Connection" {
			p.selfStarting = false
		}
	}

	for !p.done {
		p.hasData = false
		p.allDrained = true
		for _, v := range p.inPorts {
			if !v.IsClosed() {
				p.allDrained = false
			}
			if !v.IsEmpty() {
				p.hasData = true
			}
		}

		//if /*p.selfStarting ||*/ !p.allDrained {
		p.component.Execute(p) // activate component Execute logic

		if p.ownedPkts > 0 {
			panic(p.Name + "deactivated without disposing of all owned packets")
		}
		//}

		if p.allDrained || p.selfStarting {
			break
		}
		for _, v := range p.inPorts {
			if v.GetType() == "InitializationConnection" {
				v.ResetClosed()
			}
		}

	}

	p.done = p.hasData && p.allDrained
	for _, v := range p.inPorts {
		if !(v.IsClosed() && !v.IsEmpty()) {
			p.done = false
		}
	}

	for _, v := range p.outPorts {
		v.Conn.decUpstream()
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
