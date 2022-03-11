package components

import (
	"fmt"

	"github.com/jpaulm/gofbp/core"
)

type WalkChain struct {
	ipt core.InputConn
	opt core.OutputConn
}

func (wc *WalkChain) Setup(p *core.Process) {
	wc.ipt = p.OpenInPort("IN")
	wc.opt = p.OpenOutPort("OUT")
}

func (wc *WalkChain) Execute(p *core.Process) {

	var pkt = p.Receive(wc.ipt)
	chn, ok := p.GetChain(pkt, "chain1")
	if !ok {
		panic("Chain 'chain1' not found")
	}
	x := chn.First
	for x != nil {
		fmt.Print(x.Contents, "\n")
		x = x.Next
	}

	p.Send(wc.opt, pkt)
	p.Close(wc.ipt)
}
