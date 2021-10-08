package core

import (
	"fmt"
	"regexp"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	Notstarted int32 = iota
	Active
	Dormant
	SuspSend
	SuspRecv
	Terminated
)

type Network struct {
	Name  string
	procs map[string]*Process
	//procList []*Process
	//driver  Process
	logFile string
	wg      sync.WaitGroup
}

func NewNetwork(name string) *Network {
	net := &Network{
		Name:  name,
		procs: make(map[string]*Process),
	}

	return net
}

func (n *Network) NewProc(nm string, comp Component) *Process {

	proc := &Process{
		name:      nm,
		network:   n,
		logFile:   "",
		component: comp,
		status:    Notstarted,
	}

	//n.procList = append(n.procList, proc)
	n.procs[nm] = proc

	proc.inPorts = make(map[string]inputCommon)
	proc.outPorts = make(map[string]outputCommon)

	return proc
}

func (n *Network) NewConnection(cap int) *Connection {
	conn := &Connection{
		network: n,
	}
	conn.condNE.L = &conn.mtx
	conn.condNF.L = &conn.mtx
	conn.pktArray = make([]*Packet, cap)
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

	var conn InputConn

	if inPort.indexed {
		// try get an existing port
		existingConn := p2.inPorts[inPort.name]
		if existingConn == nil {
			// create one, if it doesn't exist
			existingConn = n.NewInArrayPort()
			p2.inPorts[inPort.name] = existingConn
		}
		// enforce it's an array
		arrayConn := existingConn.(InputArrayConn)

		// try get an existing connection from array
		conn = arrayConn.GetArrayItem(inPort.index)
		if conn == nil {
			// create one if it doesn't exist
			newConn := n.NewConnection(cap)
			newConn.portName = inPort.name
			newConn.fullName = p2.name + "." + inPort.name
			newConn.downStrProc = p2
			newConn.network = n
			conn = newConn
			arrayConn.SetArrayItem(conn, inPort.index)
		}
	} else {
		// try get an existing port
		existingConn, ok := p2.inPorts[inPort.name]
		if !ok {
			// create one if it doesn't exist
			newConn := n.NewConnection(cap)
			newConn.portName = inPort.name
			newConn.fullName = p2.name + "." + inPort.name
			newConn.downStrProc = p2
			newConn.network = n
			existingConn = newConn
			p2.inPorts[inPort.name] = newConn
		}
		conn = existingConn.(InputConn)
	}

	// connxn built; input port array built if necessary

	// rest of the code requires that the input is a `*Connection`
	// this probably could be unrestricted.
	connxn := conn.(*Connection)

	outPort := parsePort(out)
	if outPort.indexed {
		// try get an existing port
		anyOutConn, ok := p1.outPorts[outPort.name]
		if !ok {
			// create one if it doesn't exist
			anyOutConn = n.NewOutArrayPort()
			p1.outPorts[outPort.name] = anyOutConn
		}
		// enforce the output is an array
		arrayConn := anyOutConn.(OutputArrayConn)

		// check that nothing has connected to that item
		if existing := arrayConn.GetArrayItem(outPort.index); existing != nil {
			panic("output port " + out + " already connected")
		}

		// update the array
		arrayConn.SetArrayItem(&OutPort{
			name: out,
			conn: connxn,
		}, outPort.index)

	} else {
		// check that nothing has connected already
		if _, exists := p1.outPorts[out]; exists {
			panic("output port " + out + " already connected")
		}
		// add the connection
		p1.outPorts[out] = &OutPort{
			name: out,
			conn: connxn,
		}
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

	// FBP distinguishes between execution of the process as a whole and activating the code - the code may be deactivated and then
	// reactivated many times during the process "run"

	n.wg.Add(len(n.procs))

	defer n.wg.Wait()

	go func() {
		for {
			time.Sleep(200 * time.Millisecond)
			allTerminated := true
			deadlockDetected := true
			for _, proc := range n.procs {
				//proc.mtx.Lock()
				//defer proc.mtx.Unlock()
				status := atomic.LoadInt32(&proc.status)
				if status != Terminated {
					allTerminated = false
					if status == Active {
						deadlockDetected = false
					}
				}
			}
			if allTerminated {
				//	fmt.Println("Run terminated")
				return
			}
			if deadlockDetected {
				fmt.Println("\nDeadlock detected!")
				for key, proc := range n.procs {
					fmt.Println(key, " Status: ",
						[]string{"notStarted",
							"active",
							"dormant",
							"suspSend",
							"suspRecv",
							"terminated"}[proc.status])
				}
				panic("Deadlock!")
			}

		}
	}()

	startedAnything := false
	for _, proc := range n.procs {
		if proc.isSelfStarting() {
			proc.ensureRunning()
			startedAnything = true
		}
	}

	if !startedAnything {
		n.wg.Add(0 - len(n.procs))
		panic("No process can start")
	}
}
