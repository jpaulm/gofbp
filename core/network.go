package core

import (
	"fmt"
	"regexp"
	"strconv"
	"sync"
)

type Status int

const (
	notStarted Status = iota
	dormant
	suspSend
	suspRecv
	active
	terminated
)

type Network struct {
	Name  string
	procs map[string]*Process
	//procList []*Process
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

	//n.procList = append(n.procList, proc)
	n.procs[nm] = proc

	proc.inPorts = make(map[string]InputConn)
	proc.outPorts = make(map[string]OutputConn)

	return proc
}

func (n *Network) NewConnection(cap int) *Connection {
	conn := &Connection{
		network: n,
	}
	conn.condNE.L = &conn.mtx
	conn.condNF.L = &conn.mtx
	conn.pktArray = make([]*Packet, cap, cap)
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

	return conn
}

func (n *Network) NewOutArrayPort() *OutArrayPort {
	conn := &OutArrayPort{
		network: n,
	}

	return conn
}

func (n *Network) Connect(p1 *Process, out string, p2 *Process, in string, cap int) {

	inPort := parsePort(in)
	outPort := parsePort(out)

	var connxn *Connection
	var anyInConn InputConn

	if inPort.indexed {
		anyInConn = p2.inPorts[inPort.name]
		if anyInConn == nil {
			anyInConn = n.NewInArrayPort()
			p2.inPorts[inPort.name] = anyInConn
		}

		connxn = anyInConn.GetArrayItem(inPort.index)

		if connxn == nil {
			connxn = n.NewConnection(cap)
			connxn.portName = inPort.name
			connxn.fullName = p2.Name + "." + inPort.name
			if anyInConn == nil {
				p2.inPorts[inPort.name] = connxn
			} else {
				anyInConn.SetArrayItem(connxn, inPort.index)
			}
		}
	} else {
		if p2.inPorts[inPort.name] == nil {
			connxn = n.NewConnection(cap)
			connxn.portName = inPort.name
			connxn.fullName = p2.Name + "." + inPort.name
			p2.inPorts[inPort.name] = connxn
		} else {
			connxn = p2.inPorts[inPort.name].(*Connection)
		}
	}

	// connxn built; input port array built if necessary

	var anyOutConn OutputConn

	if outPort.indexed {
		anyOutConn = p1.outPorts[outPort.name]
		if anyOutConn == nil {
			anyOutConn = n.NewOutArrayPort()
			p1.outPorts[outPort.name] = anyOutConn
		}
		opt := new(OutPort)
		p1.outPorts[out] = anyOutConn
		opt.name = out
		anyOutConn.SetArrayItem(opt, outPort.index)
		opt.Conn = connxn
		//connxn = anyOutConn.GetArrayItem(outPort.index)

		/*
			if connxn == nil {
				connxn = n.NewConnection(cap)
				connxn.portName = inPort.name
				connxn.fullName = p2.Name + "." + inPort.name
				if anyOutConn == nil {
					p2.inPorts[inPort.name] = connxn
				} else {
					anyOutConn.SetArrayItem(connxn, inPort.index)
				}
			}
		*/
	} else {
		//var opt OutputConn
		opt := new(OutPort)
		p1.outPorts[out] = opt
		opt.name = out
		opt.Conn = connxn

		/*
			if p1.outPorts[outPort.name] == nil {
				connxn = n.NewConnection(cap)
				connxn.portName = outPort.name
				connxn.fullName = p1.Name + "." + outPort.name
				p1.outPorts[outPort.name] = connxn
			} else {
				connxn = p1.outPorts[outPort.name].(*Connection)
			}
		*/
	}

	/*
		opt := p1.outPorts[out]
		if opt != nil {
			panic("Outport port already connected")
		}
		opt = new(OutPort)
		p1.outPorts[out] = opt
		opt.name = out
	*/
	//conn := connxn
	//opt.Conn = connxn
	connxn.incUpstream()
}

type portDefinition struct {
	name    string
	index   int
	indexed bool
}

var rePort = regexp.MustCompile(`^(.+)\[(\d+)\]$`)

func parsePort(in string) portDefinition {
	matches := rePort.FindStringSubmatch(in)
	if len(matches) == 0 {
		return portDefinition{name: in}
	}
	root, indexStr := matches[1], matches[2]

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		panic(fmt.Sprintf("Invalid index in %q", in))
	}

	return portDefinition{name: root, index: index, indexed: true}
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
	//defer wg.Wait()

	// FBP distinguishes between execution of the process as a whole and activating the code - the code may be deactivated and then
	// reactivated many times during the process "run"

	wg.Add(len(n.procs))
	for _, proc := range n.procs {
		proc := proc
		//wg.Add(1)
		go func() { // Process goroutine
			defer wg.Done()
			//if len(proc.inPorts) == 0 {
			proc.Run(n)
			//}
		}()
	}

	wg.Wait()

	for key, proc := range n.procs {
		fmt.Println(key, " Status: ",
			[]string{"notStarted",
				"dormant",
				"suspSend",
				"suspRecv",
				"active",
				"terminated"}[proc.status])
	}
}
