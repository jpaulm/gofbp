package core

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var pStack []*Process

var stkLevel int
var tracing bool
var tracelocks bool

type GenNet interface {
	id() string
	NewConnection(int) *InPort
	NewInitializationConnection() *InitializationConnection
	NewInArrayPort() *InArrayPort
	NewOutArrayPort() *OutArrayPort
	Connect(*Process, string, *Process, string, int)
	Initialize(interface{}, *Process, string)
	Exit()
	GetProc(string) *Process
	SetProc(*Process, string)
	Run()
}

type Network struct {
	name  string
	procs map[string]*Process
	wg    sync.WaitGroup
}

//func (n *Network) GetNetwork() *Network {
//	return n
//}

func NewNetwork(name string) *Network {
	net := &Network{
		name:  name,
		procs: make(map[string]*Process),
		wg:    sync.WaitGroup{},
	}

	//stkLevel++
	//if stkLevel >= len(pStack) {
	pStack = append(pStack, nil)
	//}

	return net
}

func NewSubnet(name string, p *Process) *Subnet {
	net := &Subnet{}
	//stkLevel++
	//if stkLevel >= len(pStack) {
	//	pStack = append(pStack, nil)
	//}
	net.name = name
	net.procs = make(map[string]*Process)
	net.wg = sync.WaitGroup{}
	net.SetMother(p)
	pStack[stkLevel] = p
	stkLevel++
	if stkLevel >= len(pStack) {
		pStack = append(pStack, nil)
	}
	return net
}

func LockTr(sc *sync.Cond, s string, p *Process) {
	sc.L.Lock()
	if tracelocks {
		fmt.Println(p.GetName(), s)
	}
}

func UnlockTr(sc *sync.Cond, s string, p *Process) {
	sc.L.Unlock()
	if tracelocks {
		fmt.Println(p.GetName(), s)
	}
}

func BdcastTr(sc *sync.Cond, s string, p *Process) {
	sc.Broadcast()
	if tracelocks {
		fmt.Println(p.GetName(), s)
	}
}

func WaitTr(sc *sync.Cond, s string, p *Process) {
	sc.Wait()
	if tracelocks {
		fmt.Println(p.GetName(), s)
	}
}

func (n *Network) id() string { return fmt.Sprintf("%p", n) }

func (n *Network) GetProc(nm string) *Process {
	return n.procs[nm]
}

func (n *Network) SetProc(p *Process, nm string) {
	n.procs[nm] = p
}

func (n *Network) NewProc(nm string, comp Component) *Process {

	proc := &Process{
		name:      nm,
		logFile:   "",
		component: comp,
		status:    Notstarted,
		network:   n,
	}

	//proc.network = n
	n.SetProc(proc, nm)
	proc.inPorts = make(map[string]inputCommon)
	proc.outPorts = make(map[string]outputCommon)
	proc.mtx = sync.Mutex{}
	proc.canGo = sync.NewCond(&proc.mtx)

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
			//connxn.portName = inPort.name
			connxn.name = p2.name + "." + inPort.name + "[" + strconv.Itoa(inPort.index) + "]"
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
			//connxn.portName = inPort.name
			connxn.name = p2.name + "." + inPort.name
			connxn.downStrProc = p2
			connxn.network = n
			p2.inPorts[inPort.name] = connxn
		} else {
			connxn = p2.inPorts[inPort.name].(*InPort)
		}
	}

	if inPort.name == "*" {
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
		//opt.name = out
		opt = new(OutPort)
		outConn.SetArrayItem(opt, outPort.index)
		opt.name = p1.name + "." + out
	} else {
		//var opt OutputConn
		opt = new(OutPort)
		p1.outPorts[out] = opt
		opt.name = p1.name + "." + out

	}

	opt.SetSender(p1)
	opt.Conn = connxn
	opt.connected = true

	connxn.incUpstream()
	if outPort.name == "*" {
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
	//conn.portName = in
	conn.fullName = p2.name + "." + in

	conn.value = initValue

}
func (n *Network) Exit() {
	if stkLevel != 0 {
		panic("Exit - stack level incorrect: " + strconv.Itoa(stkLevel))
	}
}

func trace(s ...string) {
	if tracing {
		fmt.Print(strings.Trim(fmt.Sprint(s), "[]") + "\n")
	}
}


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
			tracing = true
		}

		i = strings.Index(rec, "<tracelocks>")
		if i > -1 && rec[i+12:i+16] == "true" {
			tracelocks = true
		}
	}

	defer n.Exit()
	defer fmt.Println(n.name + " Done")
	fmt.Println(n.name + " Starting network")

	// FBP distinguishes between execution of the process as a whole and activating the code - the code may be deactivated and then
	// reactivated many times during the process "run"

	n.wg.Add(len(n.procs))

	defer n.wg.Wait()
	var someProcsCanRun bool = false
	//time.Sleep(1 * time.Millisecond)
	for _, proc := range n.procs {

		selfStarting := true
		if proc.inPorts != nil {
			for _, conn := range proc.inPorts {
				//if conn.GetType() != "InitializationPort" {
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
	//}()
	//n.wg.Wait()
}
