package testrtn

import (
	"fmt"

	"github.com/jpaulm/gofbp/core"
)

// WriteToConsole modifieed to be a non-looper (NL)

type WriteToConsNL struct {
	ipt core.InputConn
	opt core.OutputConn
}

func (writeToConsole *WriteToConsNL) Setup(p *core.Process) {
	writeToConsole.ipt = p.OpenInPort("IN")
	writeToConsole.opt = p.OpenOutPort("OUT", "opt")
}

//func (WriteToConsNL) MustRun() {}

func (writeToConsole *WriteToConsNL) Execute(p *core.Process) {
	//fmt.Println(p.GetName() + " activated")

	//for {
	var pkt = p.Receive(writeToConsole.ipt)
	if pkt == nil {
		//break
		return
	}
	fmt.Println(pkt.Contents)
	//if writeToConsole.opt.GetType() == "OutPort" {
	_, b := writeToConsole.opt.(*core.OutPort)
	if b {
		p.Send(writeToConsole.opt, pkt)
	} else {
		p.Discard(pkt)
	}
	//}

	//fmt.Println(p.GetName() + " deactivated")
}
