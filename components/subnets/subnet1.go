package subnets

import (
	"github.com/jpaulm/gofbp/components/testrtn"
	"github.com/jpaulm/gofbp/core"
)

type Subnet1 struct{}

func (subnet *Subnet1) Setup(p *core.Process) {}

func (subnet *Subnet1) Execute(p *core.Process) {
	net := core.NewSubnet("Subnet1", p)

	proc1 := net.NewProc("SubIn", &core.SubIn{})

	proc1a := net.NewProc("Prefix", &testrtn.Prefix{})

	proc2 := net.NewProc("WriteToConsole", &testrtn.WriteToConsole{})

	proc3 := net.NewProc("SubOut", &core.SubOut{})

	net.Initialize("IN", proc1, "NAME")
	net.Connect(proc1, "OUT", proc1a, "IN", 6)
	net.Initialize("X-", proc1a, "PARAM")
	net.Connect(proc1a, "OUT", proc2, "IN", 6)
	net.Connect(proc2, "OUT", proc3, "IN", 6)
	net.Initialize("OUT", proc3, "NAME")

	net.Run()
}

//  NAME->0(proc1 SubIn)1.OUT -> IN.0(Proc2 WriteToConsole1)1.OUT -> NAME.0(proc3 SubOut);
