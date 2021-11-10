package core

import (
	"fmt"
	"sync"
)

// A subnet is like a combination of a Network and a Process...

type Subnet struct {
	Network
	//name   string
	//procs  map[string]*Process
	Mother *Process
	//wg     sync.WaitGroup
}

func (n *Subnet) SetMother(p *Process) {
	n.Mother = p
}

func (n *Subnet) GetMother() *Process {
	return n.Mother
}

//func (n *Subnet) GetNetwork() GenNet {
//	return n
//}

func (n *Subnet) GetProc(nm string) *Process {
	return n.procs[nm]
}

func (n *Subnet) SetProc(p *Process, nm string) {
	n.procs[nm] = p
}

func (n *Subnet) NewProc(nm string, comp Component) *Process {

	proc := &Process{
		name:      nm,
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
		n.SetMother(pStack[stkLevel-1])
	}

	return proc
}

func (n *Subnet) Run() {

	defer n.Exit()
	defer fmt.Println(n.name + " Subnet done")
	fmt.Println(n.name + " Starting subnet")

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
		panic("No process can start in subnet " + n.name)
	}
	//}()
	//n.wg.Wait()
}

func (n *Subnet) Exit() {
	if stkLevel > 0 {
		stkLevel--
	}
}
