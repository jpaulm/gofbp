package testrtn

import (
	"fmt"
	"strings"

	"github.com/jpaulm/gofbp/core"
)

type Selector struct {
	ipt   core.InputConn
	iptIp core.InputConn
	out1  core.OutputConn
	out2  core.OutputConn
}

func (selector *Selector) Setup(p *core.Process) {

	selector.ipt = p.OpenInPort("IN")

	selector.out1 = p.OpenOutPort("ACC")

	selector.out2 = p.OpenOutPortOptional("REJ")

	selector.iptIp = p.OpenInPort("PARAM")
}

func (Selector) MustRun() {}

func (selector *Selector) Execute(p *core.Process) {

	icpkt := p.Receive(selector.iptIp)
	param := icpkt.Contents.(string)
	i := len(param)

	p.Discard(icpkt)

	p.Close(selector.iptIp)

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
