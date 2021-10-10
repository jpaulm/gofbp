package testrtn

import (
	//"fmt"
	"strconv"

	"github.com/jpaulm/gofbp/core"
)

type Counter struct {
	ipt *core.InPort
	cnt *core.OutPort
	opt *core.OutPort
}

func (counter *Counter) Setup(p *core.Process) {
	counter.ipt = p.OpenInPort("IN")
	counter.cnt = p.OpenOutPort("COUNT")
	counter.opt = p.OpenOutPort("OUT", "opt")
}

func (Counter) MustRun() {}

func (counter *Counter) Execute(p *core.Process) {
	//fmt.Println(p.GetName() + " started")

	count := 0

	for {
		var pkt = p.Receive(counter.ipt)
		if pkt == nil {
			break
		}
		if counter.opt.IsConnected() {
			p.Send(counter.opt, pkt)
		} else {
			p.Discard(pkt)
		}

		count++
	}

	pkt := p.Create(strconv.Itoa(count))
	p.Send(counter.cnt, pkt)

	//fmt.Println(p.GetName() + " ended")
}
