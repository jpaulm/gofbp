/*Package testrtn tests gofbp code.*/
package testrtn

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
	pkt2 := p.Create("Two")
	pkt3 := p.Create("Three")
	pkt4 := p.Create("Four")
	p.Attach(pkt, "chain", pkt2)
	p.Attach(pkt2, "chain", pkt3)
	p.Attach(pkt3, "chain", pkt4)

	p.Send(cb.opt, pkt)
}
