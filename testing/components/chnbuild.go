package components

import (
	"github.com/jpaulm/gofbp/core"
)

type ChainBuild struct {
	opt core.OutputConn
}

func (cb *ChainBuild) Setup(p *core.Process) {
	cb.opt = p.OpenOutPort("OUT")
}

func (cb *ChainBuild) Execute(p *core.Process) {

	pkt := p.Create("One")
	chn := p.NewChain(pkt, "chain1")
	pkt2 := p.Create("Two")
	pkt3 := p.Create("Three")
	pkt4 := p.Create("Four")
	p.Attach(chn, pkt2)
	p.Attach(chn, pkt3)
	p.Attach(chn, pkt4)

	p.Send(cb.opt, pkt)
}
