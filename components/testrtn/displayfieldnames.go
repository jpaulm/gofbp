package testrtn

import (
	"fmt"

	"github.com/jpaulm/gofbp/core"
)

type DisplayFieldNames struct {
	ipt core.InputConn
	out core.OutputConn
}

func (dispFieldNames *DisplayFieldNames) Setup(p *core.Process) {
	dispFieldNames.ipt = p.OpenInPort("IN")
	dispFieldNames.out = p.OpenOutPortOptional("OUT")
}

func (DisplayFieldNames) MustRun() {}

func (dispFieldNames *DisplayFieldNames) Execute(p *core.Process) {

	for {
		var pkt = p.Receive(dispFieldNames.ipt)
		if pkt == nil {
			break
		}
		if pkt.PktType == core.OpenBracket {
			fmt.Println("Open Bracket", pkt.Contents)
		} else {
			if pkt.PktType == core.CloseBracket {
				fmt.Println("Close Bracket", pkt.Contents)
			} else {
				if pkt.PktType == core.Signal {
					fmt.Println("Signal", pkt.Contents)
				} else {
					//fmt.Println(pkt.Contents)
					fmt.Printf("%+v\n", pkt.Contents)
				}
			}
		}

		if dispFieldNames.out.IsConnected() {
			p.Send(dispFieldNames.out, pkt)
		} else {
			p.Discard(pkt)
		}
	}

}
