package testrtn

import (
	"fmt"
	"time"

	"github.com/jpaulm/gofbp"
)

type DelayedReceiver struct {
	ipt gofbp.InputConn
}

func (receiver *DelayedReceiver) Setup(p *gofbp.Process) {
	receiver.ipt = p.OpenInPort("IN")
}

func (DelayedReceiver) MustRun() {}

func (receiver *DelayedReceiver) Execute(p *gofbp.Process) {

	for {
		var pkt = p.Receive(receiver.ipt)
		if pkt == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
		fmt.Println("Input to DelayedReceiver:", p.Name, ">", pkt.Contents)
		p.Discard(pkt)
	}

}
