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
	inConn    *Connection
	OutConn   *Connection
	ownedPkts int
}

func (p *Process) Run(net *Network) {
	net.Wg.Add(1)
	//fmt.Println(p.name)
	for i := 0; i < 4; i++ {
		p.ProcFun(p)
	}

	net.Wg.Done()

	//wg.Done()
}

func (p *Process) Create(s string) *Packet {
	var pt *Packet = new(Packet)
	pt.contents = s
	pt.owner = p
	p.ownedPkts++
	return pt
}
