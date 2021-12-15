package testrtn

import (
	//"fmt"
	"strconv"

	"github.com/tyoung3/gofbp"
)

type Counter struct {
	ipt gofbp.InputConn
	cnt gofbp.OutputConn
	opt gofbp.OutputConn
}

func (counter *Counter) Setup(p *gofbp.Process) {
	counter.ipt = p.OpenInPort("IN")
	counter.cnt = p.OpenOutPort("COUNT")
	counter.opt = p.OpenOutPortOptional("OUT")
}

func (Counter) MustRun() {}

func (counter *Counter) Execute(p *gofbp.Process) {

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

}
