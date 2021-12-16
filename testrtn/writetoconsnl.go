package testrtn

import (
	"fmt"

	"github.com/jpaulm/gofbp"
)

// WriteToConsole modified to be a non-looper (NL)

type WriteToConsNL struct {
	ipt gofbp.InputConn
	opt gofbp.OutputConn
}

func (writeToConsole *WriteToConsNL) Setup(p *gofbp.Process) {
	writeToConsole.ipt = p.OpenInPort("IN")
	writeToConsole.opt = p.OpenOutPortOptional("OUT")
}

//func (WriteToConsNL) MustRun() {}

func (writeToConsole *WriteToConsNL) Execute(p *gofbp.Process) {

	//for {
	var pkt = p.Receive(writeToConsole.ipt)
	if pkt == nil {
		//break
		return
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

	if writeToConsole.opt.IsConnected() {
		p.Send(writeToConsole.opt, pkt)
	} else {
		p.Discard(pkt)
	}
	//}

}
