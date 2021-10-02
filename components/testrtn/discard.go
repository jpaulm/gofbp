package testrtn

import (
	//"fmt"

	"github.com/jpaulm/gofbp/core"
)

type Discard struct {
	ipt core.InputConn
	opt core.OutputConn
}

func (discard *Discard) Setup(p *core.Process) {
	discard.ipt = p.OpenInPort("IN")
	discard.opt = p.OpenOutPort("OUT", "opt")
}

//func (Discard) MustRun() {}

func (discard *Discard) Execute(p *core.Process) {
	//fmt.Println(p.GetName() + "	started")

	for {
		var pkt = p.Receive(discard.ipt)
		if pkt == nil {
			break
		}
		//fmt.Println(pkt.Contents)

		p.Discard(pkt)

	}

	//fmt.Println(p.GetName() + " ended")
}
