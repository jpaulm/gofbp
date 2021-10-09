package testrtn

import (
	//"fmt"
	"strconv"

	"github.com/jpaulm/gofbp/core"
)

type Counter struct {
	ipt core.InputConn
	opt core.OutputConn
	cnt core.OutputConn
}

func (counter *Counter) Setup(p *core.Process) {
	counter.ipt = p.OpenInPort("IN")
	counter.cnt = p.OpenOutPort("COUNT")
	counter.opt = p.OpenOutPortOptional("OUT")
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
		p.Send(counter.opt, pkt)
		count++
	}

	pkt := p.Create(strconv.Itoa(count))
	p.Send(counter.cnt, pkt)

	//fmt.Println(p.GetName() + " ended")
}
