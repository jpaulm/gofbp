package core

import (
	"container/list"
	"fmt"
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

func printSlice(s *list.List) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

func (n *Network) NewProc(x func(p *Process)) *Process {

	proc := &Process{
		Network: n,
		logFile: "",
	}

	proc.ProcFun = x
	n.procs.PushFront(*proc)

	printSlice(n.procs)

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
		go e.Run(n)
	}

	n.Wg.Wait()
}
