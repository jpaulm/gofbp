package main

import (
	"sync"
)

type Process struct {
	name string
	//procs   map[string]Process
	network *Network
	//inPorts  map[string]*InPort
	//outPorts map[string]*OutPort
	logFile string
	myFun   func()
}

/*
type Process interface {
	Name() string
	InPorts() map[string]*InPort
	OutPorts() map[string]*OutPort
	Run()
}
*/

func (n *Network) newProc(name string, crun func()) *Process {

	proc := &Process{
		name:    name,
		network: n,
		myFun:   crun,
	}

	// Set up logging
	return proc
}

func (p *Process) Run(wg *sync.WaitGroup) {

	//fmt.Println(p.name)

	p.myFun()

	wg.Done()

}
