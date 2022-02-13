package components

import (
	//"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/jpaulm/gofbp/core"
)

/*IntSenderWDelay type defines ipt and opt for process*/
type IntSenderWDelay struct {
	ipt core.InputConn
	opt core.OutputConn
}

/*Setup function initializes a source process.*/
func (intsender *IntSenderWDelay) Setup(p *core.Process) {
	intsender.ipt = p.OpenInPort("COUNT")
	intsender.opt = p.OpenOutPort("OUT")
}

/*Execute function launches a source process.*/
func (intsender *IntSenderWDelay) Execute(p *core.Process) {

	icpkt := p.Receive(intsender.ipt)
	j, _ := strconv.Atoi(icpkt.Contents.(string))
	p.Discard(icpkt)
	p.Close(intsender.ipt)

	var pkt *core.Packet
	for i := 0; i < j; i++ {
		pkt = p.Create(i)
		p.Send(intsender.opt, pkt)
		time.Sleep(time.Duration(rand.Int31n(200)) * time.Millisecond)
	}

}
