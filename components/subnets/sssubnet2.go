package subnets

import (
	"github.com/jpaulm/gofbp/components/testrtn"
	"github.com/jpaulm/gofbp/core"
)

type SSSubnet2 struct{}

func (subnet *SSSubnet2) Setup(p *core.Process) {}

func (subnet *SSSubnet2) Execute(p *core.Process) {
	net := core.NewSubnet("SSSubnet2", p)

	proc1 := net.NewProc("SubInSS", &core.SubInSS{}) // Substream-Sensitive SubIn

	proc2 := net.NewProc("Count", &testrtn.Counter{}) // count length of each substream (excl. brackets)

	proc3 := net.NewProc("SubOut", &core.SubOut{}) // Basic SubOut

	net.Initialize("IN", proc1, "NAME")
	net.Connect(proc1, "OUT", proc2, "IN", 6)
	net.Connect(proc2, "COUNT", proc3, "IN", 6)
	net.Initialize("OUT", proc3, "NAME")

	net.Run()
}
