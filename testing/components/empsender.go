package components

import (
	//"fmt"
	"strconv"

	"github.com/jpaulm/gofbp/core"
)

/*EmpSender type defines ipt and opt for process send*/
type EmpSender struct {
	cnt  core.InputConn
	data core.InputConn
	opt  core.OutputConn
}

/*Setup function initializes a source process.*/
func (empsender *EmpSender) Setup(p *core.Process) {
	empsender.cnt = p.OpenInPort("COUNT")
	empsender.data = p.OpenInPort("DATA")
	empsender.opt = p.OpenOutPort("OUT")
}

/*Execute function launches a source process.*/
func (empsender *EmpSender) Execute(p *core.Process) {

	icpkt := p.Receive(empsender.cnt)
	j, _ := strconv.Atoi(icpkt.Contents.(string))
	p.Discard(icpkt)
	p.Close(empsender.cnt)

	icpkt = p.Receive(empsender.data)
	emp := icpkt.Contents
	p.Discard(icpkt)
	p.Close(empsender.data)

	var pkt *core.Packet
	for i := 0; i < j; i++ {
		pkt = p.Create(emp)
		p.Send(empsender.opt, pkt)
	}

}
