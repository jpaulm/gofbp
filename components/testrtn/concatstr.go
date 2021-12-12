package testrtn

import (
	"github.com/jpaulm/gofbp/core"
)

type ConcatStr struct {
	ipt core.InputArrayConn
	opt core.OutputConn
}

func (concatstr *ConcatStr) Setup(p *core.Process) {
	concatstr.ipt = p.OpenInArrayPort("IN")
	concatstr.opt = p.OpenOutPort("OUT")
}

func (concatstr *ConcatStr) Execute(p *core.Process) {

	for i := 0; i < concatstr.ipt.ArrayLength(); i++ {

		for {
			conn := concatstr.ipt.GetArrayItem(i)
			if conn == nil {
				continue
			}
			var pkt = p.Receive(conn)
			if pkt == nil {
				break
			}
			//fmt.Println("Output: ", pkt.Contents)
			p.Send(concatstr.opt, pkt)
		}
	}

}
