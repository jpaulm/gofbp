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
	//net.Wg.Add(1)
	//fmt.Println(p.name)
	defer net.Wg.Done()

	//for i := 0; i < 4; i++ {
	p.ProcFun(p)
	//}
	p.done = true
	if p.InConn != nil {
		p.InConn.closed = true
	}
	if p.OutConn != nil {
		p.OutConn.closed = true
	}
	//wg.Done()
}

func (p *Process) Create(s string) *Packet {
	var pkt *Packet = new(Packet)
	pkt.Contents = s
	pkt.owner = p
	p.ownedPkts++
	return pkt
}
