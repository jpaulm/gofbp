package testrtn

import (
	"fmt"
	"strings"

	"github.com/tyoung3/gofbp"
)

type Selector struct {
	ipt   gofbp.InputConn
	iptIP gofbp.InputConn
	out1  gofbp.OutputConn
	out2  gofbp.OutputConn
}

func (selector *Selector) Setup(p *gofbp.Process) {

	selector.ipt = p.OpenInPort("IN")

	selector.out1 = p.OpenOutPort("ACC")

	selector.out2 = p.OpenOutPortOptional("REJ")

	selector.iptIP = p.OpenInPort("PARAM")
}

func (Selector) MustRun() {}

func (selector *Selector) Execute(p *gofbp.Process) {

	icpkt := p.Receive(selector.iptIP)
	param := icpkt.Contents.(string)
	i := len(param)

	p.Discard(icpkt)

	p.Close(selector.iptIP)

	for {
		var pkt = p.Receive(selector.ipt)
		if pkt == nil {
			break
		}
		fmt.Println(pkt.Contents)

		s := pkt.Contents.(string)
		if i > len(s) {
			i = len(s)
		}

		if 0 == strings.Compare(param[:i], s[:i]) {
			p.Send(selector.out1, pkt)
		} else {
			if !selector.out2.IsConnected() {
				p.Discard(pkt)
			} else {
				p.Send(selector.out2, pkt)
			}
		}
	}
}
