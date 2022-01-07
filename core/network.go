package core

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"sync"
	"sync/atomic"
)

var tracing bool
var tracelocks bool
var generateGids bool

type Network struct {
	Name  string
	procs map[string]*Process
	//conns  map[string]inputCommon
	wg     sync.WaitGroup
	mother *Process
}

func NewNetwork(name string) *Network {
	net := &Network{
		Name:  name,
		procs: make(map[string]*Process),
		wg:    sync.WaitGroup{},
	}

	return net
}

func NewSubnet(Name string, p *Process) *Network {
	net := &Network{
		Name:  Name,
		procs: make(map[string]*Process),
		wg:    sync.WaitGroup{},
	}

	net.mother = p
	return net
}

func (n *Network) GetProc(nm string) *Process {
	return n.procs[nm]
}

func (n *Network) SetProc(p *Process, nm string) {
	n.procs[nm] = p
}

func LockTr(sc *sync.Cond, s string, p *Process) {
	if tracelocks && p != nil {
		fmt.Println(p.Name, s)
	}
	sc.L.Lock()
}

func UnlockTr(sc *sync.Cond, s string, p *Process) {
	sc.L.Unlock()
	if tracelocks && p != nil {
		fmt.Println(p.Name, s)
	}
}

func BdcastTr(sc *sync.Cond, s string, p *Process) {
	sc.Broadcast()
	if tracelocks && p != nil {
		fmt.Println(p.Name, s)
	}
}

func WaitTr(sc *sync.Cond, s string, p *Process) {
	if tracelocks && p != nil {
		fmt.Println(p.Name, s)
	}
	sc.Wait()
}

func trace(p *Process, s string) {
	if tracing {
		//fmt.Print(p.Name, " "+strings.Trim(fmt.Sprint(s), "[]")+"\n")
		fmt.Print(p.Name, " "+s+"\n")
	}
}

func traceNet(n *Network, s string) {
	if tracing {
		//fmt.Print(n.Name, " "+strings.Trim(fmt.Sprint(s), "[]")+"\n")
		fmt.Print(n.Name, " "+s+"\n")
	}
}

//func (n *Network) id() string { return fmt.Sprintf("%p", n) }

func (n *Network) NewProc(nm string, comp Component) *Process {

	proc := &Process{
		Name:      nm,
		logFile:   "",
		component: comp,
		status:    Notstarted,
		network:   n,
	}
	//atomic.StoreInt32(&proc.status, Notstarted)
	n.SetProc(proc, nm)
	proc.inPorts = make(map[string]inputCommon)
	proc.outPorts = make(map[string]outputCommon)
	proc.mtx = sync.Mutex{}
	proc.canGo = sync.NewCond(&proc.mtx)
	//if stkLevel > 0 {
	//	n.mother = pStack[stkLevel-1]
	//}
	proc.gid = getGID()
	return proc
}

func (n *Network) NewConnection(cap int) *InPort {
	conn := &InPort{
		network: n,
	}
	conn.mtx = sync.Mutex{}
	conn.condNE = sync.NewCond(&conn.mtx)
	conn.condNF = sync.NewCond(&conn.mtx)
	conn.pktArray = make([]*Packet, cap)
	return conn
}

func (n *Network) NewInitializationConnection() *InitializationConnection {
	conn := &InitializationConnection{
		network: n,
	}
	conn.mtx = sync.Mutex{}
	return conn
}

func (n *Network) NewInArrayPort() *InArrayPort {
	conn := &InArrayPort{
		network: n,
	}
	conn.mtx = sync.Mutex{}
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

		connxn = anyInConn.(InputArrayConn).GetArrayItem(inPort.index)

		if connxn == nil {
			connxn = n.NewConnection(cap)
			connxn.portName = inPort.name + "[" + strconv.Itoa(inPort.index) + "]"
			connxn.fullName = p2.Name + "." + connxn.portName
			//n.conns[connxn.fullName] = connxn
			connxn.downStrProc = p2
			connxn.network = n
			if anyInConn == nil {
				p2.inPorts[inPort.name] = connxn
			} else {
				anyInConn.(InputArrayConn).setArrayItem(connxn, inPort.index)
			}
		}
	} else {
		if p2.inPorts[inPort.name] == nil {
			connxn = n.NewConnection(cap)
			connxn.portName = inPort.name
			connxn.fullName = p2.Name + "." + inPort.name
			//n.conns[connxn.fullName] = connxn
			connxn.downStrProc = p2
			connxn.network = n
			p2.inPorts[inPort.name] = connxn
		} else {
			connxn = p2.inPorts[inPort.name].(*InPort)
		}
	}

	if in == "*" {
		p2.autoInput = connxn
	}

	// connxn built; input port array built if necessary

	outPort := parsePort(out)

	var opt *OutPort

	if outPort.indexed {
		var anyOutConn = p1.outPorts[outPort.name]
		if anyOutConn == nil {
			anyOutConn = n.NewOutArrayPort()
			p1.outPorts[outPort.name] = anyOutConn
		}

		//opt := new(OutArrayPort)
		outConn := anyOutConn.(*OutArrayPort)
		//p1.outPorts[out] = anyOutConn
		//opt.Name = out
		opt = new(OutPort)
		outConn.setArrayItem(opt, outPort.index)
		opt.portName = out
		opt.fullName = p1.Name + "." + out
		//n.conns[opt.fullName] = opt.Conn
	} else {
		//var opt OutputConn
		opt = new(OutPort)
		p1.outPorts[out] = opt
		opt.network = n
		opt.portName = out
		opt.fullName = p1.Name + "." + out
		//n.conns[opt.fullName] = opt.Conn

	}

	opt.sender = p1
	opt.conn = connxn
	opt.connected = true

	connxn.incUpstream()
	if out == "*" {
		p1.autoOutput = connxn
	}
}

