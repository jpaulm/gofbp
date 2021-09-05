package core

//import (
//	"sync"
//)

type Process struct {
	Name string
	//procs   map[string]Process
	Network *Network
	//inPorts  map[string]*InPort
	//outPorts map[string]*OutPort
	logFile   string
	ProcFun   func(p *Process)
	InConn    *Connection
	OutConn   *Connection
	ownedPkts int
	done      bool
}

func (p *Process) Run(net *Network) {

	p.ProcFun(p)

	p.done = true
	if p.InConn != nil {
		p.InConn.closed = true
	}
	if p.OutConn != nil {
		p.OutConn.mtx.Lock()
		p.OutConn.UpStrmCnt--
		if p.OutConn.UpStrmCnt == 0 {
			p.OutConn.closed = true
		}
		p.OutConn.mtx.Unlock()
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
