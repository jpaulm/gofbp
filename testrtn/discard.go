package testrtn

import (
	//"fmt"

	"github.com/jpaulm/gofbp"
)

type Discard struct {
	ipt gofbp.InputConn
}

func (discard *Discard) Setup(p *gofbp.Process) {
	discard.ipt = p.OpenInPort("IN")
}

//func (Discard) MustRun() {}

func (discard *Discard) Execute(p *gofbp.Process) {

	for {
		var pkt = p.Receive(discard.ipt)
		if pkt == nil {
			break
		}
		//fmt.Println(pkt.Contents)

		p.Discard(pkt)

	}

}
