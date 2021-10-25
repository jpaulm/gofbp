package core

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	//"sync/atomic"
	//"time"
)

const (
	Notstarted int32 = iota
	Active
	Dormant
	SuspSend
	SuspRecv
	Terminated
)

type Stack []*Process

var pStack Stack

var stkLevel int

type Network struct {
	Name  string
	procs map[string]*Process
	//procList []*Process
	//driver  Process
	//logFile string
	wg      sync.WaitGroup
	Active  int32
	tracing bool
}

func NewNetwork(name string) *Network {
	net := &Network{
		Name:  name,
		procs: make(map[string]*Process),
	}

	//stkLevel++
	if stkLevel >= len(pStack) {
		pStack = append(pStack, nil)
	}

	return net
}

func NewSubnet(name string, p *Process) *Network {
	net := &Network{
		Name:  name,
		procs: make(map[string]*Process),
	}

	//stkLevel++
	if stkLevel >= len(pStack) {
		pStack = append(pStack, nil)
	}
	pStack[stkLevel] = p
	stkLevel++
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

	n.procs[nm] = proc

	proc.inPorts = make(map[string]inputCommon)
	proc.outPorts = make(map[string]outputCommon)
	if stkLevel > 0 {
		proc.Mother = pStack[stkLevel-1]
	}
	//pStack[stkLevel] = proc

	return proc
}

func (n *Network) NewConnection(cap int) *InPort {
	conn := &InPort{
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
	port := &OutArrayPort{
		network: n,
	}

	return port
}

func (n *Network) Connect(p1 *Process, out string, p2 *Process, in string, cap int) {

	inPort := parsePort(in)

	var connxn *InPort
	//var anyInConn InputConn

	if inPort.indexed {
		var anyInConn = p2.inPorts[inPort.name]
		if anyInConn == nil {

			anyInConn = n.NewInArrayPort()
			p2.inPorts[inPort.name] = anyInConn
		}

		//anyInConn = anyInConn.(InputArrayConn)
		connxn = anyInConn.(InputArrayConn).GetArrayItem(inPort.index)

		if connxn == nil {
			connxn = n.NewConnection(cap)
			connxn.portName = inPort.name
			connxn.fullName = p2.name + "." + inPort.name
			connxn.downStrProc = p2
			connxn.network = n
			if anyInConn == nil {
				p2.inPorts[inPort.name] = connxn
			} else {
				anyInConn.(InputArrayConn).SetArrayItem(connxn, inPort.index)
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
			connxn = p2.inPorts[inPort.name].(*InPort)
		}
	}

	// connxn built; input port array built if necessary

	//var anyOutConn OutputConn

	outPort := parsePort(out)

	if outPort.indexed {
		var anyOutConn = p1.outPorts[outPort.name]
		if anyOutConn == nil {
			anyOutConn = n.NewOutArrayPort()
			p1.outPorts[outPort.name] = anyOutConn
		}

		//opt := new(OutArrayPort)
		out := anyOutConn.(*OutArrayPort)
		//p1.outPorts[out] = anyOutConn
		//opt.name = out
		opt := new(OutPort)
		out.SetArrayItem(opt, outPort.index)
		opt.Conn = connxn
		opt.connected = true

	} else {
		//var opt OutputConn
		opt := new(OutPort)
		p1.outPorts[out] = opt
		opt.name = out
		opt.Conn = connxn
		opt.connected = true
		//fmt.Println(opt)
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
		panic("Invalid index in " + in)
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
func (n *Network) Exit() {
	if stkLevel > 0 {
		stkLevel--
	}
}

func (n *Network) trace(s ...string) {
	if n.tracing {
		fmt.Print(strings.Trim(fmt.Sprint(s), "[]") + "\n")
	}
}

// Deadlock detection goroutine has been commented out...

func (n *Network) Run() {

	var rec string

	f, err := os.Open("params.xml")
	if err == nil {
		defer f.Close()
		buf := make([]byte, 1024)
		for {
			n, err := f.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println(err)
				continue
			}
			if n > 0 {
				//fmt.Println(string(buf[:n]))
				rec += string(buf[:n])
			}
		}

		i := strings.Index(rec, "<tracing>")
		if i > -1 && rec[i+9:i+13] == "true" {
			n.tracing = true
		}
	}

	// Criterion being used for deadlock: no process has become or is already active in last 200 ms

	/*
		go func(n *Network) {

			for {
				atomic.StoreInt32(&n.Active, 0)
				time.Sleep(2 * time.Millisecond) // shd be 200 ms!
				//atomic.StoreInt32(&n.active, 0)
				allTerminated := true
				//deadlockDetected := true
				for _, proc := range n.procs {
					proc.mtx.Lock()
					//defer proc.mtx.Unlock()
					status := atomic.LoadInt32(&proc.status)
					if status != Terminated {
						allTerminated = false
						if status == Active {
							atomic.StoreInt32(&n.Active, 1) // in case 200 ms go by without a status change...
						}
					}
					proc.mtx.Unlock()
				}
				if allTerminated {
					//fmt.Println(n.Name, " terminated")
					return
				}
				//if deadlockDetected {
				if n.Active == 0 {
					fmt.Println("\nDeadlock detected in", n.Name+"!")
					for key, proc := range n.procs {
						proc.mtx.Lock()
						//defer proc.mtx.Unlock()
						status := atomic.LoadInt32(&proc.status)
						fmt.Println(" ",
							[]string{"NotStarted:",
								"Active:    ",
								"Dormant:   ",
								"SuspSend:  ",
								"SuspRecv:  ",
								"Terminated:"}[status], key)
						proc.mtx.Unlock()
					}
					panic("Deadlock!")
				}

			}

		}(n)
	*/

	defer n.Exit()
	defer fmt.Println(n.Name + " Done")
	fmt.Println(n.Name + " Starting")

	// FBP distinguishes between execution of the process as a whole and activating the code - the code may be deactivated and then
	// reactivated many times during the process "run"

	n.wg.Add(len(n.procs))

	defer n.wg.Wait()
	var someProcsCanRun bool = false
	//time.Sleep(1 * time.Millisecond)
	for _, proc := range n.procs {
		proc.mtx.Lock()
		defer proc.mtx.Unlock()
		proc.selfStarting = true
		if proc.inPorts != nil {
			for _, conn := range proc.inPorts {
				//if conn.GetType() != "InitializationPort" {
				_, b := conn.(*InitializationConnection)
				if !b {
					proc.selfStarting = false
				}
			}
		}
		if !proc.selfStarting {
			continue
		}

		proc.ensureRunning()
		someProcsCanRun = true
	}
	if !someProcsCanRun {
		n.wg.Add(0 - len(n.procs))
		panic("No process can start")
	}
	//}()
}