type portDefinition struct {
	name    string
	index   int
	indexed bool
}

var portPattern = regexp.MustCompile(`^(.+)\[(\d+)\]$`)

func parsePort(in string) portDefinition {
	matches := portPattern.FindStringSubmatch(in)
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

func (n *Network) Initialize(initValue interface{}, p2 *Process, in string) {

	conn := n.NewInitializationConnection()
	p2.inPorts[in] = conn
	conn.portName = in
	conn.fullName = p2.Name + "." + in
	//n.conns[conn.fullName] = conn

	conn.value = initValue
}

func (n *Network) Exit() {
	if n.mother == nil {
		traceNet(n, "Exit network")
	} else {
		traceNet(n, "Exit subnet")
	}
}

func setOptions() {
	var params struct {
		Tracing      bool `xml:"tracing"`
		TraceLocks   bool `xml:"tracelocks"`
		GenerateGIDs bool `xml:"generate-gIds"`
	}

	xmldata, err := ioutil.ReadFile("params.xml")
	if err != nil {
		return
	}

	err = xml.Unmarshal(xmldata, &params)
	if err != nil {
		return

	}

	tracing = params.Tracing
	tracelocks = params.TraceLocks
	generateGids = params.GenerateGIDs
}

func (n *Network) Run() {
	defer n.Exit()
	if n.mother == nil {
		setOptions()
	}

	defer traceNet(n, " Done")

	for {
		if n.mother != nil {
			traceNet(n, " Starting subnet activation")
		} else {
			traceNet(n, " Starting network")
		}

		// FBP distinguishes between execution of the process as a whole and activating the code - the code may be deactivated and then
		// reactivated many times during the process "run"

		n.wg = sync.WaitGroup{}
		n.wg.Add(len(n.procs))

		//defer n.wg.Wait()
		var someProcsCanRun bool = false

		for _, proc := range n.procs {

			//LockTr(proc.canGo, "test if not started L", proc)
			//atomic.StoreInt32(&proc.status, Notstarted)

			selfStarting := true
			if proc.inPorts != nil {
				for _, conn := range proc.inPorts {

					_, b := conn.(*InitializationConnection)
					if !b {
						selfStarting = false
					}
				}
			}

			if !selfStarting {
				continue
			}
			//atomic.StoreInt32(&proc.status, Notstarted)
			trace(proc, "act from start")
			proc.activate()
			someProcsCanRun = true
			//UnlockTr(proc.canGo, "test if not started U", proc)
		}
		if !someProcsCanRun {
			n.wg.Add(0 - len(n.procs))
			panic("No process can start")
		}
		n.wg.Wait()

		if n.mother == nil {
			return
		}

		traceNet(n, " subnet deactivated")

		for _, p := range n.procs {
			for _, v := range p.inPorts {
				_, b := v.(*InArrayPort)
				if b {
					for _, w := range v.(*InArrayPort).array {
						w.resetForNextExecution()
					}
				} else {
					v.resetForNextExecution()
				}
			}
		}

		p := n.mother

		allDrained, _, _ := p.inputState()
		if allDrained {
			break
		}

		for _, p := range n.procs {
			atomic.StoreInt32(&p.status, Notstarted)
			for _, v := range p.outPorts {
				_, b := v.(*OutArrayPort)
				if b {
					for _, w := range v.(*OutArrayPort).array {
						w.conn.upStrmCnt = 0
					}
				} else {
					w, b := v.(*OutPort)
					if b {
						w.conn.upStrmCnt = 0
					}
				}
			}
		}

		for _, p := range n.procs {
			for _, v := range p.outPorts {
				_, b := v.(*OutArrayPort)
				if b {
					for _, w := range v.(*OutArrayPort).array {
						w.conn.incUpstream()
					}
				} else {
					w, b := v.(*OutPort)
					if b {
						w.conn.incUpstream()
					}
				}
			}
		}
	}

	//if n.mother != nil {
	p := n.mother
	for _, v := range p.outPorts {
		_, b := v.(*OutArrayPort)
		if b {
			for _, w := range v.(*OutArrayPort).array {
				w.Close()
			}
		} else {
			w, b := v.(*OutPort)
			if b {
				w.Close()
			}
		}
	}
}

//}
