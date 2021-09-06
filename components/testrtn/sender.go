package testrtn

import (
	"fmt"
	"strconv"

	"github.com/jpaulm/gofbp/core"
)

type Sender struct {
	opt *core.OutPort
}

func NewSender() *Sender {
	return &Sender{}
}

func (sender *Sender) OpenPorts(p *core.Process) {
	sender.opt = p.OpenOutPort("OUT")
}

func (sender *Sender) Execute(p *core.Process) {
	fmt.Println(p.Name + " started")
	var pkt *core.Packet
	for i := 0; i < 15; i++ {
		pkt = p.Create("IP - # " + strconv.Itoa(i) + " (" + p.Name + ")")
		p.Send(sender.opt.Conn, pkt)
	}
	fmt.Println(p.Name + " ended")
}
