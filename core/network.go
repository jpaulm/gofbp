package core

import (
	"fmt"
	"sync"
)

type Network struct {
	Name  string
	procs map[string]Process
	//driver  Process
	logFile string
	Wg      *sync.WaitGroup
}

func NewNetwork(name string) *Network {
	net := &Network{
		Name:  name,
		procs: map[string]Process{},
		Wg:    new(sync.WaitGroup),
	}

	// Set up logging
	return net
}

func (n *Network) NewProc(x func(p *Process)) *Process {

	proc := &Process{
		Network: n,
		logFile: "",
	}

	proc.ProcFun = x

	// Set up logging
	return proc
}

func (n *Network) NewConnection() *Connection {

	conn := &Connection{
		network: n,
	}
	conn.slice = make([]Packet, 10, 10)
	return conn
}

func (n *Network) Run() {
	defer fmt.Println(n.Name + " Done")
}
