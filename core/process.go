package main

import (
	"fmt"
	"sync"
)

type Process struct {
	name string
	//procs   map[string]Process
	network *Network
	//inPorts  map[string]*InPort
	//outPorts map[string]*OutPort
	logFile string
}

/*
type Process interface {
	Name() string
	InPorts() map[string]*InPort
	OutPorts() map[string]*OutPort
	Run()
}
*/

func (p *Process) Run(wg *sync.WaitGroup) {

	fmt.Println(p.name)

	wg.Done()

}
