package testrtn

import (
	//"fmt"
	"strconv"

	"github.com/jpaulm/gofbp"
)

type Sender struct {
	ipt gofbp.InputConn
	opt gofbp.OutputConn
}

func (sender *Sender) Setup(p *gofbp.Process) {
	sender.ipt = p.OpenInPort("COUNT")
	sender.opt = p.OpenOutPort("OUT")
}

func (sender *Sender) Execute(p *gofbp.Process) {

	icpkt := p.Receive(sender.ipt)
	j, _ := strconv.Atoi(icpkt.Contents.(string))
	p.Discard(icpkt)
	p.Close(sender.ipt)

	var pkt *gofbp.Packet
	for i := 0; i < j; i++ {
		pkt = p.Create("IP - # " + strconv.Itoa(i))
		p.Send(sender.opt, pkt)
	}

}
