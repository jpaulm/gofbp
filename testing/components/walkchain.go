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
	chn := p.GetChain(pkt, "chain1")
	x := chn.First
	for x != nil {
		fmt.Println(x.Contents)
		x = x.Next
	}

	p.Send(wc.opt, pkt)
	p.Close(wc.ipt)
}
