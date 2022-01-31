/*Package testrtn tests gofbp code.*/
package testrtn

import (
	//"fmt"
	"strconv"

	"github.com/jpaulm/gofbp/core"
)

/*Sender type defines ipt and opt for process send*/
type IntSender struct {
	ipt core.InputConn
	opt core.OutputConn
}

/*Setup function initializes a source process.*/
func (intsender *IntSender) Setup(p *core.Process) {
	intsender.ipt = p.OpenInPort("COUNT")
	intsender.opt = p.OpenOutPort("OUT")
}

/*Execute function launches a source process.*/
func (intsender *IntSender) Execute(p *core.Process) {

	icpkt := p.Receive(intsender.ipt)
	j, _ := strconv.Atoi(icpkt.Contents.(string))
	p.Discard(icpkt)
	p.Close(intsender.ipt)

	var pkt *core.Packet
	for i := 0; i < j; i++ {
		pkt = p.Create(i)
		p.Send(intsender.opt, pkt)
	}

}
