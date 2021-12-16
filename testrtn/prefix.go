package testrtn

import (
	"fmt"

	"github.com/jpaulm/gofbp"
)

type Prefix struct {
	ipt   gofbp.InputConn
	iptIP gofbp.InputConn
	out   gofbp.OutputConn
}

func (prefix *Prefix) Setup(p *gofbp.Process) {

	prefix.ipt = p.OpenInPort("IN")

	prefix.out = p.OpenOutPort("OUT")

	prefix.iptIP = p.OpenInPort("PARAM")
}

func (prefix *Prefix) Execute(p *gofbp.Process) {

	icpkt := p.Receive(prefix.iptIP)
	param := icpkt.Contents.(string)

	p.Discard(icpkt)

	p.Close(prefix.iptIP)

	for {
		var pkt = p.Receive(prefix.ipt)
		if pkt == nil {
			break
		}
		fmt.Println(pkt.Contents)

		s, ok := pkt.Contents.(string)
		if !ok {
			panic("IP contents must be a string")
		}
		s = param + s
		pkt.Contents = s
		p.Send(prefix.out, pkt)

	}
}
