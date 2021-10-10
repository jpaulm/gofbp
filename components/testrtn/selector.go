package testrtn

import (
	"fmt"
	"strings"

	"github.com/jpaulm/gofbp/core"
)

type Selector struct {
	ipt   *core.InPort
	iptIp *core.InitializationConnection
	out1  *core.OutPort
	out2  *core.OutPort
}

func (selector *Selector) Setup(p *core.Process) {

	selector.ipt = p.OpenInPort("IN")

	selector.out1 = p.OpenOutPort("ACC")

	selector.out2 = p.OpenOutPort("REJ", "opt") // is optional

	selector.iptIp = p.OpenInitializationPort("PARAM")
}

func (Selector) MustRun() {}

func (selector *Selector) Execute(p *core.Process) {
	//fmt.Println(p.GetName() + " started")

	icpkt := p.Receive(selector.iptIp)
	param := icpkt.Contents.(string)
	i := len(param)

	p.Discard(icpkt)

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
				if selector.out2 == nil {
					p.Discard(pkt)
				}
			} else {
				if selector.out2 == nil {
					panic("Selector - port not specified, but not optional")
				}
			}

			p.Send(selector.out2, pkt)

		}
	}

	//fmt.Println(p.GetName() + " ended")
}
