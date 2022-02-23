package components

import (
	"fmt"

	"github.com/jpaulm/gofbp/core"
)

type ChainBrkUp struct {
	ipt core.InputConn
	opt core.OutputConn
}

func (cb *ChainBrkUp) Setup(p *core.Process) {
	cb.ipt = p.OpenInPort("IN")
	cb.opt = p.OpenOutPort("OUT")
}

func (cb *ChainBrkUp) Execute(p *core.Process) {

	var pkt = p.Receive(cb.ipt)
	chn, ok := p.GetChain(pkt, "chain1")
	if !ok {
		panic("Chain 'chain1' not found")
	}
	x := chn.First
	for x != nil {
		fmt.Println(x.Contents)
		p.Detach(chn, x)
		p.Send(cb.opt, x)
		x = chn.First
	}
	fmt.Println(pkt.Contents)
	p.Send(cb.opt, pkt)
}
