package core

//import (
//	"sync"
//)

type Process struct {
	Name string
	//procs   map[string]Process
	Network   *Network
	inPorts   map[string]*InPort
	outPorts  map[string]*OutPort
	logFile   string
	OpenPorts func(p *Process)
	ProcFun   func(p *Process)
	ownedPkts int
	done      bool
}

func (p *Process) OpenInPort(s string) *InPort {
	return p.inPorts[s]
}

func (p *Process) OpenOutPort(s string) *OutPort {
	return p.outPorts[s]
}

func (p *Process) Run(net *Network) {
	p.OpenPorts(p)

	for !p.done {
		p.ProcFun(p)
		p.done = true // fudge
	}

	//p.done = true
	for _, v := range p.inPorts {
		v.Conn.Close()
	}

	for _, v := range p.outPorts {
		v.Conn.mtx.Lock()
		v.Conn.UpStrmCnt--
		if v.Conn.UpStrmCnt == 0 {
			v.Conn.closed = true
		}
		v.Conn.mtx.Unlock()
	}

	if p.ownedPkts > 0 {
		panic(p.Name + "deactivated without disposing of all owned packets")
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
