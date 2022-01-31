package testrtn

import (
	"github.com/jpaulm/gofbp"
)

type RoundRobinSender struct {
	ipt gofbp.InputConn
	out gofbp.OutputArrayConn
}

func (rrsender *RoundRobinSender) Setup(p *gofbp.Process) {
	rrsender.ipt = p.OpenInPort("IN")
	rrsender.out = p.OpenOutArrayPort("OUT")
}

func (rrsender *RoundRobinSender) Execute(p *gofbp.Process) {

	var i = 0

	j := rrsender.out.ArrayLength()

	for {
		var pkt = p.Receive(rrsender.ipt)
		if pkt == nil {
			break
		}
		//fmt.Println("Output: ", pkt.Contents)

		opt := rrsender.out.GetArrayItem(i)

		p.Send(opt, pkt)
		i = (i + 1) % j
	}

}
