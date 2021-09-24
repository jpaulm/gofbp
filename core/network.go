package core

import (
	"fmt"
	"regexp"
	"strconv"
	"sync"
	"unsafe"
)

const (
	Notstarted int32 = iota
	Dormant
	SuspSend
	SuspRecv
	Active
	Terminated
)

var wg sync.WaitGroup

type Network struct {
	Name  string
	procs map[string]*Process
	//procList []*Process
	//driver  Process
	logFile string
	wg      *sync.WaitGroup
}

func NewNetwork(name string) *Network {
	net := &Network{
		Name: name,
	}

	net.wg = &wg

	ptr := unsafe.Pointer(&net.wg)
	fmt.Println(ptr)

	net.procs = make(map[string]*Process)

	return net
}

func (n *Network) NewProc(nm string, comp Component) *Process {

	proc := &Process{
		name:      nm,
		network:   n,
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
			connxn.fullName = p2.name + "." + inPort.name
			connxn.downStrProc = p2
			connxn.network = n
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
			connxn.fullName = p2.name + "." + inPort.name
			connxn.downStrProc = p2
			connxn.network = n
			p2.inPorts[inPort.name] = connxn
		} else {
			connxn = p2.inPorts[inPort.name].(*Connection)
		}
	}

	// connxn built; input port array built if necessary

	var anyOutConn OutputConn

	outPort := parsePort(out)

	if outPort.indexed {
		anyOutConn = p1.outPorts[outPort.name]
		if anyOutConn == nil {
			anyOutConn = n.NewOutArrayPort()
			p1.outPorts[outPort.name] = anyOutConn
		}

		opt := new(OutPort)
		//p1.outPorts[out] = anyOutConn
		opt.name = out
		anyOutConn.SetArrayItem(opt, outPort.index)
		opt.Conn = connxn

	} else {
		//var opt OutputConn
		opt := new(OutPort)
		p1.outPorts[out] = opt
		opt.name = out
		opt.Conn = connxn
	}

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
	conn.fullName = p2.name + "." + in

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

	defer func() {
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
	}()

	for _, proc := range n.procs {

		proc.status = Notstarted
		//proc.network = n
		proc.starting = true
		if proc.inPorts != nil && !isMustRun(proc.component) {
			for _, conn := range proc.inPorts {
				if conn.GetType() != "InitializationConnection" {
					proc.starting = false
				}
			}
		}
		if !proc.starting {
			continue
		}

		proc.ensureRunning()
	}
}

func isMustRun(comp Component) bool {
	_, hasMustRun := comp.(ComponentWithMustRun)
	if hasMustRun {
		fmt.Printf("%T component has MustRun method\n", comp)
		return true
	}
	return false
}
