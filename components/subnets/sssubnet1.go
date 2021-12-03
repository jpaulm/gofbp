/*Package subnets implements FBP subnet processing*/
package subnets

import (
	"github.com/jpaulm/gofbp/components/testrtn"
	"github.com/jpaulm/gofbp/core"
)

type SSSubnet1 struct{}

func (subnet *SSSubnet1) Setup(p *core.Process) {}

func (subnet *SSSubnet1) Execute(p *core.Process) {
	net := core.NewSubnet("SSSubnet1", p)

	proc1 := net.NewProc("SubInSS", &core.SubInSS{}) // Substream-Sensitive SubIn

	proc2 := net.NewProc("WriteToConsole1", &testrtn.WriteToConsole{})

	proc3 := net.NewProc("SubOutSS", &core.SubOutSS{}) // Substream-delimiter-Generating SubOut

	net.Initialize("IN", proc1, "NAME")
	net.Connect(proc1, "OUT", proc2, "IN", 6)
	net.Connect(proc2, "OUT", proc3, "IN", 6)
	net.Initialize("OUT", proc3, "NAME")

	net.Run()
}
