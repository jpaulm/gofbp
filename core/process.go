package core

//import (
//	"sync"
//)

type Process struct {
	Name       string
	Network    *Network
	inPorts    map[string]Conn
	outPorts   map[string]*OutPort
	logFile    string
	component  Component
	ownedPkts  int
	done       bool
	allDrained bool
	hasData    bool
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

	for !p.done {
		p.hasData = false
		p.allDrained = true
		for _, v := range p.inPorts {
			v.Lock()
			if !v.IsEmpty() {
				p.hasData = true
			}
			if !v.IsClosed() {
				p.allDrained = false
			}
			v.Unlock()
		}

		if len(p.inPorts) == 0 || !p.allDrained {
			p.component.Execute(p) // activate component Execute logic

			if p.ownedPkts > 0 {
				panic(p.Name + "deactivated without disposing of all owned packets")
			}
		}
		if p.allDrained {
			break
		}
	}
	p.done = true
	for _, v := range p.inPorts {
		v.Lock()
		if !(v.IsClosed() && !v.IsEmpty()) {
			p.done = false
		}
		v.Unlock()
	}

	for _, v := range p.outPorts {
		v.Conn.mtx.Lock()
		v.Conn.UpStrmCnt--
		if v.Conn.UpStrmCnt == 0 {
			v.Conn.closed = true
		}
		v.Conn.mtx.Unlock()
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
