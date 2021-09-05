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
	//driver  Process
	logFile string
}

func NewNetwork(name string) *Network {
	net := &Network{
		Name: name,
	}
	// Set up logging
	return net
}

func (n *Network) NewProc(nm string, x func(p *Process)) *Process {

	proc := &Process{
		Name:    nm,
		Network: n,
		logFile: "",
	}

	proc.ProcFun = x
	n.procList = append(n.procList, proc)

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

	var wg sync.WaitGroup
	defer wg.Wait()

	// FBP distinguishes between execution of the process as a whole and activating the code - the code may be deactivated and then
	// reactivated many times during the process "run"

	for _, proc := range n.procList {
		proc := proc
		wg.Add(1)
		go func() { // Process goroutine
			defer wg.Done()
			for !proc.done {
				proc.Run(n)
			}
		}()
	}
}
