package testrtn

import (
	"fmt"

	"github.com/jpaulm/gofbp"
)

type Receiver struct {
	ipt gofbp.InputConn
}

func (receiver *Receiver) Setup(p *gofbp.Process) {
	receiver.ipt = p.OpenInPort("IN")
}

func (Receiver) MustRun() {}

func (receiver *Receiver) Execute(p *gofbp.Process) {

	for {
		var pkt = p.Receive(receiver.ipt)
		if pkt == nil {
			break
		}

		fmt.Println("Input to Receiver:", p.Name, ">", pkt.Contents)
		p.Discard(pkt)
	}

}
