package testrtn

import (
	"fmt"

	"github.com/jpaulm/gofbp/core"
)

type Receiver struct {
	ipt core.InputConn
}

func (receiver *Receiver) Setup(p *core.Process) {
	receiver.ipt = p.OpenInPort("IN")
}

func (Receiver) MustRun() {}

func (receiver *Receiver) Execute(p *core.Process) {

	for {
		var pkt = p.Receive(receiver.ipt)
		if pkt == nil {
			break
		}

		fmt.Print("Input to Receiver:", p.Name, ">", pkt.Contents, "\n")
		p.Discard(pkt)
	}

}
