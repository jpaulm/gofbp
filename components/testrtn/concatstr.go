package testrtn

import (
	"github.com/jpaulm/gofbp/core"
)

type ConcatStr struct {
	ipt core.InputArrayConn
	opt core.OutputConn
}

func (concatstr *ConcatStr) Setup(p *core.Process) {
	concatstr.ipt = p.OpenInArrayPort("IN")
	concatstr.opt = p.OpenOutPort("OUT")
}

func (concatstr *ConcatStr) Execute(p *core.Process) {

	for _, inPort := range concatstr.ipt.GetArray() {
		for {
			if inPort == nil {
				continue
			}
			var pkt = p.Receive(inPort)
			if pkt == nil {
				break
			}
			p.Send(concatstr.opt, pkt)
		}
	}
}
