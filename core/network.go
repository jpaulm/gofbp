package core

import (
	"container/list"
	"fmt"
	"reflect"
	"sync"
)

type Network struct {
	Name  string
	procs *list.List
	//driver  Process
	logFile string
	Wg      *sync.WaitGroup
}

func NewNetwork(name string) *Network {
	net := &Network{
		Name: name,
		Wg:   new(sync.WaitGroup),
	}
	net.procs = list.New()
	// Set up logging
	return net
}

func (n *Network) NewProc(x func(p *Process)) *Process {

	proc := &Process{
		Network: n,
		logFile: "",
	}

	proc.ProcFun = x
	n.procs.PushFront(*proc)

	// Set up logging
	return proc
}

func (n *Network) NewConnection() *Connection {

	conn := &Connection{
		network: n,
	}
	conn.pktArray = make([]Packet, 10, 10)
	return conn
}

func (n *Network) Run() {
	defer fmt.Println(n.Name + " Done")
	fmt.Println(n.Name + " Starting")
	for e := n.procs.Front(); e != nil; e = e.Next() {
		var v = e.Value
		fmt.Printf("%T %T\n", e, v)

		var w *Process = reflect.Addr(v.reflect.Value())
		go w.Run(n)
	}

	n.Wg.Wait()
}
