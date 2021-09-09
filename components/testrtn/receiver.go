package testrtn

import (
	"fmt"

	"github.com/jpaulm/gofbp/core"
)

type Receiver struct {
	conn core.Conn
}

func (comp *Receiver) OpenPorts(p *core.Process) {
	comp.conn = p.OpenInPort("IN").(*Connection)
}

func (comp *Receiver) Execute(p *core.Process) {
	fmt.Println(p.Name + " started")

	for {
		var pkt = p.Receive(comp.conn)
		if pkt == nil {
			break
		}
		fmt.Println("Output: ", pkt.Contents)
		p.Discard(pkt)
	}

	fmt.Println(p.Name + " ended")
}
