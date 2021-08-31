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
	}

	var wg sync.WaitGroup
	net.wg = wg

	// Set up logging
	return net
}

func (n *Network) Run() {
	defer fmt.Println(n.name + " Done")

	var sendFun func(*Process) = Execute
	proc := n.newProc("Sender", sendFun)
	proc.OutConn = n.newConnection()

	n.wg.Add(1)
	go proc.Run(&n.wg)

	n.wg.Wait()
}
