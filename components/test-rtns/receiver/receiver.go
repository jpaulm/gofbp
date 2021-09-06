package receiver

import (
	"fmt"

	"github.com/jpaulm/gofbp/core"
)

type Component struct {
	ipt *core.InPort
}

func New() *Component {
	return &Component{}
}

func (comp *Component) OpenPorts(p *core.Process) {
	comp.ipt = p.OpenInPort("IN")
}

func (comp *Component) Execute(p *core.Process) {
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
