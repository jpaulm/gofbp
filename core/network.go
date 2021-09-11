package core

import (
	"fmt"
	"sync"
)

type Component interface {
	OpenPorts(*Process)
	Execute(*Process)
}

type Network struct {
	Name     string
	procs    map[string]*Process
	procList []*Process
	//driver  Process
	logFile string
}

func NewNetwork(name string) *Network {
	net := &Network{
		Name: name,
	}

	net.procs = make(map[string]*Process)

	return net
}

func (n *Network) NewProc(nm string, comp Component) *Process {

	proc := &Process{
		Name:      nm,
		Network:   n,
		logFile:   "",
		component: comp,
	}

	n.procList = append(n.procList, proc)
	n.procs[nm] = proc

	proc.inPorts = make(map[string]Conn)
	proc.outPorts = make(map[string]*OutPort)

	return proc
}

func (n *Network) NewConnection(cap int) *Connection {
	conn := &Connection{
		network: n,
	}
	//conn.condNE.L = &conn.mtx
	//conn.condNF.L = &conn.mtx
	conn.condNE = sync.NewCond(&conn.mtx)
	conn.condNF = sync.NewCond(&conn.mtx)
	conn.pktArray = make([]*Packet, cap, cap)
	return conn
}

func (n *Network) NewInitializationConnection() *InitializationConnection {
	conn := &InitializationConnection{
		network: n,
	}

	return conn
}

func (n *Network) Connect(p1 *Process, out string, p2 *Process, in string, cap int) {
	var conn *Connection
	anyConn := p2.inPorts[in]
	if anyConn == nil {
		conn = n.NewConnection(cap)
		conn.portName = in
		conn.fullName = p2.Name + "." + in
		p2.inPorts[in] = conn
	} else {
		var ok bool
		conn, ok = anyConn.(*Connection)
		if !ok {
			panic(fmt.Sprintf("existing connection of type %T", anyConn))
		}
	}

	opt := p1.outPorts[out]
	if opt != nil {
		panic("Outport port already connected")
	}
	opt = new(OutPort)
	p1.outPorts[out] = opt
	opt.name = out
	opt.Conn = conn
	conn.incUpstream()
}

func (n *Network) Initialize(initValue string, p2 *Process, in string) {

	//conn := p2.inPorts[in]
	//if conn == nil {
	//ipt = new(InPort)
	//ipt.Name = in
	//p2.inPorts[in] = ipt
	conn := n.NewInitializationConnection()
	p2.inPorts[in] = conn
	conn.portName = in
	conn.fullName = p2.Name + "." + in
	conn.value = initValue
	//}

}

func (n *Network) Run() {
	defer fmt.Println(n.Name + " Done")
	fmt.Println(n.Name + " Starting")

	var wg sync.WaitGroup
	defer wg.Wait()

	// FBP distinguishes between execution of the process as a whole and activating the code - the code may be deactivated and then
	// reactivated many times during the process "run"

	for _, proc := range n.procList {
		proc := proc
		wg.Add(1)
		go func() { // Process goroutine
			defer wg.Done()
			//if len(proc.inPorts) == 0 {
			proc.Run(n)
			//}
		}()
	}
}
