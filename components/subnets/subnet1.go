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

	proc2 := net.NewProc("WriteToConsole1", &testrtn.WriteToConsole{})

	proc3 := net.NewProc("SubOut", &core.SubOut{})

	net.Initialize("IN", proc1, "NAME")
	net.Connect(proc1, "OUT", proc2, "IN", 6)
	net.Connect(proc2, "OUT", proc3, "IN", 6)
	net.Initialize("OUT", proc3, "NAME")

	net.Run()
}
