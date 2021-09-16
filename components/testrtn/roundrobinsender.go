package testrtn

import (
	"fmt"

	"github.com/jpaulm/gofbp/core"
)

type RoundRobinSender struct {
	ipt core.InputConn
	opt core.OutputConn
}

func (rrsender *RoundRobinSender) OpenPorts(p *core.Process) {
	rrsender.ipt = p.OpenInPort("IN")
}

func (rrsender *RoundRobinSender) Execute(p *core.Process) {
	fmt.Println(p.Name + " started")

	for {
		var pkt = p.Receive(rrsender.ipt)
		if pkt == nil {
			break
		}
		fmt.Println("Output: ", pkt.Contents)
		p.Discard(pkt)
	}

	fmt.Println(p.Name + " ended")
}

func (rrsender *RoundRobinSender) GetMustRun() bool {
	return true
}
