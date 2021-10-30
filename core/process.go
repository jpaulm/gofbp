package core

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
)

//import (
//	"github.com/gofbp/core"
//)

type Process struct {
	name    string
	network GenNet

	inPorts  map[string]inputCommon
	outPorts map[string]outputCommon

	logFile   string
	component Component
	ownedPkts int
	//done         bool
	selfStarting bool // process has no non-IIP input ports
	//MustRun   bool
	status int32
	mtx    sync.Mutex
	//Mother *Process
}

const (
	Notstarted int32 = iota
	Active
	Dormant
	SuspSend
	SuspRecv
	Terminated
)

func (p *Process) GetName() string {
	return p.name
}

func (p *Process) OpenInPort(s string) InputConn {
	var in InputConn
	var b bool
	if len(p.inPorts) == 0 {
		panic(p.name + ": No input ports specified")
	}
	in, b = p.inPorts[s].(InputConn)
	//if in == nil {
	//	panic(p.name + ": Port name not found (" + s + ")")
	//}
	if !b {
		panic(p.name + " " + s + " InPort not connected, or found other type")
	}
	return in
}

func (p *Process) OpenInArrayPort(s string) *InArrayPort {
	var in *InArrayPort
	var b bool
	if len(p.inPorts) == 0 {
		panic(p.name + ": No input ports specified")
	}
	in, b = p.inPorts[s].(*InArrayPort)
	//if in == nil {
	//	panic(p.name + ": Port name not found (" + s + ")")
	//}
	if !b {
		panic(p.name + " " + s + " InArrayPort not connected, or found other type")
	}
	return in
}

func (p *Process) OpenOutPort(s string) OutputConn {
	var out OutputConn
	var b bool
	if len(p.outPorts) == 0 {
		panic(p.name + " " + s + " OutPort not connected")
	} else {
		out, b = p.outPorts[s].(*OutPort)
		if !b {
			panic(p.name + " " + s + " OutPort not connected, or found other type")
		}
	}

	return out

}

func (p *Process) OpenOutPortOptional(s string) OutputConn {
	var out OutputConn
	var b bool
	if len(p.outPorts) == 0 {
		out = new(NullOutPort)
		p.outPorts[s] = out
		//out.name = s
		//out.connected = false
	} else {
		out, b = p.outPorts[s].(*OutPort)
		//if !b {
		//	panic(p.name + " " + s + " OutPort not connected, or found other type")
		//}

		if !b {
			out := new(NullOutPort)
			p.outPorts[s] = out
			//out.name = s
			//out.connected = false
		}
	}

	return out
}

// not sure it makes sense to allow optional for array ports!

func (p *Process) OpenOutArrayPort(s string) *OutArrayPort {
	var out *OutArrayPort
	var b bool
	if len(p.outPorts) == 0 {
		out = new(OutArrayPort)
		p.outPorts[s] = out
		out.name = s
		out.connected = false
	} else {
		out, b = p.outPorts[s].(*OutArrayPort)
		if !b {
			panic(p.name + " " + s + " OutArrayPort not connected, or found other type")
		}
	}

	return out
}

// Send sends a packet to the output connection.
// Returns false when fails to send.
func (p *Process) Send(o OutputConn, pkt *Packet) bool {
	return o.send(p, pkt)
}

// Receive receives from the connection.
// Returns nil, when there's no more data.
func (p *Process) Receive(c InputConn) *Packet {
	return c.receive(p)
}

func (p *Process) ensureRunning() {

	if !atomic.CompareAndSwapInt32(&p.status, Notstarted, Active) {
		return
	}

	netx, b := p.network.(*Network)

	var wg *sync.WaitGroup
	if b {
		//s = netx.id()
		wg = &netx.wg
	} else {
		nets, _ := p.network.(*Subnet)
		wg = &nets.wg
	}

	//pprof.Do(context.TODO(), pprof.Labels(
	//	"network", p.network.(*Network).id(),
	//	"process", p.name,
	//), func(c context.Context) {
	go func() { // Process goroutine
		defer wg.Done()
		p.Run()
	}()
	//})
}

func (p *Process) inputState() (bool, bool) {

	allDrained := true
	hasData := false
	for _, v := range p.inPorts {
		//if v.GetType() == "InArrayPort" {
		_, b := v.(*InArrayPort)
		if b {
			//allClosed = true
			for _, w := range v.(*InArrayPort).array {
				if !w.isDrained() /* || !w.IsClosed() */ {
					allDrained = false
				}
				hasData = hasData || !w.IsEmpty()
			}
		} else {
			_, b := v.(*InPort)
			if b {
				if !v.(*InPort).isDrained() /*|| !v.IsClosed() */ {
					allDrained = false
				}
				hasData = hasData || !v.(*InPort).IsEmpty()
			}
		}
	}
	return allDrained, hasData
}

func (p *Process) Run() {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	//atomic.StoreInt32(&p.status, Dormant)
	defer atomic.StoreInt32(&p.status, Terminated)
	defer trace(p.GetName(), " terminated")
	trace(p.GetName(), " started")

	fmt.Println("Goroutine", p.GetName()+":", "no.", getGID())

	p.component.Setup(p)

	//var allDrained bool
	//var hasData bool

	//atomic.StoreInt32(&p.status, Dormant)

	allDrained, hasData := p.inputState()

	canRun := p.selfStarting || hasData || !allDrained || p.isMustRun() && allDrained

	for canRun {
		// multiple activations, if necessary!
		trace(p.GetName(), " activated")
		atomic.StoreInt32(&p.status, Active)
		//atomic.StoreInt32(&p.network.Active, 1)

		p.component.Execute(p) // single "activation"

		atomic.StoreInt32(&p.status, Dormant)
		trace(p.GetName(), " deactivated")

		if p.ownedPkts > 0 {
			panic(p.name + " deactivated without disposing of all owned packets")
		}

		allDrained, _ := p.inputState()

		if allDrained {
			canRun = false
		} else {
			for _, v := range p.inPorts {
				//_, b := v.(InitializationConnection)
				//if b {
				//v.(*InitializationConnection).resetForNextExecution()
				v.resetForNextExecution()
				//}
			}
		}

	}

	for _, v := range p.outPorts {

		if v.IsConnected() {
			v.Close()
		}
		// if anything else, just continue
	}
}

func (p *Process) isMustRun() bool {
	_, hasMustRun := p.component.(ComponentWithMustRun)
	return hasMustRun
}

// create packet containing anything!
func (p *Process) Create(x interface{}) *Packet {
	var pkt *Packet = new(Packet)
	pkt.Contents = x
	pkt.owner = p
	p.ownedPkts++
	return pkt
}

// create bracket
func (p *Process) CreateBracket(pktType int32, s string) *Packet {
	var pkt *Packet = new(Packet)
	pkt.Contents = s
	pkt.pktType = pktType
	pkt.owner = p
	p.ownedPkts++
	return pkt
}

func (p *Process) Discard(pkt *Packet) {
	if pkt == nil {
		panic("Discarding nil packet")
	}
	p.ownedPkts--
	pkt = nil
}

//https://blog.sgmansfield.com/2015/12/goroutine-ids/

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
