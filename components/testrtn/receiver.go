package testrtn

import (
	"fmt"

	"github.com/jpaulm/gofbp/core"
)

type Receiver struct {
	ipt *core.InPort
}

func (receiver *Receiver) Setup(p *core.Process) {
	receiver.ipt = p.OpenInPort("IN")
}

func (Receiver) MustRun() {}

func (receiver *Receiver) Execute(p *core.Process) {
	//fmt.Println(p.GetName() + " started")

	for {
		var pkt = p.Receive(receiver.ipt)
		if pkt == nil {
			break
		}
		fmt.Println("Output: ", pkt.Contents)
		p.Discard(pkt)
	}

	//fmt.Println(p.GetName() + " ended")
}
