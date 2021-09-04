package core

import (
	"fmt"
	"sync"
)

/********

  Going to give up on Lists - I suspect a bug in the Golang driver

***********/

type Network struct {
	Name string
	//procList *list.List
	procList []*Process
	procNo   int
	//driver  Process
	logFile string
	Wg      *sync.WaitGroup
}

func NewNetwork(name string) *Network {
	net := &Network{
		Name: name,
		Wg:   new(sync.WaitGroup),
	}
	//net.procList = list.New()
	net.procList = make([]*Process, 10, 200) // I assume it will take up 200 slots - ugghh!
	// Set up logging
	return net
}

func (n *Network) NewProc(x func(p *Process)) *Process {

	proc := &Process{
		Network: n,
		logFile: "",
	}

	proc.ProcFun = x
	n.procList[n.procNo] = proc
	n.procNo++

	// Set up logging
	return proc
}

func (n *Network) NewConnection(cap int) *Connection {

	conn := &Connection{
		network: n,
	}

	conn.mtx = sync.Mutex{}
	conn.condNE = sync.NewCond(&conn.mtx)
	conn.condNF = sync.NewCond(&conn.mtx)
	conn.pktArray = make([]*Packet, cap, cap)
	return conn
}

func (n *Network) Run() {
	defer fmt.Println(n.Name + " Done")
	fmt.Println(n.Name + " Starting")
	for i := 0; i < n.procNo; i++ {
		n.Wg.Add(1)
		go n.procList[i].Run(n)
		//n.Wg.Done()
	}
	n.Wg.Wait()
}
