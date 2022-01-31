package testrtn

import (
	"github.com/jpaulm/gofbp/core"
)

type LoadBalance struct {
	ipt core.InputConn
	out core.OutputArrayConn
}

func (loadbalance *LoadBalance) Setup(p *core.Process) {
	loadbalance.ipt = p.OpenInPort("IN")
	loadbalance.out = p.OpenOutArrayPort("OUT")
}

func (loadbalance *LoadBalance) Execute(p *core.Process) {

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

		if pkt.PktType == core.OpenBracket {
			inSubstream = true
			level++
		}

		if pkt.PktType == core.CloseBracket && level == 1 {
			inSubstream = false
			level--
		}

		opt := loadbalance.out.GetArrayItem(i)

		p.Send(opt, pkt)
	}

}
