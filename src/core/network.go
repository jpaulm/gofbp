package core

import (
	"fmt"
	"sync"
)

type Network struct {
	name  string
	procs map[string]Process
	//driver  Process
	logFile string
	wg      sync.WaitGroup
}

func NewNetwork(name string) *Network {
	net := &Network{
		name:  name,
		procs: map[string]Process{},
		//wg: new(sync.WaitGroup),
	}

	var wg sync.WaitGroup
	net.wg = wg

	// Set up logging
	return net
}

func (n *Network) Run() {
	defer fmt.Println(n.name + " Done")

	proc := n.newProc("Sender", sender)
	proc.outConn = n.newConnection()

	n.wg.Add(1)
	go proc.Run(&n.wg)

	n.wg.Wait()
}
