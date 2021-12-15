package testrtn

import (
	"github.com/tyoung3/gofbp"
)

type LoadBalance struct {
	ipt gofbp.InputConn
	out gofbp.OutputArrayConn
}

func (loadbalance *LoadBalance) Setup(p *gofbp.Process) {
	loadbalance.ipt = p.OpenInPort("IN")
	loadbalance.out = p.OpenOutArrayPort("OUT")
}

func (loadbalance *LoadBalance) Execute(p *gofbp.Process) {

	var i int
	var level int
	var inSubstream bool
	for {
		var pkt = p.Receive(loadbalance.ipt)
		if pkt == nil {
			break
		}

		if !inSubstream {
			i = loadbalance.out.GetItemWithFewestIPs()
		}

		if pkt.PktType == gofbp.OpenBracket {
			inSubstream = true
			level++
		}

		if pkt.PktType == gofbp.CloseBracket && level == 1 {
			inSubstream = false
			level--
		}

		opt := loadbalance.out.GetArrayItem(i)

		p.Send(opt, pkt)
	}

}
