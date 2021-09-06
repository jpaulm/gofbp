package sender

import (
	"fmt"
	"strconv"

	"github.com/jpaulm/gofbp/core"
)

type Component struct {
	out *core.OutPort
}

func New() *Component {
	return &Component{}
}

func (comp *Component) OpenPorts(p *core.Process) {
	comp.out = p.OpenOutPort("OUT")
}

func (comp *Component) Execute(p *core.Process) {
	fmt.Println(p.Name + " started")
	var pkt *core.Packet
	for i := 0; i < 15; i++ {
		pkt = p.Create("IP - # " + strconv.Itoa(i) + " (" + p.Name + ")")
		p.Send(comp.out.Conn, pkt)
	}
	fmt.Println(p.Name + " ended")
}
