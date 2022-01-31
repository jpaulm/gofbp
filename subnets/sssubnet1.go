package subnets

import (
	"github.com/jpaulm/gofbp"
	"github.com/jpaulm/gofbp/testrtn"
)

type SSSubnet1 struct{}

func (subnet *SSSubnet1) Setup(p *gofbp.Process) {}

func (subnet *SSSubnet1) Execute(p *gofbp.Process) {
	net := gofbp.NewSubnet("SSSubnet1", p)

	proc1 := net.NewProc("SubInSS", &gofbp.SubInSS{}) // Substream-Sensitive SubIn

	proc2 := net.NewProc("WriteToConsole1", &testrtn.WriteToConsole{})

	proc3 := net.NewProc("SubOutSS", &gofbp.SubOutSS{}) // Substream-delimiter-Generating SubOut

	net.Initialize("IN", proc1, "NAME")
	net.Connect(proc1, "OUT", proc2, "IN", 6)
	net.Connect(proc2, "OUT", proc3, "IN", 6)
	net.Initialize("OUT", proc3, "NAME")

	net.Run()
}
