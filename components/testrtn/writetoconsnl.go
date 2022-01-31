package testrtn

import (
	"fmt"

	"github.com/jpaulm/gofbp/core"
)

// WriteToConsole modified to be a non-looper (NL)

type WriteToConsNL struct {
	ipt core.InputConn
	opt core.OutputConn
}

func (writeToConsole *WriteToConsNL) Setup(p *core.Process) {
	writeToConsole.ipt = p.OpenInPort("IN")
	writeToConsole.opt = p.OpenOutPortOptional("OUT")
}

//func (WriteToConsNL) MustRun() {}

func (writeToConsole *WriteToConsNL) Execute(p *core.Process) {

	//for {
	var pkt = p.Receive(writeToConsole.ipt)
	if pkt == nil {
		//break
		return
	}
	if pkt.PktType == core.OpenBracket {
		fmt.Println("Open Bracket", pkt.Contents)
	} else {
		if pkt.PktType == core.CloseBracket {
			fmt.Println("Close Bracket", pkt.Contents)
		} else {
			fmt.Print(pkt.Contents)
		}
	}

	if writeToConsole.opt.IsConnected() {
		p.Send(writeToConsole.opt, pkt)
	} else {
		p.Discard(pkt)
	}
	//}

}
