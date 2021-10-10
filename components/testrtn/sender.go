package testrtn

import (
	//"fmt"
	"strconv"

	"github.com/jpaulm/gofbp/core"
)

type Sender struct {
	ipt *core.InitializationConnection
	opt *core.OutPort
}

func (sender *Sender) Setup(p *core.Process) {
	sender.ipt = p.OpenInitializationPort("COUNT")
	sender.opt = p.OpenOutPort("OUT")
}

func (sender *Sender) Execute(p *core.Process) {
	//fmt.Println(p.GetName() + " started")
	icpkt := p.Receive(sender.ipt)
	j, _ := strconv.Atoi(icpkt.Contents.(string))
	p.Discard(icpkt)
	var pkt *core.Packet
	for i := 0; i < j; i++ {
		pkt = p.Create("IP - # " + strconv.Itoa(i))
		p.Send(sender.opt, pkt)
	}
	//fmt.Println(p.GetName() + " ended")
}
