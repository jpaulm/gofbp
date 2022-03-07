package testrtn

import (
	//"fmt"
	"encoding/json"

	"github.com/jpaulm/gofbp/core"
)

type Marshal struct {
	ipt core.InputConn
	out core.OutputConn
}

func (marshal *Marshal) Setup(p *core.Process) {
	marshal.ipt = p.OpenInPort("IN")
	marshal.out = p.OpenOutPort("OUT")
}

func (Marshal) MustRun() {}

func (marshal *Marshal) Execute(p *core.Process) {

	for {
		var pkt = p.Receive(marshal.ipt)
		if pkt == nil {
			break
		}
		if pkt.PktType == core.OpenBracket || pkt.PktType == core.CloseBracket || pkt.PktType == core.Signal {
			p.Send(marshal.out, pkt)
		} else {
			json, _ := json.Marshal(pkt.Contents)
			p.Send(marshal.out, p.Create(string(json)))
			p.Discard(pkt)
		}
	}

}
