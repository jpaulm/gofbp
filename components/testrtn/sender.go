/*Package testrtn tests gofbp code.*/
package testrtn

import (
	//"fmt"
	"strconv"
	
	"github.com/jpaulm/gofbp/core"
)

/*Sender type defines ipt and opt for process send*/
type Sender struct {
	ipt core.InputConn
	opt core.OutputConn
}

/*Setup function initializes a source process.*/
func (sender *Sender) Setup(p *core.Process) {
	sender.ipt = p.OpenInPort("COUNT")
	sender.opt = p.OpenOutPort("OUT")
}

/*Execute function launches a source process.*/
func (sender *Sender) Execute(p *core.Process) {
	
	icpkt := p.Receive(sender.ipt)
	j, _ := strconv.Atoi(icpkt.Contents.(string))
	p.Discard(icpkt)
	p.Close(sender.ipt)

	var pkt *core.Packet
	for i := 0; i < j; i++ {
		pkt = p.Create("IP - # " + strconv.Itoa(i))
		p.Send(sender.opt, pkt)
	}
	
}
