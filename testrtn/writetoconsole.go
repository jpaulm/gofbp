package testrtn

import (
	"fmt"

	"github.com/tyoung3/gofbp"
)

type WriteToConsole struct {
	ipt gofbp.InputConn
	out gofbp.OutputConn
}

func (writeToConsole *WriteToConsole) Setup(p *gofbp.Process) {
	writeToConsole.ipt = p.OpenInPort("IN")
	writeToConsole.out = p.OpenOutPortOptional("OUT")
}

func (WriteToConsole) MustRun() {}

func (writeToConsole *WriteToConsole) Execute(p *gofbp.Process) {

	for {
		var pkt = p.Receive(writeToConsole.ipt)
		if pkt == nil {
			break
		}
		if pkt.PktType == gofbp.OpenBracket {
			fmt.Println("Open", pkt.Contents)
		} else {
			if pkt.PktType == gofbp.CloseBracket {
				fmt.Println("Close", pkt.Contents)
			} else {
				fmt.Println(pkt.Contents)
			}
		}

		if writeToConsole.out.IsConnected() {
			p.Send(writeToConsole.out, pkt)
		} else {
			p.Discard(pkt)
		}
	}

}
