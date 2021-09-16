package testrtn

import (
	"fmt"

	"github.com/jpaulm/gofbp/core"
)

type Receiver struct {
	conn core.Conn
}

func (receiver *Receiver) OpenPorts(p *core.Process) {
	receiver.conn = p.OpenInPort("IN")
}

func (receiver *Receiver) Execute(p *core.Process) {
	fmt.Println(p.Name + " started")

	for {
		var pkt = p.Receive(receiver.conn)
		if pkt == nil {
			break
		}
		fmt.Println("Output: ", pkt.Contents)
		p.Discard(pkt)
	}

	fmt.Println(p.Name + " ended")
}

func (receiver *Receiver) GetMustRun() bool {
	return true
}
