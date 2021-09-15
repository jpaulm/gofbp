package core

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
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
	conn.condNE.L = &conn.mtx
	conn.condNF.L = &conn.mtx
	conn.pktArray = make([]*Packet, cap, cap)
	//conn.array = make([]*Conn)
	return conn
}

func (n *Network) NewInitializationConnection() *InitializationConnection {
	conn := &InitializationConnection{
		network: n,
	}

	return conn
}

func (n *Network) NewInArrayPort() *InArrayPort {
	conn := &InArrayPort{
		network: n,
	}
	conn.array = make([]*Connection, 20) // set to 20 for now...
	return conn
}

func (n *Network) Connect(p1 *Process, out string, p2 *Process, in string, cap int) {

	//var i int
	var connxn *Connection
	var anyConn Conn
	if strings.Index(in, "[") > -1 {
		//match := ""
		re := regexp.MustCompile(`(.*)\[(\d*)\]`)
		parts := re.FindStringSubmatch(in)
		root := parts[1]
		indx := parts[2]
		i, err := strconv.Atoi(indx)
		if err != nil {
			panic(fmt.Sprintf("Invalid index: ", indx))
		}
		anyConn = p2.inPorts[root]
		if anyConn == nil {
			anyConn = n.NewInArrayPort()
			p2.inPorts[root] = anyConn
		}

		connxn = anyConn.GetArrayItem(i)

		if connxn == nil {
			connxn := n.NewConnection(cap)
			connxn.portName = in
			connxn.fullName = p2.Name + "." + in
			if anyConn == nil {
				p2.inPorts[in] = connxn
			} else {
				anyConn.SetArrayItem(connxn, i)
			}
		}

	} else {
		connxn = p2.inPorts[in].(*Connection)
	}
	opt := p1.outPorts[out]
	if opt != nil {
		panic("Outport port already connected")
	}
	opt = new(OutPort)
	p1.outPorts[out] = opt
	opt.name = out
	//conn := connxn
	opt.Conn = connxn
	connxn.incUpstream()
}

func (n *Network) Initialize(initValue string, p2 *Process, in string) {

	conn := n.NewInitializationConnection()
	p2.inPorts[in] = conn
	conn.portName = in
	conn.fullName = p2.Name + "." + in

	conn.value = initValue

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
