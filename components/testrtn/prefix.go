package testrtn

import (
	"fmt"

	"github.com/jpaulm/gofbp/core"
)

type Prefix struct {
	ipt   core.InputConn
	iptIP core.InputConn
	out   core.OutputConn
}

func (prefix *Prefix) Setup(p *core.Process) {

	prefix.ipt = p.OpenInPort("IN")

	prefix.out = p.OpenOutPort("OUT")

	prefix.iptIP = p.OpenInPort("PARAM")
}

func (prefix *Prefix) Execute(p *core.Process) {

	icpkt := p.Receive(prefix.iptIP)
	param := icpkt.Contents.(string)

	p.Discard(icpkt)

	p.Close(prefix.iptIP)

	for {
		var pkt = p.Receive(prefix.ipt)
		if pkt == nil {
			break
		}
		fmt.Print(pkt.Contents, "\n")

		s, ok := pkt.Contents.(string)
		if !ok {
			panic("IP contents must be a string")
		}
		s = param + s
		pkt.Contents = s
		p.Send(prefix.out, pkt)

	}
}
