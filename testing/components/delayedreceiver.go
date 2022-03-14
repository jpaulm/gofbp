package components

import (
	"math/rand"
	"time"

	"github.com/jpaulm/gofbp/core"
)

type DelayedReceiver struct {
	ipt core.InputConn
}

func (receiver *DelayedReceiver) Setup(p *core.Process) {
	receiver.ipt = p.OpenInPort("IN")
}

func (DelayedReceiver) MustRun() {}

func (receiver *DelayedReceiver) Execute(p *core.Process) {

	for {
		var pkt = p.Receive(receiver.ipt)
		if pkt == nil {
			break
		}
		time.Sleep(time.Duration(rand.Int31n(500)) * time.Millisecond)
		//fmt.Print("DelayedReceiver input:", p.Name, ">", pkt.Contents,  "\n")
		p.Discard(pkt)
	}

}
