package components

import (
	"github.com/jpaulm/gofbp/core"
)

type WSAnsReq struct {
	ipt core.InputConn
	out core.OutputConn
}

func (wsansreq *WSAnsReq) Setup(p *core.Process) {
	wsansreq.ipt = p.OpenInPort("IN")
	wsansreq.out = p.OpenOutPortOptional("OUT")
}

func (WSAnsReq) MustRun() {}

func (wsansreq *WSAnsReq) Execute(p *core.Process) {
	pkt := p.Receive(wsansreq.ipt)
	for {

		if pkt == nil {
			break
		}
		p.Send(wsansreq.out, pkt)

		pkt = p.Receive(wsansreq.ipt) // connection
		p.Send(wsansreq.out, pkt)

		pkt = p.Receive(wsansreq.ipt) //"namelist"
		p.Discard(pkt)

		pkt = p.Create("line1")
		p.Send(wsansreq.out, pkt)
		pkt = p.Create("line2")
		p.Send(wsansreq.out, pkt)
		pkt = p.Create("line3")
		p.Send(wsansreq.out, pkt)

		pkt = p.Receive(wsansreq.ipt) // close bracket
		p.Send(wsansreq.out, pkt)

		pkt = p.Receive(wsansreq.ipt)
		if pkt.Contents.(string) == "@kill" {
			p.Send(wsansreq.out, pkt)
			pkt = p.Receive(wsansreq.ipt)
		}
	}

}
