package core

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

var pStack []*Process

var stkLevel int

var tracing bool
var tracelocks bool

type Network struct {
	Name   string
	procs  map[string]*Process
	wg     sync.WaitGroup
	mother *Process
}

func NewNetwork(name string) *Network {
	net := &Network{
		Name:  name,
		procs: make(map[string]*Process),
		wg:    sync.WaitGroup{},
	}

	//stkLevel++
	//if stkLevel >= len(pStack) {
	pStack = append(pStack, nil)
	//}

	return net
}

func NewSubnet(Name string, p *Process) *Network {
	net := &Network{
		Name:  Name,
		procs: make(map[string]*Process),
		wg:    sync.WaitGroup{},
	}

	net.mother = p
	pStack[stkLevel] = p
	stkLevel++
	if stkLevel >= len(pStack) {
		pStack = append(pStack, nil)
	}
	return net
}

func (n *Network) GetProc(nm string) *Process {
	return n.procs[nm]
}

func (n *Network) SetProc(p *Process, nm string) {
	n.procs[nm] = p
}

func LockTr(sc *sync.Cond, s string, p *Process) {
	sc.L.Lock()
	if tracelocks && p != nil {
		fmt.Println(p.Name, s)
	}
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
	sc.Wait()
	if tracelocks && p != nil {
		fmt.Println(p.Name, s)
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

	n.SetProc(proc, nm)
	proc.inPorts = make(map[string]inputCommon)
	proc.outPorts = make(map[string]outputCommon)
	proc.mtx = sync.Mutex{}
	proc.canGo = sync.NewCond(&proc.mtx)
	if stkLevel > 0 {
		n.mother = pStack[stkLevel-1]
	}

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
			//connxn.fullName = p2.Name + "." + inPort.name + "[" + strconv.Itoa(inPort.index) + "]"
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
			//connxn.fullName = p2.Name + "." + inPort.name
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
		outConn.SetArrayItem(opt, outPort.index)
		opt.portName = out
		//opt.fullName = p1.Name + "." + out
	} else {
		//var opt OutputConn
		opt = new(OutPort)
		p1.outPorts[out] = opt
		opt.network = n
		opt.portName = out
		//opt.fullName = p1.Name + "." + out

	}

	opt.sender = p1
	opt.Conn = connxn
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
	//conn.fullName = p2.Name + "." + in

	conn.value = initValue
}

func (n *Network) Exit() {
	if n.mother == nil {
		if stkLevel != 0 {
			panic("Exit - stack level incorrect: " + strconv.Itoa(stkLevel))
		}
	} else {
		stkLevel--
	}
}

func trace(p *Process, s ...string) {
	if tracing {
		fmt.Print(p.Name, " "+strings.Trim(fmt.Sprint(s), "[]")+"\n")
	}
}

func setOptions() {
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
			tracing = true
		}

		i = strings.Index(rec, "<tracelocks>")
		if i > -1 && rec[i+12:i+16] == "true" {
			tracelocks = true
		}
	}
}

func (n *Network) Run() {
	defer n.Exit()
	if n.mother == nil {
		setOptions()
	}
	defer fmt.Println(n.Name + " Done")

	for {
		if n.mother != nil {
			fmt.Println(n.Name + " Starting subnet activation")
		} else {
			fmt.Println(n.Name + " Starting network")
		}

		// FBP distinguishes between execution of the process as a whole and activating the code - the code may be deactivated and then
		// reactivated many times during the process "run"

		n.wg = sync.WaitGroup{}
		n.wg.Add(len(n.procs))

		//defer n.wg.Wait()
		var someProcsCanRun bool = false
		//time.Sleep(1 * time.Millisecond)
		for _, proc := range n.procs {
			//proc.status = Notstarted
			atomic.StoreInt32(&proc.status, Notstarted)
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

			proc.activate()
			someProcsCanRun = true
		}
		if !someProcsCanRun {
			n.wg.Add(0 - len(n.procs))
			panic("No process can start")
		}
		n.wg.Wait()

		if n.mother == nil {
			return
		}
		fmt.Println(n.Name + " Subnet deactivated")

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
		//stkLevel--

		p := pStack[stkLevel-1]

		allDrained, _, _ := p.inputState()
		if allDrained {
			break
		}
		fmt.Println("XXXXXXXXXXXXXXXXXXXXXXXXXXX")
		fmt.Println(n.Name)
		for _, p := range n.procs {
			fmt.Println(p.Name)
			for _, v := range p.inPorts {
				_, b := v.(*InArrayPort)
				if b {
					for _, w := range v.(*InArrayPort).array {
						//w.resetForNextExecution()
						w.upStrmCnt = 0
					}
				} else {
					w, b := v.(*InPort)
					if b {
						//w.resetForNextExecution()
						w.upStrmCnt = 0
					}
				}
			}

			for _, v := range p.outPorts {
				_, b := v.(*OutArrayPort)
				if b {
					for _, w := range v.(*OutArrayPort).array {
						w.Conn.incUpstream()
					}
				} else {
					w, b := v.(*OutPort)
					if b {
						w.Conn.incUpstream()
					}
				}
			}
		}
	}

	//stkLevel++
}
