package core

import (
	"fmt"
	//"github.com/jpaulm/gofbp/components"
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

	//var wg sync.WaitGroup
	//net.wg = wg

	// Set up logging
	return net
}

func (n *Network) NewProc(name string, s Component) *Process {

	proc := &Process{
		Name:    name,
		Network: n,
		logFile: "",
		ProcFun: s.Execute,
	}

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

	//var sendFun func(*Process) = components.Sender.Execute
	//proc := n.newProc("Sender", sendFun)
	//proc.OutConn = n.newConnection()

	//n.wg.Add(1)
	//go proc.Run(&n.wg)

	//n.wg.Wait()
}
