package core

import (
	"sync"
)

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

func (p *Process) Run(wg *sync.WaitGroup) {

	//fmt.Println(p.name)
	for {
		p.ProcFun(p)
		break
	}

	wg.Done()
}

func (p *Process) Create(s string) *Packet {
	var pt *Packet = new(Packet)
	pt.contents = s
	pt.owner = p
	p.ownedPkts++
	return pt
}
