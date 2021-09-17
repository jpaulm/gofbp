package testrtn

import (
	"fmt"
	"strconv"

	"github.com/jpaulm/gofbp/core"
)

type Sender struct {
	ipt core.InputConn
	opt core.OutputConn
}

func (sender *Sender) OpenPorts(p *core.Process) {
	sender.ipt = p.OpenInPort("COUNT")
	sender.opt = p.OpenOutPort("OUT")
}

func (sender *Sender) Execute(p *core.Process) {
	fmt.Println(p.Name + " started")
	icpkt := p.Receive(sender.ipt)
	j, _ := strconv.Atoi(icpkt.Contents.(string))
	p.Discard(icpkt)
	var pkt *core.Packet
	for i := 0; i < j; i++ {
		pkt = p.Create("IP - # " + strconv.Itoa(i) + " (" + p.Name + ")")
		p.Send(sender.opt.(*core.OutPort).Conn, pkt)
	}
	fmt.Println(p.Name + " ended")
}

//func (sender *Sender) GetMustRun() bool {
//	return false
//}
