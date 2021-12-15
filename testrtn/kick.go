package testrtn

import (
	//"fmt"

	"github.com/tyoung3/gofbp"
)

type Kick struct {
	opt gofbp.OutputConn
}

func (kick *Kick) Setup(p *gofbp.Process) {
	kick.opt = p.OpenOutPort("OUT")
}

func (kick *Kick) Execute(p *gofbp.Process) {

	var pkt = p.Create("Kicker IP")
	p.Send(kick.opt, pkt)

}
