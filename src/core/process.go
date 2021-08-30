package core

import (
	"sync"
)

type Process struct {
	name string
	//procs   map[string]Process
	network *Network
	//inPorts  map[string]*InPort
	//outPorts map[string]*OutPort
	logFile   string
	myFun     func(p *Process)
	inConn    *Connection
	outConn   *Connection
	ownedPkts int
}

func (n *Network) newProc(name string, cRun func(*Process)) *Process {

	proc := &Process{
		name:    name,
		network: n,
		logFile: "",
		myFun:   cRun,
	}

	// Set up logging
	return proc
}

func (p *Process) Run(wg *sync.WaitGroup) {

	//fmt.Println(p.name)
	for {
		p.myFun(p)
		break
	}

	wg.Done()
}

func (p *Process) create(s string) *Packet {
	var pt *Packet = new(Packet)
	pt.contents = s
	pt.owner = p
	p.ownedPkts++
	return pt
}
