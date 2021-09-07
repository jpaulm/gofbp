package testrtn

import (
	"fmt"

	"github.com/jpaulm/gofbp/core"
)

type Receiver struct {
	ipt *core.InPort
}

func (comp *Receiver) OpenPorts(p *core.Process) {
	comp.ipt = p.OpenInPort("IN")
}

func (comp *Receiver) Execute(p *core.Process) {
	fmt.Println(p.Name + " started")

	for {
		var pkt = p.Receive(comp.ipt.Conn)
		if pkt == nil {
			break
		}
		fmt.Println("Output: ", pkt.Contents)
		p.Discard(pkt)
	}

	fmt.Println(p.Name + " ended")
}
