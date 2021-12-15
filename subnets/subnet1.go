package subnets

import (
	"github.com/tyoung3/gofbp"
	"github.com/tyoung3/gofbp/testrtn"
)

type Subnet1 struct{}

func (subnet *Subnet1) Setup(p *gofbp.Process) {}

func (subnet *Subnet1) Execute(p *gofbp.Process) {
	net := gofbp.NewSubnet("Subnet1", p)

	proc1 := net.NewProc("SubIn", &gofbp.SubIn{})

	proc2 := net.NewProc("WriteToConsole1", &testrtn.WriteToConsole{})

	proc3 := net.NewProc("SubOut", &gofbp.SubOut{})

	net.Initialize("IN", proc1, "NAME")
	net.Connect(proc1, "OUT", proc2, "IN", 6)
	net.Connect(proc2, "OUT", proc3, "IN", 6)
	net.Initialize("OUT", proc3, "NAME")

	net.Run()
}
