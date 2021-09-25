package testrtn

import (
	"fmt"
	"github.com/jpaulm/gofbp/core"
	"strings"
)

type Selector struct {
	ipt   core.InputConn
	iptIp core.InputConn
	opt1  core.OutputConn
	opt2  core.OutputConn
}

func (selector *Selector) Setup(p *core.Process) {
	selector.ipt = p.OpenInPort("IN")
	selector.opt1 = p.OpenOutPort("ACC")
	selector.opt2 = p.OpenOutPort("REJ")
	selector.iptIp = p.OpenInPort("PARAM")
}

func (Selector) MustRun() {}

func (selector *Selector) Execute(p *core.Process) {
	fmt.Println(p.GetName() + " started")

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

		if 0 == strings.Compare(s[:i], param[:i]) {
			p.Send(selector.opt1.(*core.OutPort), pkt)
		} else {
			p.Send(selector.opt2.(*core.OutPort), pkt)
		}

	}

	fmt.Println(p.GetName() + " ended")
}
