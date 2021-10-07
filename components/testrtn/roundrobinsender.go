package testrtn

import (
	"fmt"

	"github.com/jpaulm/gofbp/core"
)

type RoundRobinSender struct {
	ipt core.InputConn
	out core.OutputArrayConn
}

func (rrsender *RoundRobinSender) Setup(p *core.Process) {
	rrsender.ipt = p.OpenInPort("IN")
	rrsender.out = p.OpenOutArrayPort("OUT")
}

func (rrsender *RoundRobinSender) Execute(p *core.Process) {
	//fmt.Println(p.GetName() + " started")

	var i = 0

	j := rrsender.out.ArrayLength()

	for {
		var pkt = p.Receive(rrsender.ipt)
		if pkt == nil {
			break
		}
		fmt.Println("Output: ", pkt.Contents)

		opt := rrsender.out.GetArrayItem(i)

		p.Send(opt, pkt)
		i = (i + 1) % j
	}

	//fmt.Println(p.GetName() + " ended")
}
