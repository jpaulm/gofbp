package testrtn

import (
	//"fmt"

	"github.com/jpaulm/gofbp/core"
)

type Discard struct {
	ipt core.InputConn
}

func (discard *Discard) Setup(p *core.Process) {
	discard.ipt = p.OpenInPort("IN")
}

//func (Discard) MustRun() {}

func (discard *Discard) Execute(p *core.Process) {
	

	for {
		var pkt = p.Receive(discard.ipt)
		if pkt == nil {
			break
		}
		//fmt.Println(pkt.Contents)

		p.Discard(pkt)

	}

	
}
