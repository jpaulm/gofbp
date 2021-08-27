package main

import (
	"fmt"
	"strconv"
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

func (n *Network) newProc(name string) *Process {

	proc := &Process{
		name:    name,
		network: n,
	}

	// Set up logging
	return proc
}

func (n *Network) Run() {
	defer fmt.Println(n.name)
	for i := 0; i < 4; i++ {

		proc := n.newProc("Sender" + strconv.Itoa(i))
		n.wg.Add(1)
		go proc.Run(&n.wg)
	}

	n.wg.Wait()
}
