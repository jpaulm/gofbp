package main

import (
	"fmt"
	//"strconv"
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
	defer fmt.Println(n.name + "Done")
	//for i := 0; i < 4; i++ {

	proc := n.newProc("Sender", sender)
	n.wg.Add(1)
	go proc.Run(&n.wg)
	//}

	n.wg.Wait()
}
