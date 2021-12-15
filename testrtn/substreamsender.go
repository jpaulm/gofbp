package testrtn

import (
	"strconv"

	"github.com/tyoung3/gofbp"
)

type SubstreamSender struct {
	ipt gofbp.InputConn
	opt gofbp.OutputConn
}

func (sender *SubstreamSender) Setup(p *gofbp.Process) {
	sender.ipt = p.OpenInPort("COUNT")
	sender.opt = p.OpenOutPort("OUT")
}

func (sender *SubstreamSender) Execute(p *gofbp.Process) {

	icpkt := p.Receive(sender.ipt)
	j, _ := strconv.Atoi(icpkt.Contents.(string))
	p.Discard(icpkt)
	p.Close(sender.ipt)

	var pkt *gofbp.Packet
	pkt = p.CreateBracket(gofbp.OpenBracket, "")
	p.Send(sender.opt, pkt)

	for i := 0; i < j; i++ {
		k := i % 10

		if k == 2 {
			pkt = p.CreateBracket(gofbp.CloseBracket, "")
			p.Send(sender.opt, pkt)
		}
		if k == 3 {
			pkt = p.CreateBracket(gofbp.OpenBracket, "")
			p.Send(sender.opt, pkt)
		}

		if k == 7 || k == 0 {
			pkt = p.CreateBracket(gofbp.CloseBracket, "")
			p.Send(sender.opt, pkt)
			pkt = p.CreateBracket(gofbp.OpenBracket, "")
			p.Send(sender.opt, pkt)
		}

		pkt = p.Create("IP - # " + strconv.Itoa(i))
		p.Send(sender.opt, pkt)
	}
	pkt = p.CreateBracket(gofbp.CloseBracket, "")
	p.Send(sender.opt, pkt)

}
