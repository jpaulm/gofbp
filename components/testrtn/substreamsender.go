package testrtn

import (
	"strconv"

	"github.com/jpaulm/gofbp/core"
)

type SubstreamSender struct {
	ipt core.InputConn
	opt core.OutputConn
}

func (sender *SubstreamSender) Setup(p *core.Process) {
	sender.ipt = p.OpenInPort("COUNT")
	sender.opt = p.OpenOutPort("OUT")
}

func (sender *SubstreamSender) Execute(p *core.Process) {

	icpkt := p.Receive(sender.ipt)
	j, _ := strconv.Atoi(icpkt.Contents.(string))
	p.Discard(icpkt)
	p.Close(sender.ipt)

	var pkt *core.Packet
	pkt = p.CreateBracket(core.OpenBracket, "")
	p.Send(sender.opt, pkt)

	for i := 0; i < j; i++ {
		k := i % 10

		if k == 2 {
			pkt = p.CreateBracket(core.CloseBracket, "")
			p.Send(sender.opt, pkt)
		}
		if k == 3 {
			pkt = p.CreateBracket(core.OpenBracket, "")
			p.Send(sender.opt, pkt)
		}

		if k == 7 || k == 0 {
			pkt = p.CreateBracket(core.CloseBracket, "")
			p.Send(sender.opt, pkt)
			pkt = p.CreateBracket(core.OpenBracket, "")
			p.Send(sender.opt, pkt)
		}

		pkt = p.Create("IP - # " + strconv.Itoa(i))
		p.Send(sender.opt, pkt)
	}
	pkt = p.CreateBracket(core.CloseBracket, "")
	p.Send(sender.opt, pkt)

}
